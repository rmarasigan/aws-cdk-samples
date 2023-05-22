import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';

export class SnsEmailSubscriptionStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SNS ********** //
    // 1. Create an SNS Topic.
    const topic = new Topic(this, 'sns-topic', {
      topicName: 'sns-topic',
      displayName: 'Transaction'
    });
    topic.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Subscribe an email address to the SNS Topic.
    const subscription = new EmailSubscription('j.doe@email.com');

    // 3. Email subscription confirmation.
    // The email subscription require confirmation by visiting the link sent to the
    // email address.
    topic.addSubscription(subscription);
  }
}
