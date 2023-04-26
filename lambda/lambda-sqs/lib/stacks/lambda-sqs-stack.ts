import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { SqsEventSource} from 'aws-cdk-lib/aws-lambda-event-sources';

export class LambdaSqsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SQS Queue ********** //
    // 1. Create a deadletter queue that will
    // contain the unsuccessfully processed and
    // should have a ".fifo" to the queue name.
    const deadLetterQueue = new Queue(this, 'deadLetterQueue.fifo', {
      queueName: 'deadLetterQueue.fifo',
      contentBasedDeduplication: true
    });

    // 2. Create a queue that is configured to be a
    // FIFO queue with deadletter queue. It is needed
    // to add a ".fifo" to the queue name.
    const queue = new Queue(this, 'order-queue.fifo', {
      fifo: true,
      deadLetterQueue: {
        maxReceiveCount: 5,
        queue: deadLetterQueue
      },
      queueName: 'order-queue.fifo',
      contentBasedDeduplication: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      visibilityTimeout: cdk.Duration.seconds(60)
    });

    // ********** Lambda Functions ********** //
    // 1. Create a Lambda function to send the message
    // to an SQS queue and grant access to send messages.
    const receiveOrder = new lambda.Function(this, 'receiveOrder', {
      retryAttempts: 1,
      memorySize: 1024,
      handler: 'receiveOrder',
      functionName: 'receiveOrder',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/receiveOrder'),
      environment: {
        "QUEUE_URL": queue.queueUrl
      }
    });
    queue.grantSendMessages(receiveOrder);

    // 2. Create a Lambda function that will be triggered
    // for every event received from the SQS queue and consume
    // the messages.
    const processOrder = new lambda.Function(this, 'processOrder', {
      retryAttempts: 1,
      memorySize: 1024,
      handler: 'processOrder',
      functionName: 'processOrder',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/processOrder')
    });

    // 3. Configure the Lambda function's SQS event source.
    processOrder.addEventSource(new SqsEventSource(queue, {
      batchSize: 1,
      reportBatchItemFailures: true
    }));
  }
}
