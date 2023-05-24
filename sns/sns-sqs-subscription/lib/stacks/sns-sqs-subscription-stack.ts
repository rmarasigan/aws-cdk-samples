import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { SqsEventSource} from 'aws-cdk-lib/aws-lambda-event-sources';
import { SqsSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';

export class SnsSqsSubscriptionStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SQS Queue ********** //
    // 1. Create a deadletter queue that will contain the
    // unsuccessfully processed.
    const deadLetterQueue = new Queue(this, 'deadLetterQueue', {
      queueName: 'deadLetterQueue',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // 2. Create a queue that is configured with deadletter queue.
    const queue = new Queue(this, 'order-queue', {
      deadLetterQueue: {
        maxReceiveCount: 5,
        queue: deadLetterQueue
      },
      queueName: 'order-queue',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      visibilityTimeout: cdk.Duration.seconds(60)
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that will be triggered by
    // the SQS Queue.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      tracing: lambda.Tracing.ACTIVE,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });
    lambdaFn.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Configure the Lambda function's SQS event source.
    lambdaFn.addEventSource(new SqsEventSource(queue, {
      batchSize: 1
    }));

    // ********** SNS ********** //
    // 1. Create an SNS Topic.
    const topic = new Topic(this, 'sns-topic', {
      topicName: 'sns-topic',
      displayName: 'SQS Subscription'
    });
    topic.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Subscribe an SQS Queue to the SNS Topic.
    topic.addSubscription(new SqsSubscription(queue, {
      deadLetterQueue: deadLetterQueue
    }));
  }
}