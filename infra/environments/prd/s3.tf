resource "aws_s3_bucket" "public" {
  bucket = "public-toebeans"
  acl    = "public-read"

  cors_rule {
    allowed_origins = ["https://toebeans.tk"]
    allowed_methods = ["POST", "GET", "PUT"]
    allowed_headers = ["*"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket" "alb_log" {
  bucket = "alb-log-toebeans"

  lifecycle_rule {
    enabled = true

    expiration {
      days = "180"
    }
  }
}

resource "aws_s3_bucket_policy" "alb_log" {
  bucket = aws_s3_bucket.alb_log.id
  policy = data.aws_iam_policy_document.alb_log.json
}

data "aws_iam_policy_document" "alb_log" {
  statement {
    effect    = "Allow"
    actions   = ["s3:PutObject"]
    resources = ["arn:aws:s3:::${aws_s3_bucket.alb_log.id}/*"]

    principals {
      type        = "AWS"
      identifiers = ["582318560864"]
    }
  }
}

resource "aws_s3_bucket" "artifact" {
  bucket = "artifact-toebeans"

  lifecycle_rule {
    enabled = true

    expiration {
      days = "180"
    }
  }
}
