import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as event_source from 'aws-cdk-lib/aws-lambda-event-sources';

export class S3LambdaDynamodbStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    // 1. Create an S3 Bucket
    const bucket = new s3.Bucket(this, `item-data-${this.region}`, {
      publicReadAccess: false,
      bucketName: `item-data-${this.region}`,
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL
    });

    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB Table that will contain
    // the processed information
    const table = new dynamodb.Table(this, 'data-table', {
      tableName: 'data-table',
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'id',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'name',
        type: dynamodb.AttributeType.STRING
      }
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will process
    // each object-created event from the S3 Bucket
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
        "TABLE_NAME": table.tableName
      }
    });
    table.grantWriteData(lambdaFn);
    bucket.grantReadWrite(lambdaFn);

    lambdaFn.addEventSource(new event_source.S3EventSource(bucket, {
      filters: [{ suffix: '.json' }],
      events: [ s3.EventType.OBJECT_CREATED ]
    }));
  }
}
