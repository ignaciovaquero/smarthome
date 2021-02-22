variable "aws_region" {
  type = string
  description = "AWS region"
  default = "eu-west-3"
}

variable "aws_access_key" {
  type = string
  description = "AWS access key"
}

variable "aws_secret_key" {
  type = string
  description = "AWS secret key"
}

variable "dynamodb_endpoint" {
  type = string
  description = "DynamoDB endpoint"
  default = "https://dynamodb.eu-west-3.amazonaws.com"
}

variable "dynamo_db_tables" {
  type = list(object({
    name = string
    hash_key = string
    attributes = list(object({
      name = string
      type = string
    }))
  }))

  default = [
    {
      name = "SmartHome"
      hash_key = "room"

      attributes = [
        {
          name = "room"
          type = "S"
        }
      ]
    }
  ]
}
