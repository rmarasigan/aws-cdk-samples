import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ses from 'aws-cdk-lib/aws-ses';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Role, Effect, ServicePrincipal, PolicyStatement } from 'aws-cdk-lib/aws-iam';

export class SesSendEmailStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** SES ********** //
    // 1. Create and verify the identity. It can be
    // an email address or a domain. In here, we are
    // going to use an Email Identity.
    const identity = new ses.EmailIdentity(this, 'ses-email-identity', {
      identity: ses.Identity.email('abc@email.com'),
      feedbackForwarding: true,
    });
    identity.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // ********** IAM Policy ********** //
    // 1. Create a Role Policy where the lambda can send
    // emails to the Email Identities of SES.
    const role = new Role(this, 'lambda-ses-role', {
      roleName: 'lambda-ses-role',
      assumedBy: new ServicePrincipal('lambda.amazonaws.com')
    });

    role.addToPolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      actions: [ 'ses:SendEmail' ],
      resources: [ `arn:aws:ses:${this.region}:${this.account}:identity/*` ]
    }));

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function to send an email message
    // using the email identity as the sender.
    new lambda.Function(this, 'sendEmail', {
      role: role,
      memorySize: 1024,
      handler: 'sendEmail',
      functionName: 'sendEmail',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/sendEmail'),
      environment: {
        "EMAIL_IDENTITY": identity.emailIdentityName
      }
    });
  }
}
