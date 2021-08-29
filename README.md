# Loan Broker application

Inspired by [this blog post](https://www.enterpriseintegrationpatterns.com/ramblings/loanbroker_stepfunctions.html).

## Learnings

- _AWS SAM_ allows you to deploy multiple nested stacks at once. To do so, use the `AWS::Serverless::Application` resource type.

- The API Gateway REST API is not automatically deploying your changes. Once you have created the deployment, unless the _logical ID_ of that deployment changes, your changes will not be live.

  - You can either pre-process the template manually and add some kind of identifier at the end of the deployment resource name.
  - You can use _CloudFormation Macros_.

- **To use macros with _nested stacks_ you have to specify the `CAPABILITY_AUTO_EXPAND` CFN capability**

  - You would also have to do that if you were interacting with the CFN API directly without nested stacks and without creating a changeset.
  - How does the SAM and CDK publish to CFN ? Is there a synthetic changeset approval process?

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
