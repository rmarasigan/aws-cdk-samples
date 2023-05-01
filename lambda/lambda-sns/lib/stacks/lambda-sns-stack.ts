import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import { Role, ServicePrincipal, PolicyStatement, Effect } from 'aws-cdk-lib/aws-iam';

export class LambdaSnsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // ********** IAM Role ********** //
    // 1. Create an IAM Role for the Lambda function.
    const role = new Role(this, 'lambda-sns-role', {
      roleName: 'lambda-sns-role',
      assumedBy: new ServicePrincipal('lambda.amazonaws.com')
    });

    // 2. Add a policy for the Lambda function to be
    // able to subscribe to the SNS topic.
    role.addToPolicy(new PolicyStatement({
      resources: [ '*' ],
      effect: Effect.ALLOW,
      actions: [ 'sns:Subscribe' ]
    }));

    // ********** SNS ********** //
    // 1. Create an SNS Topic.
    const topic = new Topic(this, 'alert-topic', {
      topicName: 'alert-topic',
      displayName: 'Alert Topic'
    });

    // 2. Subscribe an email address to the SNS Topic.
    const subscription = new EmailSubscription('j.doe@email.com');

    // 3. Subscribe an email address to the SNS topic.
    // The email subscription require confirmation by visiting the link sent to the email address.
    topic.addSubscription(subscription);

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will send a
    // confirmation message to an SNS Topic and grant
    // the Lambda function to publish a message.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      role: role,
      memorySize: 1024,
      retryAttempts: 1,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      tracing: lambda.Tracing.ACTIVE,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
      environment: {
        "TOPIC_ARN": topic.topicArn
      }
    });
    topic.grantPublish(lambdaFn);
  }
}
