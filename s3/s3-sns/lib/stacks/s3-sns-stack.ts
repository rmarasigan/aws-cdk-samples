import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as s3_notification from 'aws-cdk-lib/aws-s3-notifications';
import * as sns_subscription from 'aws-cdk-lib/aws-sns-subscriptions';

export class S3SnsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    // 1. Create an S3 Bucket
    const bucket = new s3.Bucket(this, `order-data-${this.region}`, {
      publicReadAccess: false,
      bucketName: `order-data-${this.region}`,
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL
    });

    // ********** SNS ********** //
    // 1. Create an SNS Topic and add an object removal
    // event notification to the Bucket and send it
    // via email notification
    const topic = new sns.Topic(this, 's3-bucket-events', {
      topicName: 's3-bucket-events',
      displayName: 'S3 Bucket Object Events'
    });

    // 2. Set SNS as the Event Destination
    bucket.addEventNotification(s3.EventType.OBJECT_REMOVED, new s3_notification.SnsDestination(topic));

    // 3. Subscribe an email address to SNS topic
    // The email subscription require confirmation by visiting the link sent to the email address.
    topic.addSubscription(new sns_subscription.EmailSubscription('your_email@gmail.com'));
  }
}
