# MedusaJS Terraform Provider

The Terraform Medusajs provider allows you to configure your
[medusajs application](https://medusajs.com/) space with infrastructure-as-code
principles.

## Compatibility
This release **v0.0.1** is compatible with MedusaJS **v1.20.6**

## Usage

The provider is distributed via the Terraform registry. To use it you need to configure
the [`required_provider`](https://www.terraform.io/language/providers/requirements#requiring-providers) block. For example:

```hcl
terraform {
  required_providers {
    medusa = {
      source  = "ikhvost/medusa"

      # It's recommended to pin the version, e.g.:
      # version = "~> 0.0.1"
    }
  }
}

provider "medusa" {
  url      = "<url>"
  email    = "<email>"
  password = "<token>"
}
```

# Binaries

Packages of the releases are available at https://github.com/ikhvost/terraform-provider-medusa/releases 
See the [terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more information about installing third-party providers.

# Contributing

## Building the provider

Clone the repository and run the following command:

```sh
$ task build-local
```

## Debugging / Troubleshooting

There are two environment settings for troubleshooting:

- `TF_LOG=INFO` enables debug output for Terraform.

Note this generates a lot of output!
