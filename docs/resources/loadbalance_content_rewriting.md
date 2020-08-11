# fortiadc_loadbalance_content_rewriting

`fortiadc_loadbalance_content_rewriting` is a resource for managing a real server.

## Example Usage

```hcl
resource "fortiadc_loadbalance_virtual_server" "myvirtualserver" {
  name    = "myvirtualserver"
  address = "192.168.11.10"
  port    = 80
  pool    = "my-rs-pool"

  content_rewriting_enable = true
  content_rewriting_list   = [fortiadc_loadbalance_content_rewriting.myrw.name]
}

resource "fortiadc_loadbalance_content_rewriting" "myrw" {
  name = "my-content-rewriting"
  action_type = "request"
  action      = "send-403-forbidden"
}
```

## Argument Reference

* `name` - (Required) Content rewriting name.
* `action_type` - (Required) Action type.
* `action` - (Required) Action.
* `comment` - (Optional) Comment.
* `host_match` - (Optional) Host header match.
* `url_match` - (Optional) URL path match.
* `referer_match` - (Optional) Referer header match.
* `redirect` - (Optional) Redirect URL.
* `location` - (Optional) HTTP location rewrite.
* `header_name` - (Optional) Header name to add or remove.

## Attribute Reference

* `id` - Content rewriting Mkey (internal ID).