# Loan Broker application

Inspired by [this blog post](https://www.enterpriseintegrationpatterns.com/ramblings/loanbroker_stepfunctions.html).

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

- While the SFN use _JSONPath_ for path traversal, not all of its functions are available to you.
  Having said that, you should **take a look at the [_Intrinsic functions_ reference](https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-intrinsic-functions.html)**. You can do some cool stuff with them, **especially the `States.Format` one**

  1. Imagine wanting to overload _DynamoDB_ primary key, how one might add a `TYPE#` prefix?
     Here is how: `"S.$": "States.Format('LOAN#{}', $$.Execution.Name)"`

  2. Maybe you want to create an array from object values?
     Here is how: `"banks.$": "States.Array('Static', $.Resolved, 'AnotherStatic')"`

- Remember that, when specifying _DynamoDB_ conditions, the _DynamoDB_ will first identify an item to compare against, then run the `Condition Expression`.

  1. Do we pay for such identification? Is that baked into the price of the service?

- You **cannot use the `ConditionExpression` while performing `GetItem`**.
  Maybe it's because for the `ConditionExpression` to work, DDB has to inspect the item in the first place - you might as well do it yourself and check?

- You **cannot natively marshal back the result of the `GetItem` call**.
  What you can do is to use **`ResultSelector` and dig into each individual property** to retrieve the value.
  One note though, **if you decide to use the `ResultSelector`, remember that every underlying value will be a string**. The _DynamoDB_ uses the `N` or `S` to define what kind of type the value is but keeps the value as a string.

- Actually, the `ResultPath` is very handy. I imagine if you want to get the actual data returned by the service integration, you will be using it a lot.

- _JSONPath_ can be a bit cumbersome to work with. Some expression return an array of values, that would not be the problem if not for the fact that
  **I could not find a way to access an specific index within an array via _JSONPath_**.

- There is an **alternative syntax** you can use to create the `ExpressionAttributeValues` when using _DynamoDB UpdateItem_ operation.
  Usually, I used a simple assignment syntax like so: `":value": "VALUE"` but **this will not work when the "VALUE" is dynamic, e.g coming from the SFN input**.
  An alternative syntax is to **use the _attribute value_ notation** like so:

  ```text
    ":quotes": {
      "S.$": "States.JsonToString($.quotes)"
    }
  ```

  This syntax **also works for non-dynamic values**

  ```text
    ":newStatus": {
      "S": "FINISHED"
    },
  ```

- Time and time again I'm amazed by how hard is it to get an absolute path to a file using Go.
  There is no `__dirname`, you have to walk the directory tree recursively.

- Sadly, _AWS SAM_ does not have the functionality to automatically save _CloudFormation_ outputs to a file.
  I was able to circumvent it by some _Makefile_ magic. Actually it was not that bad and learned things along the way, but it would be nice to have this functionality be there out of the box.
