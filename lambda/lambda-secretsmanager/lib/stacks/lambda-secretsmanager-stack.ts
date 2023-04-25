import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Secret } from 'aws-cdk-lib/aws-secretsmanager';

export class LambdaSecretsmanagerStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Secrets Manager ********** //
    // 1. Create a Secret Manager that will contain
    // the user credentials.
    const secret = new Secret(this, 'user-credentials', {
      description: 'A sample secret manager that will contain the user credentials',
      secretName: 'user-credentials',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Lambda Functions ********** //
    // 1. Create a Lambda function to interact with
    // the AWS Secrets Manager and allow a read and write
    // permission to the Lambda function.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      retryAttempts: 1,
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
      environment: {
        "SECRET_ARN": secret.secretArn
      }
    });
    secret.grantRead(lambdaFn);
    secret.grantWrite(lambdaFn);
  }
}
