resource "aws_dynamodb_table" "table" {
  name = var.table_name
  billing_mode = "PROVISIONED"

  # AWS Free tier includes a maximum of:
  #   - 25GB
  #   - 25 read capacity units
  #   - 25 write capacity units
  # Here we are limiting read and write capacity to 20
  read_capacity = 20
  write_capacity = 20
  hash_key = var.hash_key

  dynamic "attribute" {
    for_each = var.attributes
    content {
      name = attribute.value["name"]
      type = attribute.value["type"]
    }
  }
}
