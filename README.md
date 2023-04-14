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

## [API Gateway](api-gateway/)
* [API Gateway Async → Lambda](api-gateway/api-gateway-async-lambda/README.md)
* [API Gateway → Lambda → S3 Bucket](api-gateway/api-gateway-lambda-s3/README.md)
* [API Gateway → Lambda → DynamoDB](api-gateway/api-gateway-lambda-dynamodb/README.md)
* [API Gateway → Lambda → SQS → Lambda](api-gateway/api-gateway-lambda-sqs/README.md)
* [API Gateway CORS + API Key → Lambda](api-gateway/api-gateway-cors-lambda/README.md)

## [S3](s3/)
* [S3 Static Website hosting](s3/s3-website/README.md)
* [S3 Bucket → Lambda → DynamoDB](s3/s3-lambda-dynamodb/README.md)
* [S3 Bucket → EventBridge Rule → Lambda](s3/s3-eventbridge-lambda/README.md)
* [S3 Bucket → SNS → Email](s3/s3-sns/README.md)
* [S3 Presigned URLs](s3/s3-presigned-urls/README.md)