import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { LambdaSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';

export class SnsLambdaSubscriptionStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    //  1. Create a Lambda Function that will be triggered by
    // the SNS topic.
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

    // ********** SNS ********** //
    // 1. Create an SNS Topic.
    const topic = new Topic(this, 'sns-topic', {
      topicName: 'sns-topic',
      displayName: 'Lambda Function Subscription'
    });
    topic.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Subscribe a Lambda Function to the SNS Topic.
    topic.addSubscription(new LambdaSubscription(lambdaFn));
  }
}
