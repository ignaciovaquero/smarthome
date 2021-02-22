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
  for_each = {for table in var.dynamo_db_tables: table.name => table}
  source = "./modules/dynamodb_table"

  table_name = each.value.name
  hash_key = each.value.hash_key

  attributes = each.value.attributes
}
