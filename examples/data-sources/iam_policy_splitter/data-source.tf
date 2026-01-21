data "iampolicysplitter_iam_policy_splitter" "example" {
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3BucketAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:GetObjectVersion",
          "s3:PutObjectAcl",
          "s3:GetObjectAcl",
          "s3:GetObjectVersionAcl"
        ]
        Resource = [
          "arn:aws:s3:::production-data-bucket/*",
          "arn:aws:s3:::staging-data-bucket/*",
          "arn:aws:s3:::development-data-bucket/*"
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
          "arn:aws:s3:::production-data-bucket",
          "arn:aws:s3:::staging-data-bucket",
          "arn:aws:s3:::development-data-bucket"
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
      },
      {
        Sid    = "EC2VolumeManagement"
        Effect = "Allow"
        Action = [
          "ec2:DescribeVolumes",
          "ec2:DescribeVolumeStatus",
          "ec2:DescribeVolumeAttribute",
          "ec2:CreateVolume",
          "ec2:AttachVolume",
          "ec2:DetachVolume",
          "ec2:DeleteVolume",
          "ec2:ModifyVolumeAttribute",
          "ec2:CreateSnapshot",
          "ec2:DescribeSnapshots",
          "ec2:DeleteSnapshot"
        ]
        Resource = "*"
      },
      {
        Sid    = "EC2SecurityGroupManagement"
        Effect = "Allow"
        Action = [
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSecurityGroupRules",
          "ec2:CreateSecurityGroup",
          "ec2:DeleteSecurityGroup",
          "ec2:AuthorizeSecurityGroupIngress",
          "ec2:AuthorizeSecurityGroupEgress",
          "ec2:RevokeSecurityGroupIngress",
          "ec2:RevokeSecurityGroupEgress",
          "ec2:UpdateSecurityGroupRuleDescriptionsIngress",
          "ec2:UpdateSecurityGroupRuleDescriptionsEgress"
        ]
        Resource = "*"
      },
      {
        Sid    = "EC2NetworkInterfaceManagement"
        Effect = "Allow"
        Action = [
          "ec2:DescribeNetworkInterfaces",
          "ec2:DescribeNetworkInterfaceAttribute",
          "ec2:CreateNetworkInterface",
          "ec2:DeleteNetworkInterface",
          "ec2:AttachNetworkInterface",
          "ec2:DetachNetworkInterface",
          "ec2:ModifyNetworkInterfaceAttribute"
        ]
        Resource = "*"
      },
      {
        Sid    = "RDSDatabaseAccess"
        Effect = "Allow"
        Action = [
          "rds:DescribeDBInstances",
          "rds:DescribeDBClusters",
          "rds:DescribeDBSubnetGroups",
          "rds:DescribeDBParameterGroups",
          "rds:DescribeDBClusterParameterGroups",
          "rds:DescribeDBEngineVersions",
          "rds:DescribeDBLogFiles",
          "rds:DownloadDBLogFilePortion",
          "rds:CreateDBInstance",
          "rds:ModifyDBInstance",
          "rds:DeleteDBInstance",
          "rds:RebootDBInstance",
          "rds:StartDBInstance",
          "rds:StopDBInstance"
        ]
        Resource = [
          "arn:aws:rds:us-east-1:*:db:production-db-*",
          "arn:aws:rds:us-east-1:*:db:staging-db-*",
          "arn:aws:rds:us-west-2:*:db:production-db-*",
          "arn:aws:rds:us-west-2:*:db:staging-db-*"
        ]
      },
      {
        Sid    = "LambdaFunctionManagement"
        Effect = "Allow"
        Action = [
          "lambda:ListFunctions",
          "lambda:GetFunction",
          "lambda:GetFunctionConfiguration",
          "lambda:GetFunctionCodeSigningConfig",
          "lambda:GetFunctionConcurrency",
          "lambda:GetFunctionEventInvokeConfig",
          "lambda:GetFunctionUrlConfig",
          "lambda:InvokeFunction",
          "lambda:CreateFunction",
          "lambda:UpdateFunctionCode",
          "lambda:UpdateFunctionConfiguration",
          "lambda:DeleteFunction",
          "lambda:AddPermission",
          "lambda:RemovePermission",
          "lambda:PutFunctionConcurrency",
          "lambda:DeleteFunctionConcurrency"
        ]
        Resource = [
          "arn:aws:lambda:us-east-1:*:function:api-handler-*",
          "arn:aws:lambda:us-east-1:*:function:data-processor-*",
          "arn:aws:lambda:us-west-2:*:function:api-handler-*",
          "arn:aws:lambda:us-west-2:*:function:data-processor-*"
        ]
      },
      {
        Sid    = "CloudWatchLogsAccess"
        Effect = "Allow"
        Action = [
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:GetLogEvents",
          "logs:FilterLogEvents",
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DeleteLogGroup",
          "logs:DeleteLogStream",
          "logs:PutRetentionPolicy",
          "logs:PutMetricFilter",
          "logs:DeleteMetricFilter"
        ]
        Resource = [
          "arn:aws:logs:us-east-1:*:log-group:/aws/lambda/api-handler-*",
          "arn:aws:logs:us-east-1:*:log-group:/aws/lambda/data-processor-*",
          "arn:aws:logs:us-west-2:*:log-group:/aws/lambda/api-handler-*",
          "arn:aws:logs:us-west-2:*:log-group:/aws/lambda/data-processor-*",
          "arn:aws:logs:us-east-1:*:log-group:/aws/ec2/application-logs-*",
          "arn:aws:logs:us-west-2:*:log-group:/aws/ec2/application-logs-*"
        ]
      },
      {
        Sid    = "IAMRoleManagement"
        Effect = "Allow"
        Action = [
          "iam:GetRole",
          "iam:GetRolePolicy",
          "iam:ListRolePolicies",
          "iam:ListAttachedRolePolicies",
          "iam:ListInstanceProfilesForRole",
          "iam:CreateRole",
          "iam:UpdateRole",
          "iam:DeleteRole",
          "iam:PutRolePolicy",
          "iam:DeleteRolePolicy",
          "iam:AttachRolePolicy",
          "iam:DetachRolePolicy",
          "iam:TagRole",
          "iam:UntagRole",
          "iam:ListRoleTags"
        ]
        Resource = [
          "arn:aws:iam::*:role/application-role-*",
          "arn:aws:iam::*:role/service-role-*",
          "arn:aws:iam::*:role/lambda-execution-role-*"
        ]
      },
      {
        Sid    = "SecretsManagerAccess"
        Effect = "Allow"
        Action = [
          "secretsmanager:DescribeSecret",
          "secretsmanager:GetSecretValue",
          "secretsmanager:ListSecretVersionIds",
          "secretsmanager:GetResourcePolicy",
          "secretsmanager:CreateSecret",
          "secretsmanager:UpdateSecret",
          "secretsmanager:PutSecretValue",
          "secretsmanager:DeleteSecret",
          "secretsmanager:RestoreSecret",
          "secretsmanager:RotateSecret",
          "secretsmanager:TagResource",
          "secretsmanager:UntagResource"
        ]
        Resource = [
          "arn:aws:secretsmanager:us-east-1:*:secret:database-credentials-*",
          "arn:aws:secretsmanager:us-east-1:*:secret:api-keys-*",
          "arn:aws:secretsmanager:us-west-2:*:secret:database-credentials-*",
          "arn:aws:secretsmanager:us-west-2:*:secret:api-keys-*"
        ]
      },
      {
        Sid    = "DynamoDBTableAccess"
        Effect = "Allow"
        Action = [
          "dynamodb:DescribeTable",
          "dynamodb:DescribeTableReplicaAutoScaling",
          "dynamodb:DescribeTimeToLive",
          "dynamodb:DescribeContinuousBackups",
          "dynamodb:DescribeContributorInsights",
          "dynamodb:ListTables",
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:CreateTable",
          "dynamodb:UpdateTable",
          "dynamodb:DeleteTable"
        ]
        Resource = [
          "arn:aws:dynamodb:us-east-1:*:table/user-data-*",
          "arn:aws:dynamodb:us-east-1:*:table/session-data-*",
          "arn:aws:dynamodb:us-west-2:*:table/user-data-*",
          "arn:aws:dynamodb:us-west-2:*:table/session-data-*"
        ]
      },
      {
        Sid    = "SQSQueueManagement"
        Effect = "Allow"
        Action = [
          "sqs:GetQueueAttributes",
          "sqs:GetQueueUrl",
          "sqs:ListQueues",
          "sqs:ListQueueTags",
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:ChangeMessageVisibility",
          "sqs:CreateQueue",
          "sqs:DeleteQueue",
          "sqs:SetQueueAttributes",
          "sqs:TagQueue",
          "sqs:UntagQueue",
          "sqs:PurgeQueue"
        ]
        Resource = [
          "arn:aws:sqs:us-east-1:*:task-queue-*",
          "arn:aws:sqs:us-east-1:*:notification-queue-*",
          "arn:aws:sqs:us-west-2:*:task-queue-*",
          "arn:aws:sqs:us-west-2:*:notification-queue-*"
        ]
      },
      {
        Sid    = "SNSTopicManagement"
        Effect = "Allow"
        Action = [
          "sns:GetTopicAttributes",
          "sns:ListTopics",
          "sns:ListSubscriptionsByTopic",
          "sns:ListTagsForResource",
          "sns:Publish",
          "sns:CreateTopic",
          "sns:DeleteTopic",
          "sns:SetTopicAttributes",
          "sns:Subscribe",
          "sns:Unsubscribe",
          "sns:TagResource",
          "sns:UntagResource"
        ]
        Resource = [
          "arn:aws:sns:us-east-1:*:application-notifications-*",
          "arn:aws:sns:us-east-1:*:system-alerts-*",
          "arn:aws:sns:us-west-2:*:application-notifications-*",
          "arn:aws:sns:us-west-2:*:system-alerts-*"
        ]
      },
      {
        Sid    = "KMSKeyManagement"
        Effect = "Allow"
        Action = [
          "kms:DescribeKey",
          "kms:ListKeys",
          "kms:ListAliases",
          "kms:ListGrants",
          "kms:ListKeyPolicies",
          "kms:GetKeyPolicy",
          "kms:GetKeyRotationStatus",
          "kms:Encrypt",
          "kms:Decrypt",
          "kms:ReEncrypt",
          "kms:GenerateDataKey",
          "kms:GenerateDataKeyWithoutPlaintext",
          "kms:CreateGrant",
          "kms:RevokeGrant",
          "kms:TagResource",
          "kms:UntagResource"
        ]
        Resource = [
          "arn:aws:kms:us-east-1:*:key/*",
          "arn:aws:kms:us-west-2:*:key/*"
        ]
        Condition = {
          StringEquals = {
            "kms:ViaService" = [
              "s3.us-east-1.amazonaws.com",
              "s3.us-west-2.amazonaws.com",
              "secretsmanager.us-east-1.amazonaws.com",
              "secretsmanager.us-west-2.amazonaws.com"
            ]
          }
        }
      },
      {
        Sid    = "CloudFormationStackManagement"
        Effect = "Allow"
        Action = [
          "cloudformation:DescribeStacks",
          "cloudformation:DescribeStackResources",
          "cloudformation:DescribeStackEvents",
          "cloudformation:ListStacks",
          "cloudformation:ListStackResources",
          "cloudformation:GetTemplate",
          "cloudformation:ValidateTemplate",
          "cloudformation:CreateStack",
          "cloudformation:UpdateStack",
          "cloudformation:DeleteStack",
          "cloudformation:CancelUpdateStack",
          "cloudformation:TagResource",
          "cloudformation:UntagResource"
        ]
        Resource = "*"
      },
      {
        Sid    = "APIGatewayManagement"
        Effect = "Allow"
        Action = [
          "apigateway:GET",
          "apigateway:POST",
          "apigateway:PUT",
          "apigateway:PATCH",
          "apigateway:DELETE",
          "apigateway:HEAD",
          "apigateway:OPTIONS"
        ]
        Resource = [
          "arn:aws:apigateway:us-east-1::/restapis/*",
          "arn:aws:apigateway:us-east-1::/restapis/*/*",
          "arn:aws:apigateway:us-west-2::/restapis/*",
          "arn:aws:apigateway:us-west-2::/restapis/*/*"
        ]
      }
    ]
  })
  max_chars = 2048
}

output "split_policies" {
  value = data.iampolicysplitter_iam_policy_splitter.example.split_policies
}

output "policy_count" {
  value = length(data.iampolicysplitter_iam_policy_splitter.example.split_policies)
}
