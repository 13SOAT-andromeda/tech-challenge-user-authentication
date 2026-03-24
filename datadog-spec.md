The lambda-datadog Terraform module wraps the aws_lambda_function resource and automatically configures your Lambda function for Datadog Serverless Monitoring by:

Adding the Datadog Lambda layers
Redirecting the Lambda handler
Enabling the collection and sending of metrics, traces, and logs to Datadog

```terraform
module "lambda-datadog" {
  source  = "DataDog/lambda-datadog/aws"
  version = "4.0.0"

  environment_variables = {
    "DD_API_KEY_SECRET_ARN" : "<DATADOG_API_KEY_SECRET_ARN>"
    "DD_ENV" : "<ENVIRONMENT>"
    "DD_SERVICE" : "<SERVICE_NAME>"
    "DD_SITE": "<DATADOG_SITE>"
    "DD_VERSION" : "<VERSION>"
  }

  datadog_extension_layer_version = 93

  # aws_lambda_function arguments
}
```

Replace the aws_lambda_function resource with the lambda-datadog Terraform module. Then, specify the source and version of the module.

Set the aws_lambda_function arguments:

All of the arguments available in the aws_lambda_function resource are available in this Terraform module. Arguments defined as blocks in the aws_lambda_function resource are redefined as variables with their nested arguments.

For example, in aws_lambda_function, environment is defined as a block with a variables argument. In the lambda-datadog Terraform module, the value for the environment_variables is passed to the environment.variables argument in aws_lambda_function. See inputs for a complete list of variables in this module.

Fill in the environment variable placeholders:

Replace <DATADOG_API_KEY_SECRET_ARN> with the ARN of the AWS secret where your Datadog API key is securely stored. The key needs to be stored as a plaintext string (not a JSON blob). The secretsmanager:GetSecretValue permission is required. For quick testing, you can instead use the environment variable DD_API_KEY and set your Datadog API key in plaintext.
Replace <ENVIRONMENT> with the Lambda function’s environment, such as prod or staging
Replace <SERVICE_NAME> with the name of the Lambda function’s service
Replace <DATADOG_SITE> with datadoghq.com. (Ensure the correct Datadog site is selected on this page).
Replace <VERSION> with the version number of the Lambda function
Select the version of the Datadog Extension Lambda layer to use. If left blank the latest layer version will be used.

datadog_extension_layer_version = 93
