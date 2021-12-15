resource "aws_acm_certificate" "toebeans_frontend" {
  provider                  = aws.virginia # CloudFrontでHTTPS使う場合はus-east-1になる
  domain_name               = aws_route53_record.toebeans_frontend.name
  subject_alternative_names = []
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate" "toebeans_backend" {
  domain_name               = aws_route53_record.toebeans_backend.name
  subject_alternative_names = []
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "toebeans_frontend" {
  provider                  = aws.virginia
  certificate_arn         = aws_acm_certificate.toebeans_frontend.arn
  validation_record_fqdns = [for record in aws_route53_record.toebeans_frontend_certificate : record.fqdn]
}

resource "aws_acm_certificate_validation" "toebeans_backend" {
  certificate_arn         = aws_acm_certificate.toebeans_backend.arn
  validation_record_fqdns = [for record in aws_route53_record.toebeans_backend_certificate : record.fqdn]
}
