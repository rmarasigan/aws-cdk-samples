import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ses from 'aws-cdk-lib/aws-ses';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as stepfunctions from 'aws-cdk-lib/aws-stepfunctions';
import * as stepfunctions_tasks from 'aws-cdk-lib/aws-stepfunctions-tasks';
import { Role, Effect, ServicePrincipal, PolicyStatement } from 'aws-cdk-lib/aws-iam';

export class StepFunctionsSesStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SES ********** //
    // 1. Create and verify the identity. It can be an email address
    // or a domain. In here, we are going to use an Email Identity.
    const identity = new ses.EmailIdentity(this, 'ses-email-identity', {
      identity: ses.Identity.email('abc@email.com')
    });
    identity.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // ********** IAM Policy ********** //
    // 1. Create a Role Policy where the lambda can send emails to
    // the Email Identities of SES.
    const role = new Role(this, 'lambda-ses-role', {
      roleName: 'lambda-ses-role',
      assumedBy: new ServicePrincipal('lambda.amazonaws.com')
    });
    role.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // This will allow us to use any Email Identity as it is using
    // the "*" wildcard.
    role.addToPolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      actions: [ 'ses:SendEmail' ],
      resources: [ `arn:aws:ses:${this.region}:${this.account}:identity/*` ]
    }));

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function to send an email message using the
    // email identity as the sender.
    const sendEmail = new lambda.Function(this, 'sendEmail', {
      role: role,
      memorySize: 1024,
      handler: 'sendEmail',
      functionName: 'sendEmail',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/sendEmail'),
      environment: {
        "EMAIL_IDENTITY": identity.emailIdentityName
      }
    });

    // 2. Create a Lambda Function to process the transaction
    // and start the Step Function State Machine.
    const transaction = new lambda.Function(this, 'transaction', {
      memorySize: 1024,
      handler: 'transaction',
      functionName: 'transaction',
      runtime: lambda.Runtime.GO_1_X,
      tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset('cmd/transaction')
    });

    // ********** Step Function ********** //
    // 1. Create a task for Lambda Invocation
    const task = new stepfunctions_tasks.LambdaInvoke(this, 'step-function-task', {
      inputPath: '$',
      lambdaFunction: sendEmail,
      retryOnServiceExceptions: false
    });

    // 2. Create a definition for the State Machine
    const definition = stepfunctions.Chain.start(task);

    // 3. Create a Step Function State Machine, grant the Lambda
    // Function to start the execution and add the ARN of statemachine
    // as environment variable to the 'transaction' Lambda Function.
    const statemachine = new stepfunctions.StateMachine(this, 'step-function-state-machine', {
      definition: definition,
      tracingEnabled: true,
      stateMachineName: 'step-function-state-machine',
      stateMachineType: stepfunctions.StateMachineType.EXPRESS
    });
    statemachine.grantStartExecution(transaction);
    statemachine.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);
    transaction.addEnvironment('STATE_MACHINE_ARN', statemachine.stateMachineArn);
  }
}
