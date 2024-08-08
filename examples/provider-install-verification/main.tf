terraform {
  required_providers {
    someapi = {
      source = "hashicorp.com/edu/someapi"
    }
  }
}

variable "someapi_base_url" {
  type = string
}
variable "someapi_token" {
  type = string
}

provider "someapi" {
  base_url  = var.someapi_base_url
  api_token = var.someapi_token
}

data "someapi_group" "group" {
  name = "mynewgroup"
}

output "description" {
  value = data.someapi_group.group.description
}

resource "someapi_group" "other" {
  name        = "tfgroup"
  description = "I made this with Terraform and then I changed it"
}
