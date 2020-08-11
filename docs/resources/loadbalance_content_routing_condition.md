# fortiadc_loadbalance_content_routing_condition

`fortiadc_loadbalance_content_routing_condition` is a resource for managing a content routing condition.

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

resource "fortiadc_loadbalance_content_routing_condition" "mycrcond" {
  content_routing = fortiadc_loadbalance_content_routing.mycr.name
  object          = "http-request-url"
  type            = "string"
  content         = "myvhost.domain.loc"
}
```

## Argument Reference

* `content_routing` - (Required) Parent content routing name.
* `object` - (Required) Matching object type (ex: http-host-header).
* `type` - (Required) Matching comparison (ex: string).
* `content` - (Required) Matching content.
* `reverse` - (Optional) Enable reverse. Defaults to `false`.

## Attribute Reference

* `id` - Content routing condition (internal ID).