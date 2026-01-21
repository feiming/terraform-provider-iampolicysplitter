// Copyright (c) 2026 Fabian Chong
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IamPolicySplitterDataSource{}

func NewIamPolicySplitterDataSource() datasource.DataSource {
	return &IamPolicySplitterDataSource{}
}

// IamPolicySplitterDataSource defines the data source implementation.
type IamPolicySplitterDataSource struct {
	client *http.Client
}

// IamPolicySplitterDataSourceModel describes the data source data model.
type IamPolicySplitterDataSourceModel struct {
	PolicyJson    types.String `tfsdk:"policy_json"`
	MaxChars      types.Int64  `tfsdk:"max_chars"`
	SplitPolicies types.List   `tfsdk:"split_policies"`
	Id            types.String `tfsdk:"id"`
}

// IAMPolicy represents an AWS IAM policy structure.
type IAMPolicy struct {
	Version   string               `json:"Version"`
	Statement []IAMPolicyStatement `json:"Statement"`
	Id        *string              `json:"Id,omitempty"`
}

// IAMPolicyStatement represents a single IAM policy statement.
type IAMPolicyStatement struct {
	Sid          *string                `json:"Sid,omitempty"`
	Effect       string                 `json:"Effect"`
	Principal    interface{}            `json:"Principal,omitempty"`
	NotPrincipal interface{}            `json:"NotPrincipal,omitempty"`
	Action       interface{}            `json:"Action,omitempty"`
	NotAction    interface{}            `json:"NotAction,omitempty"`
	Resource     interface{}            `json:"Resource,omitempty"`
	NotResource  interface{}            `json:"NotResource,omitempty"`
	Condition    map[string]interface{} `json:"Condition,omitempty"`
}

// StatementWithSize wraps a statement with its character count.
type StatementWithSize struct {
	Statement IAMPolicyStatement
	Size      int
}

func (d *IamPolicySplitterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_policy_splitter"
}

func (d *IamPolicySplitterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Splits an IAM policy by statements and rearranges them to fit within character limits. " +
			"This is useful when you have a large IAM policy that exceeds AWS limits (6,144 characters for managed policies, " +
			"2,048 characters for inline policies).",

		Attributes: map[string]schema.Attribute{
			"policy_json": schema.StringAttribute{
				MarkdownDescription: "The IAM policy JSON string to split",
				Required:            true,
			},
			"max_chars": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of characters allowed per policy. Defaults to 6144 (managed policy limit)",
				Optional:            true,
			},
			"split_policies": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of split policy JSON strings, each within the character limit",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for the data source",
				Computed:            true,
			},
		},
	}
}

func (d *IamPolicySplitterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IamPolicySplitterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IamPolicySplitterDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set default max_chars if not provided
	maxChars := int64(6144) // Default AWS managed policy limit
	if !data.MaxChars.IsNull() && !data.MaxChars.IsUnknown() {
		maxChars = data.MaxChars.ValueInt64()
	}

	if maxChars <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"max_chars must be greater than 0",
		)
		return
	}

	// Parse the input policy JSON
	var policy IAMPolicy
	if err := json.Unmarshal([]byte(data.PolicyJson.ValueString()), &policy); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Policy JSON",
			fmt.Sprintf("Failed to parse policy JSON: %s", err.Error()),
		)
		return
	}

	// Validate policy structure
	if policy.Version == "" {
		resp.Diagnostics.AddError(
			"Invalid Policy",
			"Policy must have a Version field",
		)
		return
	}

	if len(policy.Statement) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Policy",
			"Policy must have at least one Statement",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Splitting policy with %d statements, max_chars=%d", len(policy.Statement), maxChars))

	// Calculate size for each statement
	statementsWithSize := make([]StatementWithSize, 0, len(policy.Statement))
	for i, stmt := range policy.Statement {
		stmtPolicy := IAMPolicy{
			Version:   policy.Version,
			Statement: []IAMPolicyStatement{stmt},
			Id:        policy.Id,
		}
		stmtJSON, err := json.Marshal(stmtPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				"Internal Error",
				fmt.Sprintf("Failed to marshal statement %d: %s", i, err.Error()),
			)
			return
		}
		size := len(string(stmtJSON))
		statementsWithSize = append(statementsWithSize, StatementWithSize{
			Statement: stmt,
			Size:      size,
		})
		if ctx != nil {
			tflog.Debug(ctx, fmt.Sprintf("Statement %d size: %d characters", i, size))
		}
	}

	// Check if any single statement exceeds the limit
	for i, stmt := range statementsWithSize {
		if stmt.Size > int(maxChars) {
			resp.Diagnostics.AddError(
				"Statement Too Large",
				fmt.Sprintf("Statement %d (%d characters) exceeds the maximum character limit (%d). "+
					"Individual statements cannot be split further.", i, stmt.Size, maxChars),
			)
			return
		}
	}

	// Use bin-packing algorithm to group statements
	splitPolicies := d.packStatements(ctx, policy, statementsWithSize, int(maxChars))

	// Convert to list of strings
	splitPolicyStrings := make([]types.String, len(splitPolicies))
	for i, p := range splitPolicies {
		policyJSON, err := json.Marshal(p)
		if err != nil {
			resp.Diagnostics.AddError(
				"Internal Error",
				fmt.Sprintf("Failed to marshal split policy %d: %s", i, err.Error()),
			)
			return
		}
		splitPolicyStrings[i] = types.StringValue(string(policyJSON))
	}

	// Set the split policies
	splitPoliciesList, diags := types.ListValueFrom(ctx, types.StringType, splitPolicyStrings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.SplitPolicies = splitPoliciesList
	data.Id = types.StringValue(fmt.Sprintf("split-%d-%d", len(policy.Statement), maxChars))

	tflog.Info(ctx, fmt.Sprintf("Successfully split policy into %d policies", len(splitPolicies)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// packStatements uses a bin-packing algorithm to group statements into policies
// that fit within the character limit. Uses first-fit decreasing algorithm.
func (d *IamPolicySplitterDataSource) packStatements(ctx context.Context, originalPolicy IAMPolicy, statements []StatementWithSize, maxChars int) []IAMPolicy {
	// Sort statements by size (descending) for better bin-packing
	sortedStatements := make([]StatementWithSize, len(statements))
	copy(sortedStatements, statements)
	sort.Slice(sortedStatements, func(i, j int) bool {
		return sortedStatements[i].Size > sortedStatements[j].Size
	})

	// Track which statements are in which policy
	type policyBin struct {
		statements  []IAMPolicyStatement
		currentSize int
	}

	bins := []policyBin{}

	// First-fit decreasing algorithm
	for _, stmtWithSize := range sortedStatements {
		placed := false
		// Try to place in existing bin
		for i := range bins {
			// Create a test policy to check size
			testPolicy := IAMPolicy{
				Version:   originalPolicy.Version,
				Statement: append(bins[i].statements, stmtWithSize.Statement),
				Id:        originalPolicy.Id,
			}
			testJSON, _ := json.Marshal(testPolicy)
			testSize := len(string(testJSON))

			if testSize <= maxChars {
				bins[i].statements = append(bins[i].statements, stmtWithSize.Statement)
				bins[i].currentSize = testSize
				placed = true
				break
			}
		}

		// If couldn't place in existing bin, create new bin
		if !placed {
			newPolicy := IAMPolicy{
				Version:   originalPolicy.Version,
				Statement: []IAMPolicyStatement{stmtWithSize.Statement},
				Id:        originalPolicy.Id,
			}
			newJSON, _ := json.Marshal(newPolicy)
			bins = append(bins, policyBin{
				statements:  []IAMPolicyStatement{stmtWithSize.Statement},
				currentSize: len(string(newJSON)),
			})
		}
	}

	// Convert bins to policies
	result := make([]IAMPolicy, len(bins))
	for i, bin := range bins {
		result[i] = IAMPolicy{
			Version:   originalPolicy.Version,
			Statement: bin.statements,
			Id:        originalPolicy.Id,
		}
		if ctx != nil {
			tflog.Debug(ctx, fmt.Sprintf("Policy %d: %d statements, %d characters", i+1, len(bin.statements), bin.currentSize))
		}
	}

	return result
}
