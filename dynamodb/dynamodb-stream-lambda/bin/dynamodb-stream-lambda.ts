#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DynamodbStreamLambdaStack } from '../lib/stacks/dynamodb-stream-lambda-stack';

const app = new cdk.App();
new DynamodbStreamLambdaStack(app, 'DynamodbStreamLambdaStack', {
  /******
   * If you don't specify 'env', this stack will be environment-agnostic.
   * Account/Region-dependent features and context lookups will not work,
   * but a single synthesized template can be deployed anywhere.
   * For more information, see https://docs.aws.amazon.com/cdk/latest/guide/environments.html
   ***********/
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION
  }
});