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

resource "aws_ecs_cluster" "example" {
  name = "terratest-example"
}

data "aws_iam_policy_document" "assume-execution" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "execution" {
  name               = "terratest-example"
  assume_role_policy = data.aws_iam_policy_document.assume-execution.json
}

resource "aws_iam_role_policy_attachment" "execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = aws_iam_role.execution.name
}

resource "aws_ecs_task_definition" "example" {
  family             = "terratest"
  execution_role_arn = aws_iam_role.execution.arn

  container_definitions = <<EOF
[
  {
    "name" : "terratest",
    "image": "terratest-example"
  }
]
EOF

  cpu    = 256
  memory = 512

  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
}

resource "aws_ecs_service" "example" {
  name          = "terratest-example"
  desired_count = 0
  launch_type   = "FARGATE"

  cluster         = aws_ecs_cluster.example.id
  task_definition = aws_ecs_task_definition.example.arn
}
