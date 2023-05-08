import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as events from 'aws-cdk-lib/aws-events';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { ApiDestination } from 'aws-cdk-lib/aws-events-targets';

export class EventBridgeApiDestinationStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // ********** EventBridge Event Bus ********** //
    // 1. Create an Event Bus
    const bus = new events.EventBus(this, 'event-bus', {
      eventBusName: 'event-bus'
    });

    // 2. Create a connection and an API destination. When
    // creating an API destination, we need to specify a connection.
    const connection = new events.Connection(this, 'connection', {
      connectionName: 'connection',
      authorization: events.Authorization.apiKey('x-api-key', cdk.SecretValue.secretsManager('xxxxxxxxxxxxxxx')),
    });

    const destination = new events.ApiDestination(this, 'api-destination', {
      connection: connection,
      httpMethod: events.HttpMethod.POST,
      apiDestinationName: 'api-destination',
      endpoint: 'https://YOUR_API_ENDPOINT'
    });

    // 3. Create a rule in where to send the events, associate
    // the event bus with the said rule and add add a custom
    // event source as long as it is not starting with "aws."
    new events.Rule(this, 'event-rule', {
      eventBus: bus,
      ruleName: 'event-rule',
      eventPattern: {
        source: [ 'trigger:alarm' ]
      },
      targets: [ new ApiDestination(destination) ]
    });
    
    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will send a custom
    // event to the Event Bus to trigger an API Destination,
    // and grant the Lambda function to put an event to the
    // EventBus.
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      code: lambda.Code.fromAsset('cmd/lambdaFn'),
      environment: {
        "EVENT_BUS_NAME": bus.eventBusArn
      }
    });
    bus.grantPutEventsTo(lambdaFn);
  }
}