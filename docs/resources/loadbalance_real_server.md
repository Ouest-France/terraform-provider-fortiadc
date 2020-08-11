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
```

## Argument Reference

* `name` - (Required) Real server name.
* `address` - (Required) Real server ipv4 address.
* `address6` - (Optional) Real server ipv6 address. Defaults to `::`.
* `status` - (Optional) Real server status. Defaults to `enable`.

## Attribute Reference

* `id` - Real server Mkey (internal ID).