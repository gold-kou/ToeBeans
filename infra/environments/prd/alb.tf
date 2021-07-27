resource "aws_lb" "toebeans" {
  name                       = "toebeans"
  load_balancer_type         = "application"
  internal                   = false
  idle_timeout               = 60
#   enable_deletion_protection = true

  subnets = [
    aws_subnet.public_0.id,
    aws_subnet.public_1.id,
  ]

  access_logs {
    bucket  = aws_s3_bucket.alb_log.id
    enabled = true
  }

  security_groups = [
    module.http_sg.security_group_id,
    module.https_sg.security_group_id,
    module.http_redirect_sg.security_group_id,
  ]
}

# resource "aws_lb_listener" "http" {
#   load_balancer_arn = aws_lb.toebeans.arn
#   port              = "80"
#   protocol          = "HTTP"

#   default_action {
#     type = "fixed-response"

#     fixed_response {
#       content_type = "text/plain"
#       message_body = "これは『HTTP』です"
#       status_code  = "200"
#     }
#   }
# }

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.toebeans.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate.toebeans.arn
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.toebeans.arn
  }
}

# resource "aws_lb_listener" "redirect_http_to_https" {
#   load_balancer_arn = aws_lb.toebeans.arn
#   port              = "8080"
#   protocol          = "HTTP"

#   default_action {
#     type = "redirect"

#     redirect {
#       port        = "443"
#       protocol    = "HTTPS"
#       status_code = "HTTP_301"
#     }
#   }
# }

resource "aws_lb_target_group" "toebeans" {
  name                 = "toebeans"
  target_type          = "ip"
  vpc_id               = aws_vpc.toebeans.id
  port                 = 80
  protocol             = "HTTP"
  deregistration_delay = 300

  health_check {
    path                = "/health/readiness"
    healthy_threshold   = 5
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    matcher             = 200
    port                = "traffic-port"
    protocol            = "HTTP"
  }

  depends_on = [aws_lb.toebeans]
}

resource "aws_lb_listener_rule" "toebeans" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.toebeans.arn
  }

  condition {
    path_pattern {
      values = ["/*"]
    }
  }
}
