import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';

export class ApiGatewayLambdaDynamodbStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB Table that will contain
    // the information of Coffee.
    const table = new dynamodb.Table(this, 'data-table', {
      tableName: 'data-table',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'key',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'name',
        type: dynamodb.AttributeType.STRING
      }
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will receive
    // the API Gateway event record and save the processed
    // request to a DynamoDB table. Grant the Lambda Function
    // a write access to the DynamoDB table.
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

    // ********** API Gateway ********** //
    // 1. Create a Rest API and configure the integration
    // as LambdaIntegration.
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

    // 2. Create a Lambda Integration and add a POST request method.
    const integration = new apigw.LambdaIntegration(lambdaFn);
    api.root.addMethod('POST', integration);
  }
}
