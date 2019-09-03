# `terraform-provider-fortiadc`

A [Terraform][1] plugin for managing [FortiADC][2].

## Contents

* [Installation](#installation)
* [`fortiadc` Provider](#provider-configuration)
* [Resources](#resources)
  * [`fortiadc_loadbalance_real_server`](#fortiadc_loadbalance_real_server)
  * [`fortiadc_loadbalance_pool`](#fortiadc_loadbalance_pool)
  * [`fortiadc_loadbalance_pool_member`](#fortiadc_loadbalance_pool_member)
  * [`fortiadc_loadbalance_virtual_server`](#fortiadc_loadbalance_virtual_server)
  * [`fortiadc_loadbalance_content_routing`](#fortiadc_loadbalance_content_routing)
  * [`fortiadc_loadbalance_content_routing_condition`](#fortiadc_loadbalance_content_routing_condition)
  * [`fortiadc_loadbalance_content_rewriting`](#fortiadc_loadbalance_content_rewriting)
  * [`fortiadc_loadbalance_content_rewriting_condition`](#fortiadc_loadbalance_content_rewriting_condition)
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

### `fortiadc_loadbalance_pool_member`

A resource for managing real server pool member.

#### Example

```hcl
resource "fortiadc_loadbalance_real_server" "myrealserver" {
  name     = "myrealserver"
  address  = "192.168.10.20"
  address6 = "::"
  status   = "enable"
}

resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}

resource "fortiadc_loadbalance_pool_member" "mymember" {
  name = "${fortiadc_loadbalance_real_server.myrealserver.name}"
  pool = "${fortiadc_loadbalance_pool.mypool.name}"
  port = 80
}
```

#### Arguments

| Property                      | Description                             | Type        | Required    | Default    |
| ----------------              | -----------------------                 | -------     | ----------- | ---------- |
| `name`                        | Real server name                        | String      | true        |            |
| `pool`                        | Pool name                               | String      | true        |            |
| `status`                      | Member status (enable/disable/maintain) | String      | false       | `enable`   |
| `port`                        | Port                                    | Int         | true        |            |
| `weight`                      | Weight                                  | Int         | false       | `1`        |
| `conn_limit`                  | Connection limit                        | Int         | false       | `0`        |
| `conn_rate_limit`             | Connection Rate Limit                   | Int         | false       | `0`        |
| `recover`                     | Recover                                 | Int         | false       | `0`        |
| `warmup`                      | Warm Up                                 | Int         | false       | `0`        |
| `warmrate`                    | Warm Rate                               | Int         | false       | `100`      |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Member Mkey                                    |

### `fortiadc_loadbalance_virtual_server`

A resource for managing virtual server.

#### Example

```hcl
resource "fortiadc_loadbalance_real_server" "myrealserver" {
  name     = "myrealserver"
  address  = "192.168.10.20"
  address6 = "::"
  status   = "enable"
}

resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}

resource "fortiadc_loadbalance_pool_member" "mymember" {
  name = "${fortiadc_loadbalance_real_server.myrealserver.name}"
  pool = "${fortiadc_loadbalance_pool.mypool.name}"
  port = 80
}

resource "fortiadc_loadbalance_virtual_server" "myvirtualserver" {
  name    = "myvirtualserver"
  address = "192.168.11.10"
  port    = 80
  pool    = "${fortiadc_loadbalance_pool.mypool.name}"
}
```

#### Arguments

| Property                  | Description                             | Type        | Required    | Default                 |
| ----------------          | -----------------------                 | -------     | ----------- | ----------              |
| `name`                    | Virtual server name                     | String      | true        |                         |
| `status`                  | Status (enable/disable/maintain)        | String      | false       | `enable`                |
| `type`                    | Type                                    | String      | false       | `l4-load-balance`       |
| `address_type`            | Address type (ipv4/ipv6)                | String      | true        | `ipv4`                  |
| `address`                 | Address                                 | String      | true        |                         |
| `packet_forward_method`   | Packet forwarding method (NAT/FullNAT)  | String      | false       | `NAT`                   |
| `nat_source_pool`         | NAT source pool                         | String      | false       | ` `                     |
| `port`                    | Port                                    | Int         | true        |                         |
| `connection_limit`        | Connection limit                        | Int         | false       | `0`                     |
| `content_routing_enable`  | Enable content routing                  | Bool        | false       | `false`                 |
| `content_routing_list`    | List of content routing                 | ListString  | false       | `[]`                    |
| `connection_rate_limit`   | Connection rate limit                   | Int         | false       | `0`                     |
| `interface`               | Interface                               | String      | false       | `port1`                 |
| `profile`                 | Profile                                 | String      | false       | `LB_PROF_TCP`           |
| `method`                  | Method                                  | String      | false       | `LB_METHOD_ROUND_ROBIN` |
| `pool`                    | Pool name                               | String      | true        |                         |
| `client_ssl_profile`      | Client SSL profile                      | String      | false       | ` `                     |
| `http_to_https`           | Redirect HTTP to HTTPS when L7          | Bool        | false       | `false`                 |
| `error_msg`               | Error message on backend failure        | String      | false       | `Server-unavailable!`   |
| `error_page`              | Error page on backend failure           | String      | false       | ``                 |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Virtual server Mkey                            |

### `fortiadc_loadbalance_content_routing`

A resource for managing content routing.

#### Example

```hcl
resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}

resource "fortiadc_loadbalance_content_routing" "mycr" {
  name    = "mycr"
  pool    = "${fortiadc_loadbalance_pool.mypool.name}"
}
```

#### Arguments

| Property                  | Description                             | Type        | Required    | Default                 |
| ----------------          | -----------------------                 | -------     | ----------- | ----------              |
| `name`                    | Content routing name                    | String      | true        |                         |
| `type`                    | Type                                    | String      | false       | `l7-content-routing`    |
| `pool`                    | Real server destination pool            | String      | true        |                         |
| `ipv4`                    | Source IPv4 address for l4 mode         | String      | false       | `0.0.0.0/0`             |
| `ipv6`                    | Source IPv6 address for l4 mode         | String      | false       | `::/0`                  |
| `comment`                 | Comment                                 | String      | false       | `comments`              |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Content routing Mkey                           |


### `fortiadc_loadbalance_content_routing_condition`

A resource for managing content routing condition.

#### Example

```hcl
resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}

resource "fortiadc_loadbalance_content_routing" "mycr" {
  name    = "mycr"
  pool    = "${fortiadc_loadbalance_pool.mypool.name}"
}

resource "fortiadc_loadbalance_content_routing_condition" "mycrcond" {
  content_routing = "${fortiadc_loadbalance_content_routing.mycr.name}"
  object          = "http-request-url"
  type            = "string"
  content         = "myvhost.domain.loc"
}
```

#### Arguments

| Property           | Description                                 | Type        | Required    | Default                 |
| ----------------   | -----------------------                     | -------     | ----------- | ----------              |
| `content_routing`  | Parent content routing name                 | String      | true        |                         |
| `object`           | Matching object type (ex: http-host-header) | String      | true        |                         |
| `type`             | Matching comparison (ex: string)            | String      | true        |                         |
| `content`          | Matching content                            | String      | true        |                         |
| `reverse`          | Enable reverse                              | Bool        | false       | `false`                 |


#### Attributes

| Property             | Description                      |
| ----------------     | -----------------------          |
| `id`                 | Content routing condition Mkey   |


### `fortiadc_loadbalance_content_rewriting`

A resource for managing content rewriting.

#### Example

```hcl
resource "fortiadc_loadbalance_virtual_server" "myvirtualserver" {
  name    = "myvirtualserver"
  address = "192.168.11.10"
  port    = 80
  pool    = "my-rs-pool"

  content_rewriting_enable = true
  content_rewriting_list   = ["${fortiadc_loadbalance_content_rewriting.myrw.name}"]
}

resource "fortiadc_loadbalance_content_rewriting" "myrw" {
  name = "my-content-rewriting"
  action_type = "request"
  action      = "send-403-forbidden"
}
```

#### Arguments

| Property                  | Description                             | Type        | Required    | Default                 |
| ----------------          | -----------------------                 | -------     | ----------- | ----------              |
| `name`                    | Content rewriting name                  | String      | true        |                         |
| `action_type`             | Action type                             | String      | true        |                         |
| `action`                  | Action                                  | String      | true        |                         |
| `comment`                 | Comment                                 | String      | false       | ` `                     |
| `host_match`              | Host header match                       | String      | false       | ` `                     |
| `url_match`               | URL path match                          | String      | false       | ` `                     |
| `referer_match`           | Referer header match                    | String      | false       | ` `                     |
| `redirect`                | Redirect URL                            | String      | false       | ` `                     |
| `location`                | HTTP location rewrite                   | String      | false       | ` `                     |
| `header_name`             | Header name to add or remove            | String      | false       | ` `                     |

#### Attributes

| Property             | Description                                    |
| ----------------     | -----------------------                        |
| `id`                 | Content rewriting Mkey                         |


### `fortiadc_loadbalance_content_rewriting_condition`

A resource for managing content rewriting condition.

#### Example

```hcl
resource "fortiadc_loadbalance_virtual_server" "myvirtualserver" {
  name    = "myvirtualserver"
  address = "192.168.11.10"
  port    = 80
  pool    = "my-rs-pool"

  content_rewriting_enable = true
  content_rewriting_list   = ["${fortiadc_loadbalance_content_rewriting.myrw.name}"]
}

resource "fortiadc_loadbalance_content_rewriting" "myrw" {
  name = "my-content-rewriting"
  action_type = "request"
  action      = "send-403-forbidden"
}

resource "fortiadc_loadbalance_content_rewriting_condition" "myrwcond" {
  content_rewriting = "${fortiadc_loadbalance_content_rewriting.myrw.name}"
  object            = "http-host-header"
  type              = "string"
  content           = "myvhost.domain.loc"
}
```

#### Arguments

| Property           | Description                                 | Type        | Required    | Default                 |
| ----------------   | -----------------------                     | -------     | ----------- | ----------              |
| `content_rewriting`| Parent content rewriting name               | String      | true        |                         |
| `object`           | Condition object                            | String      | true        |                         |
| `type`             | Condition type                              | String      | true        |                         |
| `content`          | Matching content                            | String      | true        |                         |
| `reverse`          | Enable reverse                              | Bool        | false       | `false`                 |
| `ignore_case`      | Ignore matching case                        | Bool        | false       | `true`                  |


#### Attributes

| Property             | Description                      |
| ----------------     | -----------------------          |
| `id`                 | Content rewriting condition Mkey   |

## Requirements
* FortiADC >= 5.1

[1]: https://www.terraform.io
[2]: https://www.fortinet.com/products/application-delivery-controller/fortiadc.html
