import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class ApiGatewayCognitoLambdaStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** Cognito ********** //
    // 1. Create a Cognito User Pool to set up custom scopes
    // and to use access tokens to authorize API method calls.
    const user_pool = new cognito.UserPool(this, 'cognito-api-user-pool', {
      userPoolName: 'cognito-api-user-pool',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // 2. Create a Resource Server Scope. When creating a
    // resource server, you must provide a resource server
    // name and a resource server identifier.
    const readOnlyScope = new cognito.ResourceServerScope({
      scopeName: 'read-only',
      scopeDescription: 'Read-only access for users'
    });

    const resourceServer = user_pool.addResourceServer('users-resource-server', {
      identifier: 'users',
      scopes: [ readOnlyScope ],
      userPoolResourceServerName: 'users-resource-server'
    });

    // 3. Create a User Pool Client to obtain an identity
    // or access token to be included to call API methods.
    user_pool.addClient('cognito-api-user-pool-client', {
      oAuth: {
        flows: {
          clientCredentials: true,
        },
        scopes: [ cognito.OAuthScope.resourceServer(resourceServer, readOnlyScope) ]
      },
      generateSecret: true,
      userPoolClientName: 'cognito-api-user-pool-client',
      accessTokenValidity: cdk.Duration.minutes(30)
    });

    // 4. Create a User Pool Domain
    user_pool.addDomain('cognito-api-domain', {
      cognitoDomain: {
        domainPrefix: 'learningjourney'
      }
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda Function
    const lambdaFn = new lambda.Function(this, 'lambdaFn', {
      memorySize: 1024,
      handler: 'lambdaFn',
      functionName: 'lambdaFn',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/lambdaFn')
    });
    lambdaFn.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // ********** API Gateway ********** //
    // 1. Create a Rest API, configure the integration as
    // Lambda Integration and use the Cognito User Pools
    // as the authorizer of the API.
    const api = new apigw.RestApi(this, 'rest-api', {
      deploy: true,
      restApiName: 'rest-api',
      deployOptions: {
        stageName: 'prod',
        metricsEnabled: true,
        tracingEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.ERROR
      }
    });
    api.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 2. Create and configure a Cognito User Pool Authorizer.
    const authorizer = new apigw.CognitoUserPoolsAuthorizer(this, 'users-authorizer', {
      cognitoUserPools: [ user_pool ],
      authorizerName: 'users-authorizer',
      identitySource: apigw.IdentitySource.header('Authorization')
    });
    authorizer.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);

    // 3. Create a Lambda integration and add a GET request method.
    const integration = new apigw.LambdaIntegration(lambdaFn);

    api.root.addMethod('GET', integration, {
      authorizer: authorizer,
      authorizationType: apigw.AuthorizationType.COGNITO,
      authorizationScopes: [ `${resourceServer.userPoolResourceServerId}/${readOnlyScope.scopeName}`]
    });
  }
}
