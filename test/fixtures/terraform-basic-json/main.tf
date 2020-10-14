/*
Source file for plan.json in the same folder

To (re)generate plan.json:
  terraform init
  terraform plan -out temp.plan
  terraform show -json temp.plan > plan.json
*/

variable example_any {
  default     = null
  description = "An example variable that can be anything"
}

output example_any {
  value = var.example_any
}

variable example_list {
  type        = list
  default     = []
  description = "An example variable that is a list"
}

output example_list {
  value = var.example_list
}

variable example_map {
  type        = map
  default     = {}
  description = "An example variable that is a map"
}

output example_map {
  value = var.example_map
}

resource local_file "example" {
  filename = "example.txt"
  content  = "example + test"
}

resource local_file "example2" {
  filename = "example2.txt"
  content  = "test"
}
