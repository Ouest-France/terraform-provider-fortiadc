# fortiadc_loadbalance_pool_member

`fortiadc_loadbalance_pool_member` is a resource for managing a real server pool member.

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
```

## Argument Reference

* `name` - (Required) Real server name.
* `pool` - (Required) Pool name.
* `port` - (Required) Port.
* `status` - (Optional) Member status (enable/disable/maintain). Defaults to `enable`.
* `weight` - (Optional) Weight. Defaults to `1`.
* `conn_limit` - (Optional) Connection limit. Defaults to `0`.
* `conn_rate_limit` - (Optional) Connection Rate Limit. Defaults to `0`.
* `recover` - (Optional) Recover. Defaults to `0`.
* `warmup` - (Optional) Warmup. Defaults to `0`.
* `warmrate` - (Optional) Warm rate. Defaults to `100`.

## Attribute Reference

* `id` - Member Mkey (internal ID).

## Import

Pool members can be imported using their pool and name joined by a dot:

```
$ terraform import fortiadc_loadbalance_pool_member.mymember mypool.mymember
```