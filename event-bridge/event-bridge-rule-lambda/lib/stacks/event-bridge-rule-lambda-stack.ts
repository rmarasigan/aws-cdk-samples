import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Bucket } from 'aws-cdk-lib/aws-s3';
import * as events from 'aws-cdk-lib/aws-events';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as cloudtrail from 'aws-cdk-lib/aws-cloudtrail'
import { LambdaFunction } from 'aws-cdk-lib/aws-events-targets';

export class EventBridgeRuleLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    // 1. Create an instance of an existing S3 Bucket.
    const sourceBucket = Bucket.fromBucketArn(this, 'your-bucket-name-us-east-1', 'arn:aws:s3:::your-bucket-name-us-east-1');

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that wil be triggered
    // for every matched pattern.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(60),      
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });

    // ********** CloudTrail ********** //
    // 1. Create a CloudTrail Trail and add an
    // S3 Event selector that has a WRITE operation.
    const trail = new cloudtrail.Trail(this, 'capture-events-trail', {
      trailName: 'capture-events-trail'
    });
    trail.addS3EventSelector([{ bucket: sourceBucket }], { readWriteType: cloudtrail.ReadWriteType.WRITE_ONLY });

    // ********** EventBridge Rule ********** //
    // 1. Create a rule that will match a pattern where
    // an object is being uploaded to the S3 bucket and
    // add the Lambda Function as the event target.
    new events.Rule(this, 's3-upload-rule', {
      enabled: true,
      ruleName: 's3-upload-rule',
      eventPattern: {
        source: [ 'aws.s3' ],
        detail: {
          eventName: [ "PutObject" ],
          requestParameters: {
            bucketName: [ sourceBucket.bucketName ]
          }
        }
      },
      targets: [ new LambdaFunction(lambdaFn) ]
    });
  }
}