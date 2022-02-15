# fortiadc_loadbalance_pool

`fortiadc_loadbalance_pool` is a resource for managing a real server pool.

## Example Usage

```hcl
resource "fortiadc_loadbalance_pool" "mypool" {
  name              = "mypool"
  healtcheck_enable = true
  healtcheck_list   = ["LB_HLTHCK_HTTP", "LB_HLTHCK_HTTPS"]
}
```

## Argument Reference

* `name` - (Required) Pool name.
* `pool_type` - (Optional) Pool type (ipv4/ipv6). Defaults to `ipv4`.
* `healtcheck_enable` - (Optional) Enable healthchecks. Defaults to `false`.
* `healtcheck_relationship` - (Optional) Healtchecks relationship (AND/OR). Defaults to `AND`.
* `healtcheck_list` - (Optional) Healtchecks list. Defaults to `[]`.
* `real_server_ssl_profile` - (Optional) Real servers SSL profile. Defaults to `NONE`.

## Attribute Reference

* `id` - Pool Mkey (internal ID).

## Import

Pools can be imported using their name:

```
$ terraform import fortiadc_loadbalance_pool.mypool mypool
```