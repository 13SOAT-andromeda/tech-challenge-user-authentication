terraform {
  required_version = ">= 1.0.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    key     = "lambda-user-validation.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_ecr_repository" "this" {
  name = "tech-challenge-user-authentication-repo"
}

data "aws_iam_role" "lab_role" {
  name = "LabRole"
}

locals {
  function_name = "tech-challenge-user-authentication"
}

module "lambda-datadog" {
  source  = "DataDog/lambda-datadog/aws"
  version = "4.0.0"

  environment_variables = {
    JWT_SECRET          = var.jwt_secret
    JWT_REFRESH_SECRET  = var.jwt_refresh_secret
    DB_HOST             = var.db_host
    DB_USER             = var.db_user
    DB_PASSWORD         = var.db_password
    DB_NAME             = var.db_name
    DB_PORT             = var.db_port
    DB_SSLMODE          = var.db_sslmode
    DYNAMODB_TABLE_NAME = var.dynamodb_table_name
    DYNAMODB_ENDPOINT   = var.dynamodb_endpoint
    DD_API_KEY          = var.dd_key
    DD_ENV              = "production"
    DD_SERVICE          = local.function_name
    DD_SITE             = "us5.datadoghq.com"
    DD_VERSION          = "1.0.0"
  }

  datadog_extension_layer_version = 93

  function_name = local.function_name
  role          = data.aws_iam_role.lab_role.arn
  package_type  = "Image"
  image_uri     = "${data.aws_ecr_repository.this.repository_url}:${var.image_tag}"

  reserved_concurrent_executions = 3

  timeout     = 30
  memory_size = 128

  image_config_command = ["bootstrap"]
}

resource "aws_lambda_function_url" "this" {
  function_name      = local.function_name
  authorization_type = "NONE"

  depends_on = [module.lambda-datadog]
}

resource "aws_lambda_permission" "allow_public_url" {
  statement_id           = "AllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = local.function_name
  principal              = "*"
  function_url_auth_type = "NONE"

  depends_on = [module.lambda-datadog]
}

