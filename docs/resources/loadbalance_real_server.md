# fortiadc_loadbalance_real_server

`fortiadc_loadbalance_real_server` is a resource for managing a real server.

## Example Usage

```hcl
resource "fortiadc_loadbalance_real_server" "myrealserver" {
  name     = "myrealserver"
  address  = "192.168.10.20"
  address6 = "::"
  status   = "enable"
}

resource "fortiadc_loadbalance_real_server" "myrealserver2" {
  name   = "myrealserver"
  type   = "fqdn"
  fqdn   = "realserver.example.com"
  status = "enable"
}
```

## Argument Reference

* `name` - (Required) Real server name.
* `type` - (Optional) Type, `"ip"` or `"fqdn"`, defaults to `"ip"`
* `address` - (Optional) Real server ipv4 address. Only used when `type="ip"`
* `address6` - (Optional) Real server ipv6 address. Only used when `type="ip"`. Defaults to `::`.
* `fqdn` - (Optional) Real server Fully Qualified Domain. Only used when `type="fqdn"`
* `status` - (Optional) Real server status. Defaults to `enable`.

## Attribute Reference

* `id` - Real server Mkey (internal ID).

## Import

Real servers can be imported using their name:

```
$ terraform import fortiadc_loadbalance_real_server.myrealserver myrealserver
```