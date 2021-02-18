terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.28"
    }
  }
}

provider "aws" {
  region = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key

  endpoints {
    dynamodb = var.dynamodb_endpoint
  }
}

module "smarthome_table" {
  source = "./modules/dynamodb_table"

  table_name = "SmartHome"
  hash_key = "Room"

  attributes = [
    {
      name = "Room"
      type = "S"
    }
  ]
}
