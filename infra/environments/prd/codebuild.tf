data "aws_iam_policy_document" "codebuild" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:GetObjectVersion",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetRepositoryPolicy",
      "ecr:DescribeRepositories",
      "ecr:ListImages",
      "ecr:DescribeImages",
      "ecr:BatchGetImage",
      "ecr:InitiateLayerUpload",
      "ecr:UploadLayerPart",
      "ecr:CompleteLayerUpload",
      "ecr:PutImage",
      "codestar-connections:UseConnection",
      "ssm:GetParameters"
    ]
  }
}

module "codebuild_role" {
  source     = "../../modules/iam"
  name       = "codebuild"
  identifier = "codebuild.amazonaws.com"
  policy     = data.aws_iam_policy_document.codebuild.json
}

resource "aws_codebuild_project" "toebeans" {
  name         = "toebeans"
  service_role = module.codebuild_role.iam_role_arn

  source {
    type = "CODEPIPELINE"
  }

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    type            = "LINUX_CONTAINER"
    compute_type    = "BUILD_GENERAL1_SMALL"
    image           = "aws/codebuild/standard:2.0"
    privileged_mode = true
    environment_variable {
      name  = "DOCKER_HUB_ID"
      type = "PARAMETER_STORE"
      value = "docker_hub_id"
    }
    environment_variable {
      name  = "DOCKER_HUB_PASSWORD"
      type = "PARAMETER_STORE"
      value = "docker_hub_password"
    }
  }
}
