variable "table_name" {
  type = string
  description = "The name of the DynamoDB table"
}

variable "attributes" {
  type = list(object({
    name = string
    type = string
  }))
  description = "A list of the attributes for the DynamoDB table"
}

variable "hash_key" {
  type = string
  description = "The Hash Key for the DynamoDB table"
}
