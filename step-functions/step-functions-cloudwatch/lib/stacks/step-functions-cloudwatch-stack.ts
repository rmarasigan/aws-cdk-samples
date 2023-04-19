import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cw_logs from 'aws-cdk-lib/aws-logs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import * as cw_log_destination from 'aws-cdk-lib/aws-logs-destinations';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';

export class StepFunctionsCloudwatchStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a table that will contain the information
    // of the filtered events
    const table = new dynamodb.Table(this, 'error-logs', {
      tableName: 'error-logs',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'LogStream',
        type: dynamodb.AttributeType.STRING
      }
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that will be the
    // destination of filtered events
    const errorHandling = new lambda.Function(this, 'errorHandling', {
      memorySize: 1024,
      handler: 'errorHandling',
      functionName: 'errorHandling',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/errorHandling')
    });

    // ********** CloudWatch ********** //
    // 1. Reference an existing log group
    const logGroup = cw_logs.LogGroup.fromLogGroupName(this, 'cw-log-group', 'LOG_GROUP_NAME');

    // 2. Create a subscription filter
    new cw_logs.SubscriptionFilter(this, 'log_subscription', {
      logGroup: logGroup,
      destination: new cw_log_destination.LambdaDestination(errorHandling),
      filterPattern: cw_logs.FilterPattern.stringValue('$.ERROR_KEY_FIELD', '=', 'ERROR_FIELD_VALUE')
    });

    // ********** Step Function ********** //
    // 1. Create a task for DynamoDB PutItem operation
    const putErrorAlertTask = new stepfunctions_tasks.DynamoPutItem(this, 'ErrorLogPutItem', {
      table: table,
      item: {
        'ID': stepfunctions_tasks.DynamoAttributeValue.fromString(stepfunctions.JsonPath.stringAt('$.logEvents[0].id')),
        'Owner': stepfunctions_tasks.DynamoAttributeValue.fromString(stepfunctions.JsonPath.stringAt('$.owner')),
        'LogGroup': stepfunctions_tasks.DynamoAttributeValue.fromString(stepfunctions.JsonPath.stringAt('$.logGroup')),
        'LogStream': stepfunctions_tasks.DynamoAttributeValue.fromString(stepfunctions.JsonPath.stringAt('$.logStream')),
        'Message': stepfunctions_tasks.DynamoAttributeValue.fromString(stepfunctions.JsonPath.stringAt('$.logEvents[0].message')),
        'Timestamp': stepfunctions_tasks.DynamoAttributeValue.numberFromString(stepfunctions.JsonPath.jsonToString(stepfunctions.JsonPath.stringAt('$.logEvents[0].timestamp'))) 
      }
    });

    // 2. Create a definition for the State Machine
    const definition = stepfunctions.Chain.start(putErrorAlertTask);

    // 3. Create a Step Function State Machine, grant the Lambda Function
    // to start the execution and grant the State Machine to put an item
    // to the DynamoDB table
    const statemachine = new stepfunctions.StateMachine(this, 'ErrorLogStateMachine', {
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'ErrorLogStateMachine',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      stateMachineType: stepfunctions.StateMachineType.EXPRESS
    });
    table.grantWriteData(statemachine);
    statemachine.grantStartExecution(errorHandling);
    errorHandling.addEnvironment("STATE_MACHINE_ARN", statemachine.stateMachineArn);
  }
}
