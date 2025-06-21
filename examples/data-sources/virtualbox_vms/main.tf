terraform {
  required_providers {
    virtualbox = {
      source  = "registry.terraform.io/apriliantocecep/virtualbox"
    }
  }
}

provider "virtualbox" {}

data "virtualbox_vms" "vms" {}

output "vms_list" {
  value = data.virtualbox_vms.vms
}