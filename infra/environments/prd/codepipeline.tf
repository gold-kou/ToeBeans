data "aws_iam_policy_document" "codepipeline" {
  statement {
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:GetObjectVersion",
      "s3:GetBucketVersioning",
      "codebuild:BatchGetBuilds",
      "codebuild:StartBuild",
      "ecs:DescribeServices",
      "ecs:DescribeTaskDefinition",
      "ecs:DescribeTasks",
      "ecs:ListTasks",
      "ecs:RegisterTaskDefinition",
      "ecs:UpdateService",
      "iam:PassRole",
      "codestar-connections:UseConnection",
    ]
  }
}

module "codepipeline_role" {
  source     = "../../modules/iam"
  name       = "codepipeline"
  identifier = "codepipeline.amazonaws.com"
  policy     = data.aws_iam_policy_document.codepipeline.json
}

resource "aws_codestarconnections_connection" "toebeans" {
  name          = "toebeans-connection"
  provider_type = "GitHub"
}

data "aws_ssm_parameter" "github_token" {
  name = "/github_token"
}

resource "aws_codepipeline" "toebeans" {
  name     = "toebeans"
  role_arn = module.codepipeline_role.iam_role_arn
  stage {
    name = "Source"
    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["Source"]
      configuration = {
        ConnectionArn    = aws_codestarconnections_connection.toebeans.arn
        FullRepositoryId = "gold-kou/ToeBeans"
        BranchName       = "master"
        OutputArtifactFormat = "CODEBUILD_CLONE_REF"
      }
    }
  }

  stage {
    name = "Build"
    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      version          = 1
      input_artifacts  = ["Source"]
      output_artifacts = ["Build"]
      configuration = {
        ProjectName = aws_codebuild_project.toebeans.id
      }
    }
  }

  stage {
    name = "Deploy"
    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "ECS"
      version         = 1
      input_artifacts = ["Build"]
      configuration = {
        ClusterName = aws_ecs_cluster.toebeans.name
        ServiceName = aws_ecs_service.toebeans.name
        FileName    = "imagedefinitions.json"
      }
    }
  }

  artifact_store {
    location = aws_s3_bucket.artifact.id
    type     = "S3"
  }
}

resource "aws_codepipeline_webhook" "toebeans" {
  name            = "toebeans"
  target_pipeline = aws_codepipeline.toebeans.name
  target_action   = "Source"
  authentication  = "GITHUB_HMAC"
  authentication_configuration {
    secret_token = var.codepipeline_webhook_secret_token
  }
  filter {
    json_path    = "$.ref"
    match_equals = "refs/heads/{Branch}"
  }
}

resource "github_repository_webhook" "toebeans" {
  repository = "ToeBeans"
  configuration {
    url          = aws_codepipeline_webhook.toebeans.url
    secret       = var.github_webhook_secret
    content_type = "json"
    insecure_ssl = false
  }
  events = ["push"]
}
