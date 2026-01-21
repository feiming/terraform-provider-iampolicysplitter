// Copyright (c) 2026 Fabian Chong
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccIamPolicySplitterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIamPolicySplitterDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.iampolicysplitter_iam_policy_splitter.test",
						tfjsonpath.New("split_policies"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
		},
	})
}

func TestAccIamPolicySplitterDataSourceWithLargePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIamPolicySplitterDataSourceConfigLarge,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.iampolicysplitter_iam_policy_splitter.test",
						tfjsonpath.New("split_policies"),
						knownvalue.ListSizeExact(5),
					),
				},
			},
		},
	})
}

const testAccIamPolicySplitterDataSourceConfig = `
data "iampolicysplitter_iam_policy_splitter" "test" {
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
  max_chars = 6144
}
`

const testAccIamPolicySplitterDataSourceConfigLarge = `
data "iampolicysplitter_iam_policy_splitter" "test" {
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3BucketAccess1"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:GetObjectVersion",
          "s3:PutObjectAcl",
          "s3:GetObjectAcl"
        ]
        Resource = [
          "arn:aws:s3:::production-data-bucket-1/*",
          "arn:aws:s3:::staging-data-bucket-1/*",
          "arn:aws:s3:::development-data-bucket-1/*"
        ]
      },
      {
        Sid    = "S3BucketAccess2"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:GetObjectVersion",
          "s3:PutObjectAcl",
          "s3:GetObjectAcl"
        ]
        Resource = [
          "arn:aws:s3:::production-data-bucket-2/*",
          "arn:aws:s3:::staging-data-bucket-2/*",
          "arn:aws:s3:::development-data-bucket-2/*"
        ]
      },
      {
        Sid    = "S3BucketAccess3"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:GetObjectVersion",
          "s3:PutObjectAcl",
          "s3:GetObjectAcl"
        ]
        Resource = [
          "arn:aws:s3:::production-data-bucket-3/*",
          "arn:aws:s3:::staging-data-bucket-3/*",
          "arn:aws:s3:::development-data-bucket-3/*"
        ]
      },
      {
        Sid    = "S3BucketList"
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "s3:ListBucketVersions",
          "s3:GetBucketLocation",
          "s3:GetBucketAcl",
          "s3:GetBucketVersioning"
        ]
        Resource = [
          "arn:aws:s3:::production-data-bucket-1",
          "arn:aws:s3:::staging-data-bucket-1",
          "arn:aws:s3:::development-data-bucket-1",
          "arn:aws:s3:::production-data-bucket-2",
          "arn:aws:s3:::staging-data-bucket-2",
          "arn:aws:s3:::development-data-bucket-2"
        ]
      },
      {
        Sid    = "EC2InstanceManagement"
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeInstanceStatus",
          "ec2:DescribeInstanceAttribute",
          "ec2:StartInstances",
          "ec2:StopInstances",
          "ec2:RebootInstances",
          "ec2:TerminateInstances",
          "ec2:RunInstances",
          "ec2:ModifyInstanceAttribute"
        ]
        Resource = "*"
        Condition = {
          StringEquals = {
            "ec2:Region" = ["us-east-1", "us-west-2", "eu-west-1"]
          }
        }
      }
    ]
  })
  max_chars = 500
}
`
