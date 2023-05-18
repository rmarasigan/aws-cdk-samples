import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as cw_logs from 'aws-cdk-lib/aws-logs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import { SqsEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';

export class StepFunctionsSqsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SQS Queue ********** //
    // 1. Creates a deadletter queue that will contain the
    // unsuccessfully processed, duplicates are not tolerated,
    // and should have a ".fifo" to the queue name.
    const deadLetterQueue = new sqs.Queue(this, 'deadLetterQueue.fifo', {
      queueName: 'deadLetterQueue.fifo',
      contentBasedDeduplication: true
    });

    // 2. Create a queue that is configured to be a FIFO queue with
    // deadletter queue. It is needed to add a ".fifo" to the queue
    // name. The timeout should be greater or equal to the Lambda
    // Functions execution timeout.
    const queue = new sqs.Queue(this, 'transaction-queue.fifo', {
      fifo: true,
      deadLetterQueue: {
        maxReceiveCount: 5,
        queue: deadLetterQueue
      },
      queueName: 'transaction-queue.fifo',
      contentBasedDeduplication: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      visibilityTimeout: cdk.Duration.seconds(60)
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function to process the received
    // SQS Event data.
    const transaction = new lambda.Function(this, 'transaction', {
      memorySize: 1024,
      handler: 'transaction',
      functionName: 'transaction',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/transaction')
    });
    transaction.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Configure the Lambda Function to receive an event
    // from SQS Event and to process 1 record at a time.
    transaction.addEventSource(new SqsEventSource(queue, {
      batchSize: 1
    }));

    // ********** CloudWatch ********** //
    // 1. Create a log group that is to be used by state machine
    const logGroup = new cw_logs.LogGroup(this, 'sqs-send-state-machine', {
      logGroupName: 'sqs-send-state-machine',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Step Function ********** //
    // 1. Create a definition to send messages to the configured
    // SQS queue.
    const definition = new stepfunctions_tasks.SqsSendMessage(this, 'transaction-definition', {
      queue: queue,
      messageBody: stepfunctions.TaskInput.fromObject({
        'id': stepfunctions.JsonPath.stringAt('$.id'),
        'order_line_id': stepfunctions.JsonPath.stringAt('$.line_id'),
        'total': stepfunctions.JsonPath.numberAt('$.total'),
        'change': stepfunctions.JsonPath.numberAt('$.change'),
        'customer': stepfunctions.JsonPath.objectAt('$.customer')
      }),
      messageGroupId: 'process.transaction',
      integrationPattern: stepfunctions.IntegrationPattern.REQUEST_RESPONSE
    });

    // 2. Create a Step Function State Machine to send a
    // message(s) to the configured SQS queue.
    new stepfunctions.StateMachine(this, 'SQSSendMessageStateMachine', {
      logs: {
        level: stepfunctions.LogLevel.ALL,
        destination: logGroup,
        includeExecutionData: true
      },
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'SQSSendMessageStateMachine',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      stateMachineType: stepfunctions.StateMachineType.EXPRESS
    });
  }
}