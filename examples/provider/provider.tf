terraform {
  required_providers {
    medusa = {
      source  = "ikhvost/medusa"
      version = "0.0.1"
    }
  }
}

provider "medusa" {
  url      = "<url>"
  email    = "<email>"
  password = "<token>"
}
