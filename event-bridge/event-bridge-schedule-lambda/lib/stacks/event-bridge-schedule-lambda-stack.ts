import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as events from 'aws-cdk-lib/aws-events';
import { LambdaFunction } from 'aws-cdk-lib/aws-events-targets';

export class EventBridgeScheduleLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will be triggered
    // based on the EventBridge scheduled rule.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });

    // ********** EventBridge Schedule Rule ********** //
    // 1. Create an EventBridge schedule rule that will trigger
    // the Lambda Function every 1 minute.
    const scheduleRule = new events.Rule(this, 'schedule-rule', {
      ruleName: 'schedule-rule',
      enabled: true,
      targets: [ new LambdaFunction(lambdaFn) ],

      // Run every 1 minute
      // 1. Using "rate"
      schedule: events.Schedule.rate(cdk.Duration.minutes(1))

      // Run every 1 minute between 8:00 AM and 9:00 PM (PH Time)
      // 2. Using "cron"
      // schedule: events.Schedule.cron({ minute: '0/1', hour: '2-21' })
      
      // 3. Using "expression"
      // schedule: events.Schedule.expression('cron(0/1 2-21 ? * * *)')
    });
  }
}