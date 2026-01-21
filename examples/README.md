# Examples

This directory contains examples demonstrating how to use the `iampolicysplitter` Terraform provider. These examples are used for documentation and can also be run/tested manually via the Terraform CLI.

## Available Examples

### Data Sources

#### `iampolicysplitter_iam_policy_splitter`

The IAM policy splitter data source splits large IAM policies into multiple smaller policies that comply with AWS character limits.

**Location:** `data-sources/iam_policy_splitter/data-source.tf`

**Example Usage:**

```hcl
data "iampolicysplitter_iam_policy_splitter" "example" {
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = "arn:aws:s3:::example-bucket/*"
      }
    ]
  })
  max_chars = 2048  # Inline policy limit (or 6144 for managed policies)
}

output "split_policies" {
  value = data.iampolicysplitter_iam_policy_splitter.example.split_policies
}
```

**Features Demonstrated:**
- Splitting a large IAM policy with multiple statements
- Configuring character limits (2048 for inline policies, 6144 for managed policies)
- Accessing the resulting split policies as a list

## Documentation Generation

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** - Example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** - Example file for the named data source page
* **resources/`full resource name`/resource.tf** - Example file for the named resource page

## Running Examples

To run an example:

1. Navigate to the example directory:
   ```bash
   cd examples/data-sources/iam_policy_splitter
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Review the plan:
   ```bash
   terraform plan
   ```

4. Apply (if needed):
   ```bash
   terraform apply
   ```
