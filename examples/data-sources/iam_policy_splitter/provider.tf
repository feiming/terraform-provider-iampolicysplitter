terraform {
  required_providers {
    iampolicysplitter = {
      source  = "feiming/iampolicysplitter"
      version = "~> 0.0.2"
    }
  }
}

provider "iampolicysplitter" {
  # no provider-specific config
}