resource "aws_acm_certificate" "toebeans" {
  domain_name               = aws_route53_record.toebeans.name
  subject_alternative_names = []
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "toebeans" {
  certificate_arn         = aws_acm_certificate.toebeans.arn
  validation_record_fqdns = [for record in aws_route53_record.toebeans_certificate : record.fqdn]
}
