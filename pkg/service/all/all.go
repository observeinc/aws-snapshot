package all

import (
	_ "github.com/observeinc/aws-snapshot/pkg/service/apigateway"
	_ "github.com/observeinc/aws-snapshot/pkg/service/autoscaling"
	_ "github.com/observeinc/aws-snapshot/pkg/service/cloudformation"
	_ "github.com/observeinc/aws-snapshot/pkg/service/cloudfront"
	_ "github.com/observeinc/aws-snapshot/pkg/service/cloudwatchlogs"
	_ "github.com/observeinc/aws-snapshot/pkg/service/dynamodb"
	_ "github.com/observeinc/aws-snapshot/pkg/service/ec2"
	_ "github.com/observeinc/aws-snapshot/pkg/service/ecs"
	_ "github.com/observeinc/aws-snapshot/pkg/service/efs"
	_ "github.com/observeinc/aws-snapshot/pkg/service/eks"
	_ "github.com/observeinc/aws-snapshot/pkg/service/elasticache"
	_ "github.com/observeinc/aws-snapshot/pkg/service/elasticbeanstalk"
	_ "github.com/observeinc/aws-snapshot/pkg/service/elasticloadbalancing"
	_ "github.com/observeinc/aws-snapshot/pkg/service/eventbridge"
	_ "github.com/observeinc/aws-snapshot/pkg/service/firehose"
	_ "github.com/observeinc/aws-snapshot/pkg/service/iam"
	_ "github.com/observeinc/aws-snapshot/pkg/service/kinesis"
	_ "github.com/observeinc/aws-snapshot/pkg/service/kms"
	_ "github.com/observeinc/aws-snapshot/pkg/service/lambda"
	_ "github.com/observeinc/aws-snapshot/pkg/service/organizations"
	_ "github.com/observeinc/aws-snapshot/pkg/service/rds"
	_ "github.com/observeinc/aws-snapshot/pkg/service/redshift"
	_ "github.com/observeinc/aws-snapshot/pkg/service/route53"
	_ "github.com/observeinc/aws-snapshot/pkg/service/s3"
	_ "github.com/observeinc/aws-snapshot/pkg/service/secretsmanager"
	_ "github.com/observeinc/aws-snapshot/pkg/service/securityhub"
	_ "github.com/observeinc/aws-snapshot/pkg/service/sns"
	_ "github.com/observeinc/aws-snapshot/pkg/service/sqs"
	_ "github.com/observeinc/aws-snapshot/pkg/service/synthetics"
)
