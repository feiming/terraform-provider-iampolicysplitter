// Copyright (c) 2026 Fabian Chong
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"testing"
)

func TestPackStatements(t *testing.T) {
	ds := &IamPolicySplitterDataSource{}

	// Create a context for testing (logger not needed for unit tests)
	ctx := context.Background()

	// Create a test policy with multiple statements
	originalPolicy := IAMPolicy{
		Version: "2012-10-17",
		Statement: []IAMPolicyStatement{
			{
				Effect: "Allow",
				Action: []string{"s3:GetObject"},
				Resource: "arn:aws:s3:::bucket1/*",
			},
			{
				Effect: "Allow",
				Action: []string{"s3:PutObject"},
				Resource: "arn:aws:s3:::bucket2/*",
			},
			{
				Effect: "Allow",
				Action: []string{"s3:DeleteObject"},
				Resource: "arn:aws:s3:::bucket3/*",
			},
		},
	}

	// Calculate sizes
	statementsWithSize := make([]StatementWithSize, len(originalPolicy.Statement))
	for i, stmt := range originalPolicy.Statement {
		stmtPolicy := IAMPolicy{
			Version:   originalPolicy.Version,
			Statement: []IAMPolicyStatement{stmt},
		}
		stmtJSON, _ := json.Marshal(stmtPolicy)
		statementsWithSize[i] = StatementWithSize{
			Statement: stmt,
			Size:      len(string(stmtJSON)),
		}
	}

	// Test with a very small limit to force splitting
	maxChars := 500
	result := ds.packStatements(ctx, originalPolicy, statementsWithSize, maxChars)

	if len(result) == 0 {
		t.Fatal("Expected at least one policy, got 0")
	}

	// Verify all policies are within limit
	for i, policy := range result {
		policyJSON, err := json.Marshal(policy)
		if err != nil {
			t.Fatalf("Failed to marshal policy %d: %v", i, err)
		}
		size := len(string(policyJSON))
		if size > maxChars {
			t.Errorf("Policy %d exceeds limit: %d > %d", i, size, maxChars)
		}
	}

	// Verify all statements are present
	totalStatements := 0
	for _, policy := range result {
		totalStatements += len(policy.Statement)
	}
	if totalStatements != len(originalPolicy.Statement) {
		t.Errorf("Expected %d statements total, got %d", len(originalPolicy.Statement), totalStatements)
	}
}

func TestPackStatementsWithLargeLimit(t *testing.T) {
	ds := &IamPolicySplitterDataSource{}

	// Create a context for testing (logger not needed for unit tests)
	ctx := context.Background()

	originalPolicy := IAMPolicy{
		Version: "2012-10-17",
		Statement: []IAMPolicyStatement{
			{
				Effect:   "Allow",
				Action:   []string{"s3:GetObject"},
				Resource: "arn:aws:s3:::bucket1/*",
			},
		},
	}

	statementsWithSize := make([]StatementWithSize, len(originalPolicy.Statement))
	for i, stmt := range originalPolicy.Statement {
		stmtPolicy := IAMPolicy{
			Version:   originalPolicy.Version,
			Statement: []IAMPolicyStatement{stmt},
		}
		stmtJSON, _ := json.Marshal(stmtPolicy)
		statementsWithSize[i] = StatementWithSize{
			Statement: stmt,
			Size:      len(string(stmtJSON)),
		}
	}

	// Test with a large limit - should result in single policy
	maxChars := 10000
	result := ds.packStatements(ctx, originalPolicy, statementsWithSize, maxChars)

	if len(result) != 1 {
		t.Errorf("Expected 1 policy with large limit, got %d", len(result))
	}

	if len(result[0].Statement) != 1 {
		t.Errorf("Expected 1 statement in policy, got %d", len(result[0].Statement))
	}
}
