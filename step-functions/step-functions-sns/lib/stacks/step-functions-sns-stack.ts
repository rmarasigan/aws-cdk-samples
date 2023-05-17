import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import * as cw_logs from 'aws-cdk-lib/aws-logs';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';

export class StepFunctionsSnsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SNS ********** //
    // 1. Create an SNS Topic and subscribe an email
    // address to it.
    const topic = new Topic(this, 'step-function-topic', {
      topicName: 'step-function-topic',
      displayName: 'Step Function Topic Event'
    });
    topic.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Subscribe an email address to the SNS topic.
    // The email subscription require confirmation by visiting the link
    // sent to the email address.
    topic.addSubscription(new EmailSubscription('your_email@email.com'));

    // ********** CloudWatch ********** //
    // 1 Create a log group that is to be used by State Machine.
    const logGroup = new cw_logs.LogGroup(this, 'sns-publish-state-machine', {
      logGroupName: 'sns-publish-state-machine',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Step Function ********** //
    // 1. Create a definition to publish messages to the
    // configured SNS topic.
    const definition = new stepfunctions_tasks.SnsPublish(this, 'step-function-definition', {
      topic: topic,
      subject: stepfunctions.JsonPath.stringAt('$.subject'),
      message: stepfunctions.TaskInput.fromJsonPathAt('$.content'),
      resultPath: stepfunctions.JsonPath.DISCARD,
      integrationPattern: stepfunctions.IntegrationPattern.REQUEST_RESPONSE
    });

    // 2. Create a Step Function State Machine to publish a 
    // message(s) to the SNS topic.
    new stepfunctions.StateMachine(this, 'SNSPublishStateMachine', {
      logs: {
        level: stepfunctions.LogLevel.ALL,
        destination: logGroup,
        includeExecutionData: true
      },
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'SNSPublishStateMachine',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      stateMachineType: stepfunctions.StateMachineType.EXPRESS
    });
  }
}
