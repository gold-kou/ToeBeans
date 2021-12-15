resource "aws_ecs_cluster" "toebeans" {
  name = "toebeans"
}

resource "aws_ecs_task_definition" "toebeans" {
  family                   = "toebeans"
  cpu                      = "256"
  memory                   = "512"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  container_definitions    = file("./container_definitions.json")
  execution_role_arn       = module.ecs_task_execution_role.iam_role_arn
}

resource "aws_ecs_service" "toebeans" {
  name                              = "toebeans"
  cluster                           = aws_ecs_cluster.toebeans.arn
  task_definition                   = aws_ecs_task_definition.toebeans.arn
  desired_count                     = 0 // MEMO 節約のため一時的に0
  launch_type                       = "FARGATE"
  platform_version                  = "1.3.0"
  health_check_grace_period_seconds = 60

  network_configuration {
    assign_public_ip = false
    security_groups  = [module.app_sg.security_group_id]

    subnets = [
      aws_subnet.private_0.id,
      aws_subnet.private_1.id,
    ]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.toebeans.arn
    container_name   = "toebeans"
    container_port   = 8080
  }

  lifecycle {
    ignore_changes = [task_definition]
  }
}

module "app_sg" {
  source      = "../../modules/security_group"
  name        = "app-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 80
  cidr_blocks = [aws_vpc.toebeans.cidr_block]
}

data "aws_iam_policy" "ecs_task_execution_role_policy" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "ecs_task_execution" {
  source_json = data.aws_iam_policy.ecs_task_execution_role_policy.policy

  statement {
    effect    = "Allow"
    actions   = [
      "ssm:GetParameters",
      "kms:Decrypt",
      "rds:*",
      "ses:SendEmail",
      "ssmmessages:CreateControlChannel",
      "ssmmessages:CreateDataChannel",
      "ssmmessages:OpenControlChannel",
      "ssmmessages:OpenDataChannel"
    ]
    resources = ["*"]
  }
}

module "ecs_task_execution_role" {
  source     = "../../modules/iam"
  name       = "ecs-task-execution"
  identifier = "ecs-tasks.amazonaws.com"
  policy     = data.aws_iam_policy_document.ecs_task_execution.json
}
