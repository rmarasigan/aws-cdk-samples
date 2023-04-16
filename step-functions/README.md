# What is AWS Step Functions?

**AWS Step Functions** is a serverless orchestration service that lets you integrate with AWS Lambda functions and other AWS services to build business-critical applications. Through Step Functions' graphical console, you see your application’s workflow as a series of event-driven steps. At each step of a given workflow, Step Functions manages input, output, error handling, and retries, so that developers can focus on higher-value business logic for their applications.

Step Functions is based on state machines and tasks. A ***state machine*** is a workflow. A ***task*** is a state in a workflow that represents a single unit of work that another AWS service performs. Each step in a workflow is a state.

### State Machine
In computer science, a **state machine** is defined as a type of computational device that is able to store various status values and update them based on inputs. AWS Step Functions builds upon this very concept and uses the term state machine to refer to an application workflow. Developers can build a state machine in Step Functions with JSON files by using the [**Amazon States Language**](https://docs.aws.amazon.com/step-functions/latest/dg/concepts-amazon-states-language.html).

You can choose a *standard workflow* for processes that are long-running or that require human intervention. *Express workflows* are well-suited for short-running (fewer than five minutes), high-volume processes.

### State
A **state** represents a step in your workflow. States can perform a variety of functions:
* Perform work in the state machine (*Task state*)
* Choose between different paths in a workflow (*Choice state*)
* Stop the workflow with failure or success (a *Fail* or *Succeed state*)
* Pass output or some fixed data to another state (*Pass state*)
* Pause the workflow for a specified amount of time (*Wait state*)
* Begin parallel branches of execution (*Parallel state*)
* Repeat execution for each item of input (*Map state*)

### Task State
A **task state** (typically just referred to as a task) within your state machine is used to complete a single unit of work. Tasks can be used to call the API actions of over two hundred Amazon and AWS services.

#### Activity tasks
**Activity tasks** let you connect a step in your workflow to a batch of code that is running elsewhere. This external batch of code, called an activity worker, polls Step Functions for work, asynchronously completes the work using your code, and returns results. Activity tasks are common with asynchronous workflows in which some human intervention is required (to verify a user account, for example).

#### Service tasks
**Service tasks** let you connect steps in your workflow to specific AWS services. Step Functions sends requests to other services, waits for the task to complete, and then continues to the next step in the workflow. They can be used easily for automated steps, such as executing a Lambda function.

The states that you decide to include in your state machine and the relationships between your states form the core of your Step Functions workflow.

## Use Cases of Step Functions

### Automate extract, transform, and load (ETL) processes
Ensure that multiple long-running ETL jobs run in order and complete successfully, without the need for manual orchestration.

* Automate steps of an ETL process
* Build a data processing pipeline for streaming data

### Automate security and IT functions
Create automated workflows, including manual approval steps, for security incident response.

* Respond to operational events in your AWS account
* Synchronize data between source and destination S3 buckets
* Orchestrate a security incident response for IAM policy creation

### Orchestrate microservices
Combine multiple AWS Lambda functions into responsive serverless applications and microservices.

* Combine Lambda functions to build a web-based application
* Invoke a business process in response to an event using Express Workflows
* Combine Lambda functions to build a web-based application - with a human approval

### Orchestrate largescale parallel workloads
Iterate over and process large data-sets such as security logs, transaction data, or image and video files.

* Large scale data processing
* Extract data from PDF or images for processing
* Run an ETL pipeline with multiple jobs in parallel
* Split and transcode video using massive parallelization

## Reference
* [AWS Step Functions Overview](https://www.datadoghq.com/knowledge-center/aws-step-functions/)
* [AWS Step Functions Use Cases](https://aws.amazon.com/step-functions/use-cases/)
* [AWS Step Functions Documentation](https://docs.aws.amazon.com/step-functions/?icmpid=docs_homepage_serverless)
* [What are AWS Step Functions? (and why you should love them)](https://www.youtube.com/watch?v=zCIpWFYDJ8s)