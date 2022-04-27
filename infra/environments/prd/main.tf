terraform {
  required_version = "1.0.11"

  required_providers {
    aws = "3.48.0"
  }

  backend "s3" {
    bucket = "tfstate-toebeans-tcpip2"
    key    = "terraform.tfstate"
    region = "ap-northeast-1"
  }
}

module "http_sg" {
  source      = "../../modules/security_group"
  name        = "http-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 80
  cidr_blocks = ["0.0.0.0/0"]
}

module "https_sg" {
  source      = "../../modules/security_group"
  name        = "https-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 443
  cidr_blocks = ["0.0.0.0/0"]
}

module "http_redirect_sg" {
  source      = "../../modules/security_group"
  name        = "http-redirect-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 8080
  cidr_blocks = ["0.0.0.0/0"]
}
