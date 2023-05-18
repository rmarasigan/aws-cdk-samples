#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { ApiGatewayCognitoLambdaStack } from '../lib/stacks/api-gateway-cognito-lambda-stack';

const app = new cdk.App();
new ApiGatewayCognitoLambdaStack(app, 'ApiGatewayCognitoLambdaStack', {
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