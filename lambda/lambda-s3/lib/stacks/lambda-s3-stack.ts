import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Bucket, BlockPublicAccess } from 'aws-cdk-lib/aws-s3';

export class LambdaS3Stack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    // 1. Create an S3 Bucket that will contain the
    // order details.
    const bucket = new Bucket(this, `order-data-${this.region}`, {
      bucketName: `order-data-${this.region}`,
      autoDeleteObjects: true,
      publicReadAccess: false,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function to interact with
    // the AWS S3 Bucket and grant a read and write
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
        "BUCKET_NAME": bucket.bucketName
      }
    });
    bucket.grantReadWrite(lambdaFn);
  }
}