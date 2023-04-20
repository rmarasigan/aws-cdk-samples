import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cw_logs from 'aws-cdk-lib/aws-logs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as cw_log_destination from 'aws-cdk-lib/aws-logs-destinations';

export class CloudwatchSubscriptionLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function that will be the
    // destination of filtered events
    const handleSubscription = new lambda.Function(this, 'handleSubscription', {
      memorySize: 1024,
      handler: 'handleSubscription',
      functionName: 'handleSubscription',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/handleSubscription')
    });

    // ********** CloudWatch ********** //
    // 1. Reference an existing log group
    const logGroup = cw_logs.LogGroup.fromLogGroupName(this, 'cw-log-group', 'LOG_GROUP_NAME');

    // 2. Create a subscription filter
    new cw_logs.SubscriptionFilter(this, 'log_subscription', {
      logGroup: logGroup,
      destination: new cw_log_destination.LambdaDestination(handleSubscription),
      filterPattern: cw_logs.FilterPattern.stringValue('$.KEY_FIELD', '=', 'VALUE_FIELD')
    });
  }
}
