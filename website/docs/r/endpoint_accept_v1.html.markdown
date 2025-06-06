---
layout: "sci"
page_title: "SAP Cloud Infrastructure: sci_endpoint_accept_v1"
sidebar_current: "docs-sci-resource-endpoint-accept-v1"
description: |-
  Accept an Archer endpoint request.
---

# sci\_endpoint\_accept\_v1

Use this resource to accept an Archer endpoint connection request within the
SAP Cloud Infrastructure environment. This is typically used in scenarios where services
are shared across projects and require explicit acceptance to establish
connectivity.

~> **Note:** The `terraform destroy` command will reject the endpoint
connection.

## Example Usage

```hcl
provider "sci" {
  alias = "remote"

  tenant_name = "remote-project"
}

resource "sci_endpoint_v1" "endpoint_1" {
  provider = sci.remote

  name = "ep1"
  service_id = sci_endpoint_service_v1.service_1.id
  target {
   network = "8cd158e4-7b10-43ad-9b92-4251b10b5ee8"
  }

  tags = ["tag1","tag2"]
}

resource "sci_endpoint_service_v1" "service_1" {
  name             = "svc1"
  visibility       = "public"
  require_approval = true
  ip_addresses     = ["192.168.1.2"]
  network_id       = "a7ec6c35-4e17-4e97-aa2b-0d93e56bb6c7"
  port             = 80
}

resource "sci_endpoint_accept_v1" "accept_1" {
  service_id  = sci_endpoint_service_v1.service_1.id
  endpoint_id = sci_endpoint_v1.endpoint_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to accept the endpoint. If omitted,
  the `region` argument of the provider is used. Changing this forces a new
  resource to be created.

* `service_id` - (Required) The ID of the service to which the endpoint is
  connected.

* `endpoint_id` - (Required) The ID of the endpoint to accept.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A combined ID of the Archer service and endpoint separated by a slash.
* `status` - The current status of the Archer service endpoint acceptance.

## Import

An Archer endpoint consumer can be imported using the endpoint `id` and
`service_id` separated by a slash, e.g.:

```shell
$ terraform import sci_endpoint_accept_v1.accept_1 301317d8-9067-439f-b90f-9916beaf087c/74931fd2-90ff-41c0-93f2-f536eb3c2412
```
