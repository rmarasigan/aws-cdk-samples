import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayLambdaS3Stack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    const bucket = new s3.Bucket(this, `api-data-${this.region}`, {
      bucketName: `api-data-${this.region}`,
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL
    });

    // ********** Lambda Function ********** //
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      retryAttempts: 2,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
      environment: {
        "BUCKET_NAME": bucket.bucketName
      }
    });
    bucket.grantWrite(lambdaFn);

    // ********** API Gateway ********** //
    const api = new apigw.RestApi(this, 'rest-api', {
      deploy: true,
      restApiName: 'rest-api',
      deployOptions: {
        stageName: 'prod',
        tracingEnabled: true,
        metricsEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      }
    });

    const integration = new apigw.LambdaIntegration(lambdaFn);
    api.root.addMethod('POST', integration);
  }
}
