resource "aws_kms_key" "toebeans" {
  description             = "Customer Master Key"
  enable_key_rotation     = true
  is_enabled              = true
  deletion_window_in_days = 30
}

resource "aws_kms_alias" "toebeans" {
  name          = "alias/toebeans"
  target_key_id = aws_kms_key.toebeans.key_id
}

