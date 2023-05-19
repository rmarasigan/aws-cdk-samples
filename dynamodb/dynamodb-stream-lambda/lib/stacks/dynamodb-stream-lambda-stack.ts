import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { DynamoEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';
import { Table, BillingMode, AttributeType, StreamViewType } from 'aws-cdk-lib/aws-dynamodb';

export class DynamodbStreamLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB Table that will contain an order
    // information and enable the "stream" specification to
    // process events with a Lambda function.
    const table = new Table(this, 'orders', {
      tableName: 'orders',
      stream: StreamViewType.NEW_IMAGE,
      billingMode: BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'referenceId',
        type: AttributeType.STRING
      },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      contributorInsightsEnabled: true
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that will be triggered by
    // an event from the DynamoDB table.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      tracing: lambda.Tracing.ACTIVE,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
    });
    lambdaFn.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Configure the event source mapping in AWS Lambda to
    // read events from the DynamoDB Stream.
    lambdaFn.addEventSource(new DynamoEventSource(table, {
      enabled: true,
      batchSize: 1,
      startingPosition: lambda.StartingPosition.LATEST
    }));
  }
}
