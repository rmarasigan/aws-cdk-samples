import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayCorsLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      retryAttempts: 2,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });

    // ********** API Gateway ********** //
    const api = new apigw.RestApi(this, 'rest-api', {
      deploy: true,
      restApiName: 'rest-api',
      apiKeySourceType: apigw.ApiKeySourceType.HEADER,
      deployOptions: {
        stageName: 'prod',
        tracingEnabled: true,
        metricsEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      },
      defaultCorsPreflightOptions: {
        allowOrigins: apigw.Cors.ALL_ORIGINS,
        allowMethods: apigw.Cors.ALL_METHODS
      }
    });

    // Setting API key
    const plan = api.addUsagePlan('api-usage-plan', {
      name: 'api-usage-plan'
    });

    // It will automatically generate an API key.
    const key = api.addApiKey('api-key', {
      apiKeyName: 'api-key'
    });

    plan.addApiKey(key);
    plan.addApiStage({ stage: api.deploymentStage });

    const integration = new apigw.LambdaIntegration(lambdaFn);
    api.root.addMethod('POST', integration, {
      apiKeyRequired: true
    });
  }
}
