data "aws_route53_zone" "toebeans" {
  name = "toebeans.tk"
}

resource "aws_route53_record" "toebeans_frontend" {
  type    = "A"
  name    = data.aws_route53_zone.toebeans.name
  zone_id = data.aws_route53_zone.toebeans.id

  alias {
    name                   = aws_cloudfront_distribution.toebeans_frontend.domain_name
    zone_id                = aws_cloudfront_distribution.toebeans_frontend.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "toebeans_backend" {
  type    = "A"
  name    = "api.${data.aws_route53_zone.toebeans.name}"
  zone_id = data.aws_route53_zone.toebeans.zone_id

  alias {
    name                   = aws_lb.toebeans.dns_name
    zone_id                = aws_lb.toebeans.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "toebeans_frontend_certificate" {
  depends_on = [aws_acm_certificate.toebeans_frontend]
  for_each = {
    for dvo in aws_acm_certificate.toebeans_frontend.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  allow_overwrite = true
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  zone_id = data.aws_route53_zone.toebeans.id
  ttl     = 60
}

resource "aws_route53_record" "toebeans_backend_certificate" {
  depends_on = [aws_acm_certificate.toebeans_backend]
  for_each = {
    for dvo in aws_acm_certificate.toebeans_backend.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  allow_overwrite = true
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  zone_id = data.aws_route53_zone.toebeans.id
  ttl     = 60
}