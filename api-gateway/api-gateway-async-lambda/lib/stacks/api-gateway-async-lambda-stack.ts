import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayAsyncLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will receive
    // the API Gateway event record.
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
    // 1. Create a Rest API and configure the
    // integration as LambdaIntegration.
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

    // 2. Create a Lambda Integration and disable the proxy
    // option, and add an integration request parameter header
    // mapping "X-Amz-Invocation-Type: 'Event'" to make it asynchronous.
    const integration = new apigw.LambdaIntegration(lambdaFn, {
      proxy: false,
      requestParameters: {
        'integration.request.header.X-Amz-Invocation-Type': "'Event'"
      },
      requestTemplates: {
        "application/json": JSON.stringify({
          body: "$util.escapeJavaScript($input.body)"
        })
      },
      integrationResponses: [
       {
        statusCode: '200',
        responseTemplates: {
          'application/json': '$input.path("$.body")'
        }
       } 
      ]
    });

    // 3. Add a POST request method and a method response of
    // 200 OK HTTP Status.
    const restAPI = api.root.addMethod('POST', integration, {
      methodResponses: [{
        statusCode: '200',
        responseParameters: {
          'method.response.header.X-Amz-Invocation-Type': true
        }
      }]
    });
  }
}
