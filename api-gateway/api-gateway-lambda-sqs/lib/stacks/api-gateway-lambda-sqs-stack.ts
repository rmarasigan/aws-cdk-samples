import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as eventsource from 'aws-cdk-lib/aws-lambda-event-sources';

export class ApiGatewayLambdaSqsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SQS Queue ********** //
    // 1. Create a deadletter queue that will contain the unsuccessfully
    // processed and should have a ".fifo" to the queue name.
    const deadLetterQueue = new sqs.Queue(this, 'deadLetterQueue.fifo', {
      fifo: true,
      contentBasedDeduplication: true,
      queueName: 'deadLetterQueue.fifo',
    });

    // 2. Create a queue that is configured to be a
    // FIFO queue with deadletter queue. It is needed
    // to add a ".fifo" to the queue name.
    const queue = new sqs.Queue(this, 'item-queue.fifo', {
      fifo: true,
      deadLetterQueue: {
        maxReceiveCount: 5,
        queue: deadLetterQueue
      },
      queueName: 'item-queue.fifo',
      contentBasedDeduplication: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      visibilityTimeout: cdk.Duration.seconds(60)
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function to send the message
    // to an SQS queue and grant access to send messages.
    const receiveData = new lambda.Function(this, 'receiveData', {
      memorySize: 1024,
      retryAttempts: 2,
      handler: 'receiveData',
      functionName: 'receiveData',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/receiveData'),
      environment: {
        "QUEUE_URL": queue.queueUrl
      }
    });
    queue.grantSendMessages(receiveData);

    // 2. Create a Lambda function that will be triggered
    // for every event received from the SQS queue and consume
    // the messages.
    const processData = new lambda.Function(this, 'processData', {
      memorySize: 1024,
      retryAttempts: 2,
      handler: 'processData',
      functionName: 'processData',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/processData')
    });

    // 3. Configure the Lambda function's SQS event source.
    processData.addEventSource(new eventsource.SqsEventSource(queue, {
      batchSize: 1,
      reportBatchItemFailures: true
    }));

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
    const integration = new apigw.LambdaIntegration(receiveData);
    api.root.addMethod('POST', integration);
  }
}
