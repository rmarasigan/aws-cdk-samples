# What is Amazon SES?

![ses](../assets/img/ses.png)

**Amazon SES** is a cloud-based email service for sending both transactional and mass emails. It lets you send transactional emails, marketing messages, or any other type of high-quality content your customers. It is an email platform that provides an easy, cost-effective way for you to send and receive email using your own email addresses and domains.

## Why Amazon SES?
### Deliverability Rate
Deliverability rate is one of the main parameters to consider when choosing an email sending services. Amazon takes reputation and whitelisting seriously by supporting all three authentication (DKIM, SPF, and DMARC). In addition you can track your sending activity and manage your reputation.

### Personalization
Content personalization with replacement tags.

### E-mail Reception
With Amazon SES, you can not only send emails but also retrieve them. In this case, you have a set of flexible options, as well as usage of the received message as a trigger in AWS Lambda.

### Pricing
For all other cases, Amazon’s policy of "pay for what you actually use" applies. For emails over the 62,000 limit, you will pay as much as $0.1 for every 1,000 emails you send + $0.12 for each GB of attachments. If you need a dedicated IP address for more security, it will cost and additional $24.95 per month. 

## Use Cases
### Automate transactional messages
Keep your customers up to date by sending automated emails, such as purchase or shipping notifications, order status updates, and policy change notices.

### Deliver marketing emails globally
Tell customers around the world about products and services through newsletters, special offers, and engaging content.

### Send timely notifications to customers
Send customers timely notifications about their interaction with your products and services, including daily reminders, weekly usage reports, and newsletters.

### Send bulk email communications
Deliver messages — including notifications and announcements — to large groups, and track results using configuration sets.

## Pros and Cons
### Pros
* High deliverability and reliability along with high sending rates
* No need for additional maintenance once you have all set up
* Best quality to price ratio
* A comprehensive set of tools for both email receiving and further management

### Cons
* A compllicated initial configuration
* Initial limitations before you get approved and verify your sending domains
* Amazon SES is a simple sending servicec, not a marketing platform
* Amazon SES does not provide you with an email list storage

## Setting Up
### Verify your email address or domain
You have to verify the identities that you plan to send email from. In Amazon SES, an identity can be an email address or an entire domain. When you verify a domain, you can use Amazon SES to send email from any address on that domain.

Remember that email addresses are case sensitive. Another important thing in using AWS is its Region. They are physical locations of Amazon data centers available to any customer. These region don't limit the usage of AWS services for you, as their purpose is distributing workloads. In addition, email verification is connected to the region.

#### Verify an email address identity
1. Check the inbox of the email address used to create your identity and look for an email from `no-reply-aws@amazon.com`.
2. Open the email and click the link to complete the verification process for the email address. After it is complete, the **Identity status** updates to **Verified**

### Production access
When you first start using Amazon SES, your account is in a ***sandbox environment***. While your account is in the sandbox, you can only send email to addresses that you have verified. This means that you can send:
* Emails to verified email addresses and domains or to the Amazon SES mailbox simulator
* Up to 200 messages in 24 hours
* Only 1 message per second

### Integration
At this step, you should define exactly how you would like to send your emails.

#### Using SMTP
With SMTP you can:
* Send an email from your application if you are using a framework that supports the SMTP authentication
* Send messages from software packages you already use
* Send emails right from your email client
* Integrate it with a server where you host your app

### Using API
API integration requires definite technical skills and can be used to:
* Make raw query requests and responses
* Use an AWS SDK
* Use a Command Line Interface (CLI)

## Reference
* [Amazon SES](https://aws.amazon.com/ses/)
* [Amazon SES pricing](https://aws.amazon.com/ses/pricing/)
* [What is Amazon SES?](https://docs.aws.amazon.com/ses/latest/dg/Welcome.html)
* [Verifying an email address identity](https://docs.aws.amazon.com/ses/latest/dg/creating-identities.html#just-verify-email-proc)
* [What Is Amazon SES and How to Use It?](https://mailtrap.io/blog/amazon-ses-explained/)
* [Using the Amazon SES API to send email](https://docs.aws.amazon.com/ses/latest/dg/send-email-api.html)
* [Migrating to Amazon SES from another email-sending solution](https://docs.aws.amazon.com/ses/latest/dg/send-email-getting-started-migrate.html)
* [Amazon SES Tutorial | How To Send Emails Using AWS SES | AWS Training | Edureka](https://www.youtube.com/watch?v=gVRTKuMFc0c)