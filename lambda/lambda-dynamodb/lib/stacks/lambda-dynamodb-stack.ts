import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { DynamoEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';
import { AttributeType, BillingMode, StreamViewType, Table } from 'aws-cdk-lib/aws-dynamodb';

export class LambdaDynamodbStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB table that will contain the order
    // information and enable the "stream" specification to
    // process events with a Lambda function.
    const table = new Table(this, 'order-information', {
      tableName: 'order-table',
      stream: StreamViewType.NEW_IMAGE,
      billingMode: BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'referenceId',
        type: AttributeType.STRING
      },
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Lambda Functions ********** //
    // 1. Create a Lambda function to receive the order and
    // put the item in the DynamoDB table.
    const receiveOrder = new lambda.Function(this, 'receiveOrder', {
      retryAttempts: 1,
      memorySize: 1024,
      handler: 'receiveOrder',
      functionName: 'receiveOrder',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(30),
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/receiveOrder'),
      environment: {
        "TABLE_NAME": table.tableName
      }
    });
    table.grantWriteData(receiveOrder);

    // 2. Create a Lambda function to process the received order which
    // will be triggered by an event source from the DynamoDB table.
    const processedOrder = new lambda.Function(this, 'processedOrder', {
      retryAttempts: 1,
      memorySize: 1024,
      handler: 'processedOrder',
      functionName: 'processedOrder',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/processedOrder'),
      environment: {
        "TABLE_NAME": table.tableName
      }
    });
    table.grantReadWriteData(processedOrder);
    
    // 3. Configure the event source mapping in AWS Lambda
    // to read events from the DynamoDB.
    processedOrder.addEventSource(new DynamoEventSource(table, {
      enabled: true,
      batchSize: 1,
      retryAttempts: 1,
      startingPosition: lambda.StartingPosition.LATEST
    }));
  }
}
