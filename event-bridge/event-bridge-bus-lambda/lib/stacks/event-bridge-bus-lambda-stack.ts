import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as events from 'aws-cdk-lib/aws-events';
import { LambdaFunction, EventBus } from 'aws-cdk-lib/aws-events-targets';

export class EventBridgeBusLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    
    // ********** EventBridge Event Bus ********** //
    // 1. Create an event bus
    const bus = new events.EventBus(this, 'transaction-event-bus', {
      eventBusName: 'transaction-event-bus'
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will send a custom event
    // to the EventBridge Event Bus and grant a permission to send events.
    const transactionFn = new lambda.Function(this, 'transactionFn', {
      memorySize: 1024,
      handler: 'transactionFn',
      functionName: 'transactionFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/transactionFn'),
      environment: {
        "EVENT_BUS_NAME": bus.eventBusName
      }
    });
    bus.grantPutEventsTo(transactionFn);

    // 2. Create a Lambda functions that will be triggered
    // based on the EventBridge Event Rule received event.
    const paymentFn = new lambda.Function(this, 'paymentFn', {
      memorySize: 1024,
      handler: 'paymentFn',
      functionName: 'paymentFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/paymentFn')
    });

    const cancelFn = new lambda.Function(this, 'cancelFn', {
      memorySize: 1024,
      handler: 'cancelFn',
      functionName: 'cancelFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/cancelFn')
    });

    // ********** EventBridge Rule ********** //
    // 1. Create a rule in where to send the transaction events,
    // associate the event bus with the said rule and add a custom
    // event source as long as it is not starting with "aws.".
    new events.Rule(this, 'transaction-payment-rule', {
      eventBus: bus,
      ruleName: 'transaction-payment-rule',
      eventPattern: {
        source: [ 'transaction:payment' ]
      },
      targets: [ new LambdaFunction(paymentFn) ]
    });

    new events.Rule(this, 'transaction-cancel-rule', {
      eventBus: bus,
      ruleName: 'transaction-cancel-rule',
      eventPattern: {
        source: [ 'transaction:cancel' ]
      },
      targets: [ new LambdaFunction(cancelFn) ]
    });
  }
}
