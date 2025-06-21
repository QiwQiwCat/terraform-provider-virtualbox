terraform {
  required_providers {
    virtualbox = {
      source  = "registry.terraform.io/apriliantocecep/virtualbox"
    }
  }
}

provider "virtualbox" {}

data "virtualbox_ostypes" "ostypes" {}

output "ostypes_list" {
  value = data.virtualbox_ostypes.ostypes
}