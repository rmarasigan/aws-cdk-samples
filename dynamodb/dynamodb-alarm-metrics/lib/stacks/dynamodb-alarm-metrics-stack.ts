import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import * as cw from 'aws-cdk-lib/aws-cloudwatch';
import { SnsAction } from 'aws-cdk-lib/aws-cloudwatch-actions';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import { Table, Operation, BillingMode, AttributeType } from 'aws-cdk-lib/aws-dynamodb';

export class DynamodbAlarmMetricsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a DynamoDB Table that will contain
    // the order information and monitor the read,
    // write, and throttled requests using CloudWatch.
    const table = new Table(this, 'order', {
      tableName: 'orders',
      billingMode: BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'referenceId',
        type: AttributeType.STRING
      },
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // Sum of Consumed Read Capacity per minute
    const metricConsumedReadCapacity = table.metricConsumedReadCapacityUnits({
      region: this.region,
      account: this.account,
      statistic: cw.Stats.SUM,
      period: cdk.Duration.seconds(60),
      label: `[sum: ${cw.Stats.SUM}] ${table.tableName}_read`,
    });

    // Sum of Consumed Write Capacity per minute
    const metricsConsumedWriteCapacity = table.metricConsumedWriteCapacityUnits({
      region: this.region,
      account: this.account,
      statistic: cw.Stats.SUM,
      period: cdk.Duration.seconds(60),
      label: `[sum: ${cw.Stats.SUM}] ${table.tableName}_write`
    });

    // Sum of Throttled Reuqests per minute for Scan, Query,
    // GetItem, PutItem, UpdateItem and DeleteItem operations.
    const metricsThrottledRequests = table.metricThrottledRequestsForOperations({
      region: this.region,
      account: this.account,
      statistic: cw.Stats.SUM,
      period: cdk.Duration.seconds(60),
      label: `[sum: ${cw.Stats.SUM}] ${table.tableName}_throttle`,
      operations: [ Operation.SCAN, Operation.QUERY, Operation.GET_ITEM, Operation.PUT_ITEM, Operation.UPDATE_ITEM, Operation.DELETE_ITEM ]
    });

    // ********** SNS ********** //
    // 1. Create an SNS topic.
    const topic = new Topic(this, 'sns-topic', {
      topicName: 'sns-topic',
      displayName: 'Order DynamoDB Alarm'
    });
    topic.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Subscribe an email address to the SNS Topic.
    // The email subscription require confirmation by visiting
    // the link sent to the email address.
    const subscription = new EmailSubscription('j.deo@email.com');
    topic.addSubscription(subscription);

    // ********** CloudWatch Alarms ********** //
    // 1. Create an alarm that will send an Amazon SNS Email when the
    // alarm changes state.
    // It will ALARM when the metric if greater than the threshold
    // configured.
    const readAlarm = new cw.Alarm(this, 'OrderDynamoDBReadAlarm', {
      metric: metricConsumedReadCapacity,
      threshold: 10,
      actionsEnabled: true,
      evaluationPeriods: 1,
      alarmName: 'OrderDynamoDBReadAlarm',
      alarmDescription: 'Alarm when read capacity is greater than the threshold',
      comparisonOperator: cw.ComparisonOperator.GREATER_THAN_THRESHOLD
    });
    readAlarm.addAlarmAction(new SnsAction(topic));
    readAlarm.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // It will ALARM when the metric if greater than the threshold
    // configured.
    const writeAlarm = new cw.Alarm(this, 'OrderDynamoDBWriteAlarm', {
      metric: metricsConsumedWriteCapacity,
      threshold: 10,
      actionsEnabled: true,
      evaluationPeriods: 1,
      alarmName: 'OrderDynamoDBWriteAlarm',
      alarmDescription: 'Alarm when write capacity is greater than the threshold',
      comparisonOperator: cw.ComparisonOperator.GREATER_THAN_THRESHOLD
    });
    writeAlarm.addAlarmAction(new SnsAction(topic));
    writeAlarm.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // It will ALARM when the metric if greater than the threshold
    // for the throttled requests across all operations configured.
    const throttleAlarm = new cw.Alarm(this, 'OrderDynamoDBThrottleAlarm', {
      metric: metricsThrottledRequests,
      threshold: 5,
      actionsEnabled: true,
      evaluationPeriods: 1,
      alarmName: 'OrderDynamoDBThrottleAlarm',
      alarmDescription: 'Alarm when throttle requests is greater than the threshold',
      comparisonOperator: cw.ComparisonOperator.GREATER_THAN_THRESHOLD
    });
    throttleAlarm.addAlarmAction(new SnsAction(topic));
    throttleAlarm.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);
  }
}
