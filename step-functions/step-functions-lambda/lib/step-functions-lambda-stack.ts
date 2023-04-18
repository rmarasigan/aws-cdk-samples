import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';

export class StepFunctionsLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    const transaction = new lambda.Function(this, 'transaction', {
      memorySize: 1024,
      handler: 'transaction',
      functionName: 'transaction',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/transaction')
    });

    const processTransaction = new lambda.Function(this, 'processTransaction', {
      memorySize: 1024,
      handler: 'processTransaction',
      functionName: 'processTransaction',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/processTransaction')
    });

    // ********** Step Function ********** //
    const processTransactionTask = new stepfunctions_tasks.LambdaInvoke(this, 'processTransactionTask', {
      inputPath: '$',
      lambdaFunction: processTransaction,
      retryOnServiceExceptions: false
    });

    const definition = stepfunctions.Chain.start(processTransactionTask);
    const statemachine = new stepfunctions.StateMachine(this, 'TransactionStateMachine', {
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'TransactionStateMachine',
      stateMachineType: stepfunctions.StateMachineType.STANDARD
    });
    statemachine.grantStartExecution(transaction);
    transaction.addEnvironment("STATE_MACHINE_ARN", statemachine.stateMachineArn);

    // ********** API Gateway ********** //
    const api = new apigw.RestApi(this, 'rest-api', {
      deploy: true,
      restApiName: 'rest-api',
      deployOptions: {
        stageName: 'prod',
        tracingEnabled: true,
        metricsEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      }
    });

    const integration = new apigw.LambdaIntegration(transaction);
    api.root.addMethod('POST', integration);
  }
}
