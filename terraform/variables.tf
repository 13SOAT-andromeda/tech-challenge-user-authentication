variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "db_host" {
  description = "The database host"
  type        = string
}

variable "dynamodb_table_name" {
  description = "The name of the DynamoDB table"
  type        = string
  default     = "user-authentication-token"
}

variable "dynamodb_endpoint" {
  description = "The endpoint URL of the DynamoDB service"
  type        = string
}

variable "image_tag" {
  description = "ECR image tag"
  type        = string
}

variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

variable "jwt_refresh_secret" {
  description = "JWT refresh secret key"
  type        = string
  sensitive   = true
}

variable "db_user" {
  description = "Database user"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_port" {
  description = "Database port"
  type        = string
  default     = "5432"
}

variable "db_sslmode" {
  description = "Database SSL mode"
  type        = string
  default     = "require"
}

variable "dd_key" {
  description = "Datadog API Key"
  type        = string
  sensitive   = true
}

variable "vpc_id" {
  description = "VPC ID where the Lambda will be deployed"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private Subnet IDs for the Lambda function"
  type        = list(string)
}
