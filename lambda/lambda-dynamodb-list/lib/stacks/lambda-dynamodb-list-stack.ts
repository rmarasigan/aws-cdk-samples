import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Table, BillingMode ,AttributeType } from 'aws-cdk-lib/aws-dynamodb';

export class LambdaDynamodbListStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    
    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB table that will contain the order
    // information and its order line information.
    const table = new Table(this, 'order-information', {
      tableName: 'order-information',
      billingMode: BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'referenceId',
        type: AttributeType.STRING
      },
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that will create an order,
    // insert/update the item and grant a read/write
    // permission in the DynamoDB table.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
      environment: {
        "TABLE_NAME": table.tableName
      }
    });
    table.grantReadWriteData(lambdaFn);
  }
}
