resource "aws_vpc" "toebeans" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "toebeans"
  }
}

resource "aws_internet_gateway" "toebeans" {
  vpc_id = aws_vpc.toebeans.id
}

# public
resource "aws_subnet" "public_0" {
  vpc_id                  = aws_vpc.toebeans.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "public_0"
  }
}

resource "aws_subnet" "public_1" {
  vpc_id                  = aws_vpc.toebeans.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = true
  tags = {
    Name = "public_1"
  }
}

resource "aws_route_table" "public_0" {
  vpc_id = aws_vpc.toebeans.id
}

resource "aws_route_table" "public_1" {
  vpc_id = aws_vpc.toebeans.id
}

resource "aws_route_table_association" "public_0" {
  subnet_id      = aws_subnet.public_0.id
  route_table_id = aws_route_table.public_0.id
}

resource "aws_route_table_association" "public_1" {
  subnet_id      = aws_subnet.public_1.id
  route_table_id = aws_route_table.public_1.id
}

resource "aws_route" "public_0" {
  route_table_id         = aws_route_table.public_0.id
  gateway_id             = aws_internet_gateway.toebeans.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_route" "public_1" {
  route_table_id         = aws_route_table.public_1.id
  gateway_id             = aws_internet_gateway.toebeans.id
  destination_cidr_block = "0.0.0.0/0"
}

# private
resource "aws_subnet" "private_0" {
  vpc_id                  = aws_vpc.toebeans.id
  cidr_block              = "10.0.65.0/24"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = false
  tags = {
    Name = "private_0"
  }
}

resource "aws_subnet" "private_1" {
  vpc_id                  = aws_vpc.toebeans.id
  cidr_block              = "10.0.66.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = false
  tags = {
    Name = "private_1"
  }
}

# MEMO nat gatewayはお金かかるので開発中は一時的に消す
resource "aws_route_table" "private_0" {
  vpc_id = aws_vpc.toebeans.id
}

resource "aws_route_table" "private_1" {
  vpc_id = aws_vpc.toebeans.id
}

resource "aws_route_table_association" "private_0" {
  subnet_id      = aws_subnet.private_0.id
  route_table_id = aws_route_table.private_0.id
}

resource "aws_route_table_association" "private_1" {
  subnet_id      = aws_subnet.private_1.id
  route_table_id = aws_route_table.private_1.id
}

resource "aws_route" "private_0" {
  route_table_id         = aws_route_table.private_0.id
  nat_gateway_id         = aws_nat_gateway.nat_gateway_0.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_route" "private_1" {
  route_table_id         = aws_route_table.private_1.id
  nat_gateway_id         = aws_nat_gateway.nat_gateway_1.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_eip" "nat_gateway_0" {
  vpc        = true
  depends_on = [aws_internet_gateway.toebeans]
}

resource "aws_eip" "nat_gateway_1" {
  vpc        = true
  depends_on = [aws_internet_gateway.toebeans]
}

resource "aws_nat_gateway" "nat_gateway_0" {
  allocation_id = aws_eip.nat_gateway_0.id
  subnet_id     = aws_subnet.public_0.id
  depends_on    = [aws_internet_gateway.toebeans]
}

resource "aws_nat_gateway" "nat_gateway_1" {
  allocation_id = aws_eip.nat_gateway_1.id
  subnet_id     = aws_subnet.public_1.id
  depends_on    = [aws_internet_gateway.toebeans]
}
# ここまで
