terraform {
  required_providers {
    thinkplace = {
      source = "hashicorp.com/edu/thinkplace"
    }
  }
}

variable "thinkplace_base_url" {
  type = string
}
variable "thinkplace_token" {
  type = string
}

provider "thinkplace" {
  base_url  = var.thinkplace_base_url
  api_token = var.thinkplace_token
}

data "thinkplace_group" "group" {
  name = "mynewgroup"
}

output "description" {
  value = data.thinkplace_group.group.description
}

resource "thinkplace_group" "other" {
  name        = "tfgroup"
  description = "I made this with Terraform and then I changed it"
}
