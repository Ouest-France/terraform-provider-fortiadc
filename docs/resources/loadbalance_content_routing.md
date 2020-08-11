# fortiadc_loadbalance_content_routing

`fortiadc_loadbalance_content_routing` is a resource for managing content routing.

## Example Usage

```hcl
resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}

resource "fortiadc_loadbalance_content_routing" "mycr" {
  name    = "mycr"
  pool    = fortiadc_loadbalance_pool.mypool.name
}
```

## Argument Reference

* `name` - (Required) Content routing name.
* `pool` - (Required) Real server destination pool.
* `type` - (Optional) Content routing type. Defaults to `l7-content-routing`.
* `ipv4` - (Optional) Source IPv4 address for l4 mode. Defaults to `0.0.0.0/0`.
* `ipv6` - (Optional) Source IPv6 address for l4 mode. Defaults to `::/0`.
* `comment` - (Optional) Comment. Defaults to `comments`.

## Attribute Reference

* `id` - Content routing Mkey (internal ID).