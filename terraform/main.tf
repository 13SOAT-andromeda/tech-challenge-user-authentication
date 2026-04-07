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

data "aws_vpc" "vpc" {
  filter {
    name   = "tag:Name"
    values = ["vpc"]
  }
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.vpc.id]
  }

  filter {
    name   = "tag:Name"
    values = ["*private*"]
  }
}

resource "aws_security_group" "lambda" {
  name        = "tech-challenge-user-auth-lambda-sg"
  description = "Security group for user auth lambda"
  vpc_id      = data.aws_vpc.vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "tech-challenge-user-auth-lambda-sg"
  }
}

resource "aws_lambda_function" "this" {
  function_name = "tech-challenge-user-authentication"
  role          = data.aws_iam_role.lab_role.arn
  package_type  = "Image"
  image_uri     = "${data.aws_ecr_repository.this.repository_url}:${var.image_tag}"

  reserved_concurrent_executions = 3

  timeout     = 30
  memory_size = 128

  vpc_config {
    subnet_ids         = data.aws_subnets.private.ids
    security_group_ids = [aws_security_group.lambda.id]
  }

  environment {
    variables = {
      JWT_SECRET            = var.jwt_secret
      JWT_REFRESH_SECRET    = var.jwt_refresh_secret
      DB_HOST               = var.db_host
      DB_USER               = var.db_user
      DB_PASSWORD           = var.db_password
      DB_NAME               = var.db_name
      DB_PORT               = var.db_port
      DB_SSLMODE            = var.db_sslmode
      DYNAMODB_TABLE_NAME   = var.dynamodb_table_name
      DYNAMODB_ENDPOINT     = var.dynamodb_endpoint
      DD_API_KEY            = var.dd_key
      DD_ENV                = "production"
      DD_SERVICE            = "tech-challenge-user-authentication"
      DD_SITE               = "us5.datadoghq.com"
      DD_VERSION            = "1.0.0"
      DD_APPSEC_SCA_ENABLED = "false"
    }
  }

  image_config {
    command = ["bootstrap"]
  }
}

resource "aws_lambda_function_url" "this" {
  function_name      = aws_lambda_function.this.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "allow_public_url" {
  statement_id           = "AllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.this.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}
