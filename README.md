# Terraform Provider IAM Policy Splitter

A Terraform provider that splits IAM policies by statements and rearranges them to fit within AWS IAM policy character limits.

This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Features

- **Splits IAM policies by statements**: Automatically extracts individual statements from an IAM policy
- **Character counting**: Accurately counts characters for each statement and policy
- **Bin-packing algorithm**: Uses first-fit decreasing algorithm to efficiently group statements into multiple policies
- **Configurable limits**: Supports both managed policy (6,144 chars) and inline policy (2,048 chars) limits

## Use Cases

When you have a large IAM policy that exceeds AWS limits:
- **Managed policies**: Maximum 6,144 characters
- **Inline policies**: Maximum 2,048 characters

This provider helps you automatically split large policies into multiple smaller policies that comply with AWS limits.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Basic Example

```hcl
data "iampolicysplitter_iam_policy_splitter" "example" {
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject", "s3:PutObject"]
        Resource = "arn:aws:s3:::example-bucket/*"
      },
      {
        Effect   = "Allow"
        Action   = ["s3:ListBucket"]
        Resource = "arn:aws:s3:::example-bucket"
      }
    ]
  })
  max_chars = 2048  # Inline policy limit
}

output "split_policies" {
  value = data.iampolicysplitter_iam_policy_splitter.example.split_policies
}
```

### Data Source: `iampolicysplitter_iam_policy_splitter`

#### Arguments

- `policy_json` (Required, String): The IAM policy JSON string to split
- `max_chars` (Optional, Number): Maximum number of characters allowed per policy. Defaults to 6144 (managed policy limit)

#### Attributes

- `split_policies` (Computed, List of Strings): List of split policy JSON strings, each within the character limit
- `id` (Computed, String): Identifier for the data source

### Example: Using Split Policies

```hcl
data "iampolicysplitter_iam_policy_splitter" "large_policy" {
  policy_json = var.large_iam_policy_json
  max_chars   = 6144
}

# Create multiple managed policies from the split
resource "aws_iam_policy" "split_policies" {
  count = length(data.iampolicysplitter_iam_policy_splitter.large_policy.split_policies)
  
  name   = "policy-${count.index + 1}"
  policy = data.iampolicysplitter_iam_policy_splitter.large_policy.split_policies[count.index]
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
