terraform {
    backend "s3" {
        bucket = "demo-terraform-state"
        key    = "terraform.tfstate"
        region = "us-east-1"
        dynamodb_table = "tf_state_lock"
    }
    required_version = "= 1.14.8"
    required_providers {
      aws = {
        source = "hashicorp/aws"
        version = ">= 6.0.0"
      }
    }
}

provider "aws" {
    region = "us-east-1"
}

module "demo" {
    source = "../.."
    environment = "development"    
    project = "demo"
    lambda = {
        function_name = "demo-function"
        role_arn      = aws_iam_role.lambda_execution_role.arn
        handler       = "index.handler"
        runtime       = "nodejs18.x"
        filename      = "lambda.zip"
    }
}