import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as eventbridge from 'aws-cdk-lib/aws-events';
import * as event_target from 'aws-cdk-lib/aws-events-targets';

export class S3EventbridgeLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    new s3.Bucket(this, `order-data-${this.region}`, {
      publicReadAccess: false,
      eventBridgeEnabled: true,
      bucketName: `order-data-${this.region}`,
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
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });

    // ********** EventBridge Rule ********** //
    const s3BucketEventRule = new eventbridge.Rule(this, 's3-bucket-event-rule', {
      enabled: true,
      ruleName: 's3-bucket-event-rule',
      eventPattern: {
        source: [ 'aws.s3' ],
        account: [ this.account ],
        region: [ this.region ]
      }
    });
    s3BucketEventRule.addTarget(new event_target.LambdaFunction(lambdaFn));
  }
}
