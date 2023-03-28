import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayAsyncLambdaStack extends cdk.Stack {
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
      deployOptions: {
        stageName: 'prod',
        tracingEnabled: true,
        metricsEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      }
    });

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
