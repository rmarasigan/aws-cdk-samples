import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigwv2 from '@aws-cdk/aws-apigatewayv2-alpha';
import * as s3_deployment from 'aws-cdk-lib/aws-s3-deployment';
import * as apigwv2_integration from '@aws-cdk/aws-apigatewayv2-integrations-alpha';

export class S3PresignedUrlsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ********** S3 Bucket ********** //
    // 1. Create a Bucket and configure it with CORS
    // to allow access to the bucket
    const bucket = new s3.Bucket(this, `presigned-bucket-${this.region}`, {
      cors: [
        {
          allowedHeaders: [ '*' ],
          allowedOrigins: [ '*' ],
          allowedMethods: [ s3.HttpMethods.GET, s3.HttpMethods.PUT, s3.HttpMethods.HEAD ]
        }
      ],
      versioned: true,
      autoDeleteObjects: true,
      publicReadAccess: false,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      bucketName: `presigned-bucket-${this.region}`,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL
    });

    // 2. Create a Bucket that will contain the HTML index
    // of the frontend
    const app_bucket = new s3.Bucket(this, `static-web-${this.region}`, {
      bucketName: `static-web-${this.region}`,
      cors: [
        {
          allowedHeaders: [ '*' ],
          allowedOrigins: [ '*' ],
          allowedMethods: [ s3.HttpMethods.GET, s3.HttpMethods.PUT, s3.HttpMethods.HEAD ]
        }
      ],
      publicReadAccess: true,
      autoDeleteObjects: true,
      websiteIndexDocument: 'index.html',
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });

    // 3. Deploy the bucket with the frontend assets
    new s3_deployment.BucketDeployment(this, `static-web-app`, {
      destinationBucket: app_bucket,
      sources: [ s3_deployment.Source.asset('web/static') ],
    });

    // ********** Lambda Function ********** //
    // 1. Create a Lambda function that will have access to
    // the S3 Bucket
    const getPresignedURL = new lambda.Function(this, 'getPresignedURL', {
      memorySize: 1024,
      retryAttempts: 2,
      handler: 'getPresignedURL',
      functionName: 'getPresignedURL',
      runtime: lambda.Runtime.GO_1_X,
      reservedConcurrentExecutions: 1,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/getPresignedURL'),
      environment: {
        "BUCKET_NAME": bucket.bucketName
      }
    });
    bucket.grantReadWrite(getPresignedURL);

    // ********** API Gateway ********** //
    // 1. Create an HTTP API with CORS, deployment stage,
    // and routes configured
    const api = new apigwv2.HttpApi(this, 'http-api', {
      apiName: 'http-api',
      corsPreflight: {
        allowHeaders: [ '*' ],
        allowOrigins: [ '*' ],
        allowMethods: [ apigwv2.CorsHttpMethod.OPTIONS, apigwv2.CorsHttpMethod.GET, apigwv2.CorsHttpMethod.PUT ]
      }
    });

    new apigwv2.HttpStage(this, 'http-stage', {
      httpApi: api,
      stageName: 'prod',
      autoDeploy: true
    });

    api.addRoutes({
      path: '/{proxy+}',
      methods: [ apigwv2.HttpMethod.OPTIONS, apigwv2.HttpMethod.GET, apigwv2.HttpMethod.PUT ],
      integration: new apigwv2_integration.HttpLambdaIntegration('http-lambda-integration', getPresignedURL)
    });
  }
}
