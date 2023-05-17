import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cw_logs from 'aws-cdk-lib/aws-logs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';

export class StepFunctionsDynamodbStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** DynamoDB Table ********** //
    // 1. Create a table that will contain the information
    // of the users.
    const table = new dynamodb.Table(this, 'users-table', {
      tableName: 'users-table',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: 'ID',
        type: dynamodb.AttributeType.STRING
      }
    });

     // ********** CloudWatch ********** //
    // 1. Create a log group that is to be used by State Machine.
    const logGroup = new cw_logs.LogGroup(this, 'user-state-machine', {
      logGroupName: 'user-state-machine',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // ********** Step Function ********** //
    let attributevalue = stepfunctions_tasks.DynamoAttributeValue;

    // 1. Create a task for the DynamoDB PutItem Operation to
    // create a new item containing user information.
    const insertUser = new stepfunctions_tasks.DynamoPutItem(this, 'insert-user', {
      table: table,
      item: {
        "ID":  attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.id')),
        "Username": attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.username')),
        "Password": attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.password')),
        "FirstName": attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.first_name')),
        "LastName": attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.last_name')),
        "Role": attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.role'))
      }
    });

    // 2. Create a task for the DynamoDB UpdateItem Operations to update
    // an existing item's attributes.
    const updateUsername = new stepfunctions_tasks.DynamoUpdateItem(this, 'update-username', {
      table: table,
      key: {
        'ID': attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.id'))
      },
      conditionExpression: 'attribute_exists(ID)',
      expressionAttributeValues: {
        ':username': attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.username'))
      },
      updateExpression: 'SET Username = :username'
    });

    const updatePassword = new stepfunctions_tasks.DynamoUpdateItem(this, 'update-password', {
      table: table,
      key: {
        'ID': attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.id'))
      },
      conditionExpression: 'attribute_exists(ID)',
      expressionAttributeValues: {
        ':password': attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.password'))
      },
      updateExpression: 'SET Password = :password'
    });

    // 3. Create a task for the DynamoDB DeleteItem Operation to delete a
    // single item in a table by using the primary key.
    const deleteUser = new stepfunctions_tasks.DynamoDeleteItem(this, 'delete-user', {
      table: table,
      key: {
        'ID': attributevalue.fromString(stepfunctions.JsonPath.stringAt('$.id'))
      },
      resultPath: stepfunctions.JsonPath.DISCARD
    });

    // 4. Create a definition that has a "Choice" state that can take a
    // different path through the workflow based on the '$.transaction_type'
    // input value.
    const choice = new stepfunctions.Choice(this, 'transaction-type');
    
    const definition = choice.when(stepfunctions.Condition.stringEquals('$.transaction_type', 'delete'), deleteUser)
                             .when(stepfunctions.Condition.stringEquals('$.transaction_type', 'insert'), insertUser)
                             .when(stepfunctions.Condition.stringEquals('$.transaction_type', 'update:username'), updateUsername)
                             .when(stepfunctions.Condition.stringEquals('$.transaction_type', 'update:password'), updatePassword);
    
    // 5. Create a Step Function State Machine and grant the State
    // Machine a permission of read and write to the DynamoDB Table.
    const statemachine = new stepfunctions.StateMachine(this, 'UserStateMacine', {
      logs: {
        level: stepfunctions.LogLevel.ALL,
        destination: logGroup,
        includeExecutionData: true
      },
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'UserStateMachine',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      stateMachineType: stepfunctions.StateMachineType.EXPRESS
    });
    table.grantReadWriteData(statemachine);          
  }
}
