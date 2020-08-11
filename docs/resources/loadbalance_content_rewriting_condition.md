# fortiadc_loadbalance_content_rewriting_condition

`fortiadc_loadbalance_content_rewriting_condition` is a resource for managing a content rewriting condition.

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

resource "fortiadc_loadbalance_content_rewriting_condition" "myrwcond" {
  content_rewriting = fortiadc_loadbalance_content_rewriting.myrw.name
  object            = "http-host-header"
  type              = "string"
  content           = "myvhost.domain.loc"
}
```

## Argument Reference

* `content_rewriting` - (Required) Parent content rewriting name.
* `object` - (Required) Condition object.
* `type` - (Required) Condition type.
* `content` - (Required) Matching content.
* `reverse` - (Optional) Enable/disable reverse. Defaults to `false`.
* `ignore_case` - (Optional) Ignore matching case. Defaults to `true`.

## Attribute Reference

* `id` - Content rewriting condition Mkey (internal ID).