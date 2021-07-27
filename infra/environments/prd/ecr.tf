resource "aws_ecr_repository" "toebeans" {
  name = "toebeans"
}

resource "aws_ecr_lifecycle_policy" "toebeans" {
  repository = aws_ecr_repository.toebeans.name

  policy = <<EOF
  {
    "rules": [
      {
        "rulePriority": 1,
        "description": "Keep last 30 release tagged images",
        "selection": {
          "tagStatus": "tagged",
          "tagPrefixList": ["release"],
          "countType": "imageCountMoreThan",
          "countNumber": 30
        },
        "action": {
          "type": "expire"
        }
      }
    ]
  }
EOF
}
