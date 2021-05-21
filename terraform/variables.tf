variable "aws_region" {
  type = string
  description = "AWS region"
  default = "eu-west-3"
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
      name = "ControlPlane"
      hash_key = "Room"

      attributes = [
        {
          name = "Room"
          type = "S"
        }
      ]
    },
    {
      name = "Authentication"
      hash_key = "Username"
      attributes = [
        {
          name = "Username"
          type = "S"
        }
      ]
    },
    {
      name = "TemperatureOutside"
      hash_key = "Date"

      attributes = [
        {
          name = "Date"
          type = "S"
        }
      ]
    },
    {
      name = "TemperatureInside"
      hash_key = "Date"

      attributes = [
        {
          name = "Date"
          type = "S"
        }
      ]
    }
  ]
}
