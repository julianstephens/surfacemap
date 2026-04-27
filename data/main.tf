terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">=6.0.0"
    }
  }
  required_version = "= 1.14.8"
}

locals {
  tags = {
    Environment = var.environment
    Project     = var.project
  }
}