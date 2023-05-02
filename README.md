# AWS CDK Samples

This repository contains a collection of AWS CDK integration samples for various services. These samples are designed to provide documentation and guidance on how to implement integrations using AWS CDK. Please note that while these samples have been tested and validated, they are not intended for use in a production environment.

## Pre-requisites
* Go
* AWS CLI
* AWS CDK
* Node.js
* TypeScript
* AWS Account

## AWS Configuration
Configure your workstation with your credentials and an AWS region.
```bash
dev@dev:~$ aws configure
```

To create multiple accounts for AWS CLI:
```bash
dev@dev:~$ aws configure --profile profile_name
```

Provide your AWS access key ID, secret access key and default region when prompted. You can switch between accounts by passing the profile on the command.
```bash
dev@dev:~$ aws s3 ls --profile profile_name
```

When no `--profile` parameter provided in the command, `default` profile will be used.

### AWS SDK and CDK Defaults
* **AWS Region**: `us-east-1`
  * You can change the default value of the AWS region by updating the `AWS_REGION` in one of the files inside the ***aws_wrapper***
* **Stack Account and Region**: `process.env.CDK_DEFAULT_ACCOUNT/REGION`
  * You can change the default value in the TS file inside the ***bin*** folder

# Table of Contents

### [API Gateway](api-gateway/)
* [API Gateway Async → Lambda](api-gateway/api-gateway-async-lambda/README.md)
* [API Gateway CORS + API Key → Lambda](api-gateway/api-gateway-cors-lambda/README.md)
* [API Gateway → Lambda → DynamoDB](api-gateway/api-gateway-lambda-dynamodb/README.md)
* [API Gateway → Lambda → S3 Bucket](api-gateway/api-gateway-lambda-s3/README.md)
* [API Gateway → Lambda → SQS → Lambda](api-gateway/api-gateway-lambda-sqs/README.md)

### [CloudWatch](cloudwatch/)
* [CloudWatch Subscription → Lambda](cloudwatch/cloudwatch-subscription-lambda/README.md)

### [Event Bridge](event-bridge/)
* [EventBridge → Lambda](event-bridge/event-bridge-rule-lambda/README.md)

### [Lambda](lambda/)
* [Lambda → DynamoDB](lambda/lambda-dynamodb/README.md)
* [Lambda → S3](lambda/lambda-s3/README.md)
* [Lambda → Secrets Manager](lambda/lambda-secretsmanager/README.md)
* [Lambda → SNS](lambda/lambda-sns/README.md)
* [Lambda → SQS](lambda/lambda-sqs/README.md)

### [S3](s3/)
* [S3 Bucket → EventBridge Rule → Lambda](s3/s3-eventbridge-lambda/README.md)
* [S3 Bucket → Lambda → DynamoDB](s3/s3-lambda-dynamodb/README.md)
* [S3 Bucket → SNS → Email](s3/s3-sns/README.md)
* [S3 Presigned URLs](s3/s3-presigned-urls/README.md)
* [S3 Static Website hosting](s3/s3-website/README.md)

### [Step Functions](step-functions/)
* [Step Function with CloudWatch](step-functions/step-functions-cloudwatch/README.md)
* [Step Function with Lambda](step-functions/step-functions-lambda/README.md)