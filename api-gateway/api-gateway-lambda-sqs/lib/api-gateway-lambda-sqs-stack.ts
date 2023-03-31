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
    const deadLetterQueue = new sqs.Queue(this, 'deadLetterQueue.fifo', {
      fifo: true,
      contentBasedDeduplication: true,
      queueName: 'deadLetterQueue.fifo',
    });

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

    processData.addEventSource(new eventsource.SqsEventSource(queue, {
      batchSize: 1,
      reportBatchItemFailures: true
    }));

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

    const integration = new apigw.LambdaIntegration(receiveData);
    api.root.addMethod('POST', integration);
  }
}
