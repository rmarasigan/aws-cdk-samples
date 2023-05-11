import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LogGroup } from 'aws-cdk-lib/aws-logs';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';

export class StepFunctionsApiGatewayStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** CloudWatch Logs ********** //
    // 1. Create a log group to enable the logging
    // for the State Machine.
    const logGroup = new LogGroup(this, 'sf-state-machine-logs', {
      logGroupName: 'sf-state-machine-logs',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Step Function ********** //
    // 1. Create a Step Function State Machine that
    // will be triggered by an API Gateway to start
    // the execution.
    const statemachine = new stepfunctions.StateMachine(this, 'sf-state-machine', {
      logs: {
        destination: logGroup,
        level: stepfunctions.LogLevel.ALL
      },
      tracingEnabled: true,
      stateMachineName: 'sf-state-machine',
      stateMachineType: stepfunctions.StateMachineType.EXPRESS,
      definition: stepfunctions.Chain.start(new stepfunctions.Pass(this, 'Pass'))
    });

    // ********** API Gateway ********** //
    // 1. Create a Rest API Gateway
    const api = new apigw.RestApi(this, 'rest-api', {
      deploy: true,
      restApiName: 'rest-api',
      deployOptions: {
        stageName: 'prod',
        metricsEnabled: true,
        tracingEnabled: true,
        dataTraceEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      }
    });

    // 2. Configure the API Gateway Step Function Integration
    const integration = apigw.StepFunctionsIntegration.startExecution(statemachine, {
      requestTemplates: {
        "application/json": JSON.stringify({
          "stateMachineArn": statemachine.stateMachineArn,
          "input": "$util.escapeJavaScript($input.json('$'))"
        })
      },
      integrationResponses: [
        {
          statusCode: '200',
          responseTemplates: {
            'application/json': '$input.body'
          }
        },
        {
          statusCode: '400',
          responseTemplates: {
            'application/json': '$input.body'
          }
        }
      ]
    });
    api.root.addMethod('POST', integration);
  }
}
