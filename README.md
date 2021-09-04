# Loan Broker application

Inspired by [this blog post](https://www.enterpriseintegrationpatterns.com/ramblings/loanbroker_stepfunctions.html).

WIP

## Learnings

- _AWS SAM_ allows you to deploy multiple nested stacks at once. To do so, use the `AWS::Serverless::Application` resource type.

  - One thing that is **suboptimal** is the **lack of visibility in the errors produced by a specific nested stack**.
    All you get in terms of the errors when you deploy is information that the "Embedded stack XXX was not successfully updated".
    To diagnose problem you need to check the console. You could also use the [describe-stack-events CLI command](https://docs.aws.amazon.com/cli/latest/reference/cloudformation/describe-stack-events.html#examples).

- The API Gateway REST API is not automatically deploying your changes. Once you have created the deployment, unless the _logical ID_ of that deployment changes, your changes will not be live.

  - You can either pre-process the template manually and add some kind of identifier at the end of the deployment resource name.
  - You can use _CloudFormation Macros_.

- **To use macros with _nested stacks_ you have to specify the `CAPABILITY_AUTO_EXPAND` CFN capability**.

  - You would also have to do that if you were interacting with the CFN API directly. According to the [AWS documentation](https://docs.aws.amazon.com/AWSCloudFormation/latest/APIReference/API_UpdateStack.html)

    > If your stack template contains one or more macros, and you choose to update a stack directly from the processed template, without first reviewing the resulting changes in a change set, you must acknowledge this capability.

  - The _AWS CDK_ and _AWS SAM_ are deploying your infrastructure going through the _changeset_ flow. The _changeset_ is created and then executed.

- By using the _CloudFormation Macro_ you can either

  1. Transform your whole template by specifying the macro name within the `Transform` array.

     ```yaml
     Transform:
       - FirstMacro
       - SecondMacro
     Resources:
       # Resources
     ```

  2. Transform the `Properties` section of a given resource by using the `Fn::Transform` _intrinsic function_ and specifying the macro name.

     ```yaml
     MyResource:
       Type: MyResourceType
       Properties:
         # ... Properties
         Fn::Transform:
           Name: MyMacroName
           Parameters: {}
     ```

- **By default, the _Template_ section within CFN console is displaying the CFN template BEFORE the macro was applied**.
  **To view the processed template, toggle the 'View processed template' switch on that tab**.

- APIGW validation is weird

  - You can create a model, associate it with a method, but the request body will not by validated against that model
  - To actually enable request validation you have to create the `AWS::ApiGateway::RequestValidator` resource and also associate it with a given method.
  - I'm not 100% sure why you have to attach two resources to a given method to enable validation.

- Testing _Step Functions_ is hard.

  - You could use the [Step Functions Local](https://docs.aws.amazon.com/step-functions/latest/dg/sfn-local.html).

    1. You are not able to test permissions of the SFN.
    2. I would argue that **local SFN emulator is really good for testing simple data transformations** and not the service integrations.

  - You could test the whole flow end-to-end but, usually, such endeavors end with the test taking a lot of time.

    1. You can test the permissions of the SFN.

- Use the `ResultPath: null` if you do not want to pollute the step output with the response of the service your step integrates with.
  Very useful for _DynamoDB_ integrations.

- The are **multiple ways to invoke lambda with SFNs**.

  1. Using the `"Resource": "arn:aws:states:::lambda:invoke"` and specifying `FunctionName` within the `Parameters` block.

     - You can use the `.waitForTaskToken`.
     - You have to explicitly specify the payload your function will receive via the `Payload` property.

  2. Using the `"Resource": "FUNCTION_ARN"` way.
     - You **cannot use the `Parameters` block**
     - The **SFN input is _implicitly_ passed to the function payload**.
