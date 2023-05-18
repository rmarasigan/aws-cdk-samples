# API Gateway + Cognito → Lambda

![API Gateway with Cognito to Lambda](assets/img/api-gateway-cognito-lambda.png)

Invoking the Lambda Function with AWS API Gateway only if it is authorized to do so. The API Gateway is configured with an authorizer using a Cognito User Pool to control who can access the REST API. The client must first provide a client id and client secret to the token endpoint in exchange for access tokens. The access token is used to authorize the API calls based on the custom scopes. Once the REST API receives the request, it will pass the access token and scopes to AWS Cognito for verification and confirm if the tokens and scopes are valid. The Lambda Function will be triggered and will return a successful response back to the client.

### Configuring an Amazon Cognito User Pool for your API
#### Amazon Cognito
1. Create a Cognito User Pool
2. Create Resource Server and custom scopes in the user pool
3. Create an App Client with Client Credentials Flow
4. Create a Domain in the user pool

#### API Gateway
1. Create an Amazon REST API Gateway
2. Create an Authorizer in the REST API with the Cognito User Pool created in Amazon Cognito
3. Create a Lambda Integration
4. Create resources and methods in the said resources with Cognito Authorizer and OAuth Scopes created in Amazon Cognito User Pool

### [Token Endpoint](https://www.rfc-editor.org/rfc/rfc6749#section-3.2)
When you add a **domain** to the Cognito User Pool, Amazon Cognito activates an OAuth 2.0 token endpoint that is dedicated to the User Pool. The **Cognito domain** will be used as the token endpoint where the client app sends a request in exchange for the access tokens.

### Getting the access token
1. Get the Amazon Cognito Domain.

    a. Amazon Cognito → User pools → *Your Cognito User Pool* → App integration → Cognito domain

2. Add "`/oauth2/token`" to the Cognito domain endpoint and it only supports the `HTTP POST` method.

    ```
    https://domain_name.auth.us-east-1.amazoncognito.com/oauth2/token
    ```

3. The `Content-Type` must always be `application/x-www-form-urlencoded`.

4. Add an `Authorization` header that contains the client's secret. The authorization header is a Base64Encoded (`client_id:client_secret`).

    ```
    Authorization: Basic xxxxxxxxxxxxxxxxxxxxxx
    ```

5. Add the following to the request body:

    a. `grant_type`: `client_credentials`

    b. `scope`: your_cognito_custom_scope

### Call the REST API Gateway
1. Get the `access_token` from the response of the `/oauth2/token` endpoint.

2. Add a header of `Authorization` that contains the `access_token`.

    ```
    Authorization: access_token_value
    ```

### Sample Response
```json
{
    "content": "Hello from the Lambda function!"
}
```

### AWS CDK API / Developer Reference
* [AWS Lambda](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_lambda-readme.html)
* [Amazon Cognito](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_cognito-readme.html)
* [Amazon API Gateway](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_apigateway-readme.html)
* [OAuth 2.0 Client Credentials Flow with AWS Cognito in AWS CDK](http://buraktas.com/oauth-client-credentials-flow-aws-cdk/)
* [Amazon Cognito and API Gateway Based Machine to Machine Authorization](https://github.com/aws-samples/amazon-cognito-and-api-gateway-based-machine-to-machine-authorization-using-aws-cdk)

### AWS SDK v2 API / Developer Reference
* [Getting Started with the AWS SDK for Go V2](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/)

### AWS Documentation Developer Guide
* [Token endpoint](https://docs.aws.amazon.com/cognito/latest/developerguide/token-endpoint.html)
* [Client Authentication](https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication)
* [Using tokens with user pools](https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html)
* [Using AWS Lambda with Amazon API Gateway](https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway.html)
* [Defining resource servers for your user pool](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-define-resource-servers.html)
* [Create an Amazon Cognito user pool for a REST API](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-create-cognito-user-pool.html)
* [Control access to a REST API using Amazon Cognito user pools as authorizer](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-integrate-with-cognito.html)

### Useful commands
The `cdk.json` file tells the CDK Toolkit how to execute your app.

* `npm install`     install projects dependencies
* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `cdk deploy`      deploy this stack to your default AWS account/region
* `cdk diff`        compare deployed stack with current state
* `cdk synth`       emits the synthesized CloudFormation template
* `cdk bootstrap`   deployment of AWS CloudFormation template to a specific AWS environment (account and region)
* `cdk destroy`     destroy this stack from your default AWS account/region

## Deploy

### Using `make` command
1. Install all the dependencies, bootstrap your project, and synthesized CloudFormation template.
    ```bash
    # Without passing "profile" parameter
    dev@dev:~:aws-cdk-samples/api-gateway/api-gateway-cognito-lambda$ make init

    # With "profile" parameter
    dev@dev:~:aws-cdk-samples/api-gateway/api-gateway-cognito-lambda$ make init profile=[profile_name]
    ```

2. Deploy the project.
    ```bash
    # Without passing "profile" parameter
    dev@dev:~:aws-cdk-samples/api-gateway/api-gateway-cognito-lambda$ make deploy

    # With "profile" parameter
    dev@dev:~:aws-cdk-samples/api-gateway/api-gateway-cognito-lambda$ make deploy profile=[profile_name]
    ```