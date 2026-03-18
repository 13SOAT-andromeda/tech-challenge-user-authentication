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

variable "project_env" {
  description = "Project environment"
  type        = string
  default     = "dev"
}
