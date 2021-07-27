data "aws_route53_zone" "toebeans" {
  name = "toebeans.tk"
}

resource "aws_route53_record" "toebeans" {
  zone_id = data.aws_route53_zone.toebeans.zone_id
  name    = data.aws_route53_zone.toebeans.name
  type    = "A"

  alias {
    name                   = aws_lb.toebeans.dns_name
    zone_id                = aws_lb.toebeans.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "toebeans_certificate" {
  for_each = {
    for dvo in aws_acm_certificate.toebeans.domain_validation_options : dvo.domain_name => {
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