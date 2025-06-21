terraform {
  required_providers {
    virtualbox = {
      source  = "registry.terraform.io/apriliantocecep/virtualbox"
    }
  }
}

provider "virtualbox" {}

resource "virtualbox_vm" "ubuntu" {
  name = "ubuntu-vm-64"
  cpus = 2
  memory = 2048
}

output "vm_ubuntu" {
  value = virtualbox_vm.ubuntu
}