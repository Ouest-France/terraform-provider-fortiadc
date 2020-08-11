# fortiadc_loadbalance_virtual_server

`fortiadc_loadbalance_virtual_server` is a resource for managing a virtual server.

## Example Usage

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
  name = fortiadc_loadbalance_real_server.myrealserver.name
  pool = fortiadc_loadbalance_pool.mypool.name
  port = 80
}

resource "fortiadc_loadbalance_virtual_server" "myvirtualserver" {
  name    = "myvirtualserver"
  address = "192.168.11.10"
  port    = 80
  pool    = fortiadc_loadbalance_pool.mypool.name
}
```

## Argument Reference

* `name` - (Required) Virtual server name.
* `address` - (Required) IP Address.
* `port` - (Required) Port.
* `pool` - (Required) Pool name.
* `status` - (Optional) Status (enable/disable/maintain). Defaults to `enable`.
* `type` - (Optional) Virtual server type. Defaults to `l4-load-balance`.
* `address_type` - (Optional) IP Address type. Defaults to `ipv4`.
* `packet_forward_method` - (Optional) Packet forwarding method (NAT/FullNAT). Defaults to `NAT`.
* `nat_source_pool` - (Optional) NAT source pool.
* `connection_limit` - (Optional) Connection limit. Defaults to `0`.
* `content_routing_enable` - (Optional) Enable content routing. Defaults to `false`.
* `content_routing_list` - (Optional) List of content routing. Defaults to `[]`.
* `connection_rate_limit` - (Optional) Connection rate limit. Defaults to `0`.
* `interface` - (Optional) Network interface. Defaults to `port1`.
* `profile` - (Optional) Profile. Defaults to `LB_PROF_TCP`.
* `method` - (Optional) Loadbalancing method. Defaults to `LB_METHOD_ROUND_ROBIN`.
* `client_ssl_profile` - (Optional) Client SSL profile.
* `http_to_https` - (Optional) Enable/disable redirect of HTTP to HTTPS when L7. Defaults to `false`.
* `error_msg` - (Optional) Error message on backend failure.
* `error_page` - (Optional) Error page on backend failure.
* `persistence` - (Optional) Persistence configuration.

## Attribute Reference

* `id` - Virtual server Mkey (internal ID).