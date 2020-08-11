# FortiADC Provider

The FortiADC provider is used to interact with the resources supported by [FortiADC](https://www.fortinet.com/products/application-delivery-controller/fortiadc) API. The provider needs to be configured with the proper credentials and address before it can be used.

## Example Usage

```hcl
# Provider configuration
terraform {
  required_providers {
    fortiadc = {
      source  = "Ouest-France/fortiadc"
    }
  }
}

provider "fortiadc" {
  address  = "https://fortiadc.mydomain.com"
  user     = "myuser"
  password = "mypassword"
}

# Real server definition
resource "fortiadc_loadbalance_real_server" "myrealserver" {
  name    = "myrealserver"
  address = "192.168.1.55"
  status  = "enable"
}

...
```

## Argument Reference

* `address` - (Required) This is the FortiADC address formatted like `https://myfortisrv.mydomain.com`.

* `user` - (Required) This is the FortiADC user to access the API.

* `password` - (Required) This is the FortiADC password to access the API.

* `insecure` - (Optional) This enable or disable TLS certificate verification, defaults to `false`.

## Requirements

* FortiADC >= 5.1