variable "name" {}
variable "vpc_id" {}
variable "port" {}
variable "cidr_blocks" {
  type = list(string)
}