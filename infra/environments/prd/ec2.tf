# 踏み台
module "bastion_sg" {
  source      = "../../modules/security_group"
  name        = "bastion-sg"
  vpc_id      = aws_vpc.toebeans.id
  port        = 22
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_key_pair" "toebeans-ssh-key-pair" {
  key_name   = "toebeans-ssh-key-pair"
  public_key = "${file(var.ssh_public_key_path)}"
}

resource "aws_instance" "bastion" {
  ami                         = "ami-0c3fd0f5d33134a76"
  instance_type               = "t3.micro"
  associate_public_ip_address = "true"
  key_name                    = "${aws_key_pair.toebeans-ssh-key-pair.id}"
  subnet_id                   = aws_subnet.public_0.id
  vpc_security_group_ids      = [module.bastion_sg.security_group_id]
  user_data = "${file("userdata.sh")}"
  tags = {
    Name = "toebeans-bastion"
  }
}
