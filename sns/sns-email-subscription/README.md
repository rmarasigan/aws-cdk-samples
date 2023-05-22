# SNS Email Subscription

![SNS Email Subscription](assets/img/sns-email-subscription.png)

Publishing a message to the Amazon Simple Notification Service via AWS Console or AWS CLI to send an e-mail to the endpoint.

**NOTE**: The e-mail subscription require confirmation by visiting the link sent to the e-mail address.

### SNS Subscription Confirmation

![SNS Subscription Confirmation](assets/img/sns-subscription-confirmation.png)

### Invoking SNS function via AWS Console
1. Go to Amazon SNS → Topics → *Your SNS Topic* → **Publish message**
2. Enter the following information

    a. Subject (optional)

    b. Message body (see [sample payload](#sample-payload))

3. Click on the **Publish message** in the bottom right corner

### Invoking Lambda function via AWS CLI
1. Use the following command and replace the placeholder `sns_topic_arn` with the actual SNS Topic ARN

    ```bash
    aws sns publish \
    --topic-arn sns_topic_arn \
    --message file://sns-sample-message.txt
    --subject "Cash on Delivery (COD) #210327LL6J2NE7"
    ```

    It will publish the specified message to the specified SNS Topic. The message comes from a text file, which enables you to include line breaks.

### Sample Payload
```
Hello John,
  
We regret to inform you that your Cash on Delivery (COD) payment request for order #210327LL6J2NE7 has been declined. We have notified the seller to cancel the shipping of your item(s).
```

### Sample SNS Published Email

![SNS Published Email](assets/img/sns-published-email.png)

### AWS CDK API / Developer Reference
* [Amazon Simple Notification Service](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_sns-readme.html)
* [Amazon Simple Notification Service Subscriptions](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_sns_subscriptions-readme.html)

### AWS Documentation Developer Guide
* [Publish](https://docs.aws.amazon.com/sns/latest/api/API_Publish.html)
* [Subscribe](https://docs.aws.amazon.com/sns/latest/api/API_Subscribe.html)
* [Amazon SNS FAQs](https://aws.amazon.com/sns/faqs/)
* [Email notifications](https://docs.aws.amazon.com/sns/latest/dg/sns-email-notifications.html)
* [Amazon Simple Notification Service endpoints and quotas](https://docs.aws.amazon.com/general/latest/gr/sns.html)

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
    dev@dev:~:aws-cdk-samples/sns/sns-email-subscription$ make init

    # With "profile" parameter
    dev@dev:~:aws-cdk-samples/sns/sns-email-subscription$ make init profile=[profile_name]
    ```

2. Deploy the project.
    ```bash
    # Without passing "profile" parameter
    dev@dev:~:aws-cdk-samples/sns/sns-email-subscription$ make deploy

    # With "profile" parameter
    dev@dev:~:aws-cdk-samples/sns/sns-email-subscription$ make deploy profile=[profile_name]
    ```