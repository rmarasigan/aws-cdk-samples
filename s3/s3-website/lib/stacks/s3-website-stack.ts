import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3_deployment from 'aws-cdk-lib/aws-s3-deployment';

export class S3WebsiteStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // 1. Create a Bucket that will contain the HTML index
    // of the frontend
    const bucket = new s3.Bucket(this, `static-web-${this.region}`, {
      bucketName: `static-web-${this.region}`,
      publicReadAccess: true,
      websiteIndexDocument: 'index.html'
    });

    // 2. Deploy the bucket with the frontend assets
    new s3_deployment.BucketDeployment(this, `static-web-app`, {
      destinationBucket: bucket,
      sources: [ s3_deployment.Source.asset('web/static') ],
    });
  }
}