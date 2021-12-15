resource "aws_cloudfront_distribution" "toebeans_frontend" {
  origin {
    domain_name = aws_s3_bucket.front_bucket.bucket_regional_domain_name
    origin_id = aws_s3_bucket.front_bucket.id
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.toebeans_frontend.cloudfront_access_identity_path
    }
  }

  enabled =  true

  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods = [ "GET", "HEAD" ]
    cached_methods = [ "GET", "HEAD" ]
    target_origin_id = aws_s3_bucket.front_bucket.id
        
    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
      min_ttl = 0
      default_ttl = 3600
      max_ttl = 86400
  }

  custom_error_response {
    error_caching_min_ttl = 300
    error_code            = 403
    response_code         = 200
    response_page_path    = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations = [ "JP" ]
    }
  }
  
  viewer_certificate {
    cloudfront_default_certificate = true # ↓の行だとcycle errorになってしまう。いったんdefaultにしてからコンソールで変更する。
    # acm_certificate_arn      = aws_acm_certificate.toebeans_frontend.id
    minimum_protocol_version = "TLSv1.2_2019"
    ssl_support_method       = "sni-only"
  }

  # バージニアのACMを指定できないため以下のコメントアウトを外すとエラーになる。
  # aliases = ["toebeans.tk"]
}

resource "aws_cloudfront_origin_access_identity" "toebeans_frontend" {}