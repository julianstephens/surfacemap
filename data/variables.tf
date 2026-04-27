variable "lambda" {
  description = "Configuration for the Lambda function"
  type = object({
    function_name = string
    role_arn      = string
    handler       = string
    runtime       = string
    filename      = string
    tags          = map(string)
  })
}

variable "environment" {
  description = "The environment for the resources (e.g., dev, staging, prod)"
  type        = string
}

variable "project" {
  description = "The name of the project"
  type        = string
}