resource "aws_cloudwatch_log_group" "toebeans" {
  name              = "/ecs/toebeans"
  retention_in_days = 90
}
