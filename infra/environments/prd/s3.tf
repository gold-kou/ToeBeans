# ALB log
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

# artifact
resource "aws_s3_bucket" "artifact" {
  bucket = "artifact-toebeans"

  lifecycle_rule {
    enabled = true

    expiration {
      days = "180"
    }
  }
}

# Backend
resource "aws_s3_bucket" "postings" {
  bucket = "toebeans-postings-tcpip" // MEMO AWSアカウントごとに変える
  cors_rule {
    allowed_origins = ["https://(api.)?(toebeans.tk)$"]
    allowed_methods = ["POST", "GET", "PUT", "DELETE"]
    allowed_headers = ["*"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket" "icons" {
  bucket = "toebeans-icons-tcpip" // MEMO AWSアカウントごとに変える
  cors_rule {
    allowed_origins = ["https://(api.)?(toebeans.tk)$"]
    allowed_methods = ["POST", "GET", "PUT", "DELETE"]
    allowed_headers = ["*"]
    max_age_seconds = 3000
  }
}

# Frontend
resource "aws_s3_bucket" "front_bucket" {
  bucket = "toebeans-front-bucket-tcpip"
  acl = "private" // CloudFrontのみからのアクセスとするためprivate

  website {
    index_document = "index.html"
    # error_document = "error.html" // いずれ用意したい
  } 
}

resource "aws_s3_bucket_policy" "front_bucket_policy" {
    bucket = aws_s3_bucket.front_bucket.id
    policy = data.aws_iam_policy_document.toebeans_frontend.json
}

data "aws_iam_policy_document" "toebeans_frontend" {
  statement {
    sid = "Allow CloudFront"
    effect = "Allow"
    principals {
        type = "AWS"
        identifiers = [aws_cloudfront_origin_access_identity.toebeans_frontend.iam_arn]
    }
    actions = [
        "s3:GetObject"
    ]

    resources = [
        "${aws_s3_bucket.front_bucket.arn}/*"
    ]
  }
}