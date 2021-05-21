terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.42"
    }
  }
}

provider "aws" {
  region = var.aws_region
  skip_credentials_validation = true
  skip_requesting_account_id = true

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
