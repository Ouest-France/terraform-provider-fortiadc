# `terraform-provider-fortiadc`

A [Terraform][1] plugin for managing [FortiADC][2].

## Contents

* [Installation](#installation)
* [`fortiadc` Provider](#provider-configuration)
* [Resources](#resources)
  * [`fortiadc_loadbalance_real_server`](#fortiadc_loadbalance_real_server)
  * [`fortiadc_loadbalance_real_server_pool`](#fortiadc_loadbalance_real_server_pool)
* [Requirements](#requirements)

## Installation

Download and extract the [latest
release](https://github.com/Ouest-France/terraform-provider-fortiadc/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

## Provider Configuration

### Example

Example provider.
```hcl
provider "fortiadc" {
  address  = "https://fortiadc.mydomain.com"
  user     = "myuser"
  password = "mypassword"
}
```

| Property            | Description                | Type    | Required    | Default    |
| ----------------    | -----------------------    | ------- | ----------- | ---------- |
| `address`           | fortiadc server address    | String  | true        |            |
| `user`              | fortiadc username          | String  | true        |            |
| `password`          | fortiadc password          | String  | true        |            |
| `insecure`          | disable tls verify         | Bool    | false       | `false`    |

## Resources
### `fortiadc_loadbalance_real_server`

A resource for managing real server.

#### Example

```hcl
resource "fortiadc_loadbalance_real_server" "myrealserver" {
  name     = "myrealserver"
  address  = "192.168.10.20"
  address6 = "::"
  status   = "enable"
}
```

#### Arguments

| Property            | Description                | Type    | Required    | Default    |
| ----------------    | -----------------------    | ------- | ----------- | ---------- |
| `name`              | Real server name           | String  | true        |            |
| `address`           | Real server ipv4 address   | String  | true        |            |
| `address6`          | Real server ipv4 address6  | String  | false       | `::`       |
| `status`            | Real server status         | String  | false       | `enable`   |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Real server Mkey                               |

### `fortiadc_loadbalance_pool`

A resource for managing real server pool.

#### Example

```hcl
resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}
```

#### Arguments

| Property                      | Description                          | Type        | Required    | Default    |
| ----------------              | -----------------------              | -------     | ----------- | ---------- |
| `name`                        | Pool name                            | String      | true        |            |
| `pool_type`                   | Pool type (ipv4/ipv6)                | String      | false       | `ipv4`     |
| `healtcheck_enable`           | Enable healthchecks                  | Bool        | false       | `false`    |
| `healtcheck_relationship`     | Healtchecks relationship (AND/OR)    | String      | false       | `AND`      |
| `healtcheck_list`             | Healtchecks list                     | ListString  | false       | `[]`       |
| `real_server_ssl_profile`     | Real servers SSL profile             | String      | false       | `NONE`     |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Pool Mkey                                      |


## Requirements
* FortiADC == 4.8.x

[1]: https://www.terraform.io
[2]: https://www.fortinet.com/products/application-delivery-controller/fortiadc.html
