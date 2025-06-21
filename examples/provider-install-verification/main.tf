terraform {
  required_providers {
    virtualbox = {
      source  = "registry.terraform.io/apriliantocecep/virtualbox"
    }
  }
}

provider "virtualbox" {}