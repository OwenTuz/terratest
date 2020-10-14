/*
Source file for plan.json in the same folder

To (re)generate plan.json:
  terraform init
  terraform plan -out temp.plan
  terraform show -json temp.plan > plan.json
*/

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  assign_generated_ipv6_cidr_block = false
  cidr_block                       = "10.10.0.0/16"
  enable_dns_support               = true
  instance_tenancy                 = "default"
  tags = {
    Name = "terraform-network-example"
  }
}

output "main_vpc_id" {
  description = "The main VPC ID"
  value       = aws_vpc.main
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
  }
  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "91.189.0.0/24"
  }
  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_subnet" "private" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.10.1.0/24"
  map_public_ip_on_launch = false
  vpc_id                  = aws_vpc.main.id
  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_route_table_association" "private" {
  route_table_id = aws_route_table.private.id
  subnet_id      = aws_subnet.private.id
}

output "private_subnet_id" {
  description = "The private subnet ID"
  value       = "aws_subnet.private"
}

resource "aws_subnet" "public" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.10.2.0/24"
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_route_table_association" "public" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public.id
}

output "public_subnet_id" {
  description = "The public subnet ID"
  value       = "aws_subnet.public"
}

resource "aws_internet_gateway" "main_gateway" {
  tags = {
    Name = "terraform-network-example"
  }
}

resource "aws_eip" "nat" {
  vpc        = true
  depends_on = [aws_internet_gateway.main_gateway]
}

resource "aws_nat_gateway" "nat" {
  depends_on    = [aws_internet_gateway.main_gateway]
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public.id
}
