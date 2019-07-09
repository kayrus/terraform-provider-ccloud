---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_kubernetes_v1"
sidebar_current: "docs-ccloud-resource-kubernetes-v1"
description: |-
  Manages a Kubernikus cluster
---

# ccloud\_kubernetes\_v1

Manages a Kubernikus (Kubernetes as a Service) cluster.

~> **Note:** All arguments and attributes, including basic auth username and
passwords as well as certificate outputs will be stored in the raw state as
plaintext.
[Read more about sensitive data in state](https://www.terraform.io/docs/state/sensitive-data.html).

~> Changing the arguments of Kubernikus node pools (except the `size` or
`config` arguments) will result in the node pool downscaling, deleting and
creating a new node pool with the new argument specified.

## Example Usage

### Kubernikus cluster with two node pools

```hcl
resource "ccloud_kubernetes_v1" "demo" {
  name           = "demo"
  ssh_public_key = "ssh-rsa AAAABHTmDMP6w=="

  node_pool {
    name              = "payload0"
    flavor            = "m1.xlarge_cpu"
    size              = 2
    availability_zone = "eu-de-1d"
    taints            = ["key=value:NoSchedule"]
    labels            = ["label=value"]
  }

  node_pool {
    name              = "payload1"
    flavor            = "m1.xlarge_cpu"
    size              = 1
    availability_zone = "eu-de-1b"
    taints            = ["key=value:NoSchedule"]
    labels            = ["label=value"]
  }
}
```

### Kubernikus cluster in the project with multiple networks and routers

```hcl
data "openstack_networking_network_v2" "fip_network_1" {
  name     = "fip_network_1"
  external = true
}

data "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

data "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${data.openstack_networking_network_v2.network_1.id}"
}

data "openstack_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "ccloud_kubernetes_v1" "demo" {
  name           = "demo"
  ssh_public_key = "ssh-rsa AAAABHTmDMP6w=="

  openstack {
    lb_floating_network_id = "${data.openstack_networking_network_v2.fip_network_1.id}"
    network_id             = "${data.openstack_networking_network_v2.network_1.id}"
    lb_subnet_id           = "${data.openstack_networking_subnet_v2.subnet_1.id}"
    router_id              = "${data.openstack_networking_router_v2.router_1.id}"
    security_group_name    = "default"
  }

  node_pool {
    name              = "payload0"
    flavor            = "m1.xlarge_cpu"
    size              = 2
    availability_zone = "eu-de-1d"
    taints            = ["key=value:NoSchedule"]
    labels            = ["label=value"]
  }
}

resource "local_file" "kubeconfig" {
  sensitive_content = "${ccloud_kubernetes_v1.demo.kube_config_raw}"
  filename          = "kubeconfig"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Kubernikus client. If
  omitted, the `region` argument of the provider is used. Changing this forces
  a new resource to be created.

* `name` – (Required) Name of the cluster. Changing this forces a new resource
  to be created.

* `is_admin` - (Optional) Whether to create a Kubernetes cluster in the admin
  environment. Defaults to `false`. Changing this forces a new resource to be
  created.

* `advertise_address` - (Optional) The IP address on which to advertise the
  API server to members of the cluster. Defaults to `1.1.1.1`, which is a
  default virtual address, routed by the Kubernikus Wormhole tunnel on worker
  nodes. Changing this forces a new resource to be created.

* `cluster_cidr` - (Optional) CIDR Range for Pods in cluster. If not specified,
  generated automatically. Changing this forces a new resource to be created.

* `service_cidr` - (Optional) CIDR Range for Services in cluster. If not
  specified, generated automatically. Changing this forces a new resource to be
  created.

* `dns_address` - (Optional) The IP address of the `kube-dns` service. If not
  specified, generated automatically. Changing this forces a new resource to be
  created.

* `dns_domain` - (Optional) The DNS domain, served by the `kube-dns` service.
  If not specified, generated automatically. Changing this forces a new resource
  to be created.

* `ssh_public_key` - (Optional) The SSH public key, which should be used to
  authenticate the default SSH user (`core` for CoreOS images).

* `version` - (Optional) The version of the Kubernetes master.

* `node_pools` - (**Deprecated** - use `node_pool` instead)

* `node_pool` - (Optional) The list of Kubernetes node pools (worker pools).
  The `node_pool` object structure is documented below.

* `openstack` - (Optional) The advanced Openstack options. Required, when
  Kubernikus cannot automatically detect network settings, e.g. when multiple
  networks and routers are available. The `openstack` object structure is
  documented below.

The `node_pool` block supports:

* `name` - (Required) The unique node pool name. Changing this forces a new node
  pool to be created.

* `flavor` - (Required) The name of the desired flavor for the node pool compute
  instance. Changing this forces a new node pool to be created.

* `image` - (Optional) The name of the desired image for the node pool compute
  instance. If not specified, the default is used. Changing this forces a new
  node pool to be created.

* `size` - (Optional) The size of the node pool. Defaults to `0`.

* `availability_zone` - (Optional) The availability zone in which to create the
  the node pool. If not specified, detected automatically. Changing this forces
  a new node pool to be created.

* `taints` - (Optional) The list of Kubernetes node taints to be assigned on the
  node pool compute instance.

* `labels` - (Optional) The list of Kubernetes node labels to be assigned on the
  node pool compute instance.

* `config` - (Optional) Node pool extra options.

The node pool `config` block supports:

* `allow_reboot` - (Optional) Allow automatic drain and reboot of nodes. Enables
  OS updates. Required by security policy. Defaults to `true`.

* `allow_replace` - (Optional) Allow automatic drain and replacement of nodes.
  Enables Kubernetes upgrades. Defaults to `true`.

The `openstack` block supports:

* `lb_floating_network_id` - (Optional) The network ID of the floating IP pool.
  Specify this if there are multiple floating IP networks available. Changing
  this forces a new resource to be created.

* `network_id` - (Optional) The ID of the private network. Specify this if there
  are multiple private networks available. Changing this forces a new resource
  to be created.

* `lb_subnet_id` - (Optional) The private subnet ID of the loadbalancer. The
  subnet ID has to be a part of private network ID, set in the `network_id`.
  Specify this if there are multiple private subnets available. Changing this
  forces a new resource to be created.

* `project_id` - (Optional) The ID of the OpenStack project, where the
  Kubernikus cluster should be created. Available only withing the `is_admin`
  flag set to `true`. If not specified, detected automatically. Changing this
  forces a new resource to be created.

* `router_id` - (Optional) The ID of the network router. Specify this if there
  are multiple network routers available. Changing this forces a new resource to
  be created.

* `security_group_name` - (Optional) The security group name to associate with
  the compute instance, created within the node pool.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the cluster.
* `region` - See Argument Reference above.
* `name` – See Argument Reference above.
* `is_admin` - See Argument Reference above.
* `advertise_address` - See Argument Reference above.
* `cluster_cidr` - See Argument Reference above.
* `service_cidr` - See Argument Reference above.
* `dns_address` - See Argument Reference above.
* `dns_domain` - See Argument Reference above.
* `ssh_public_key` - See Argument Reference above.
* `node_pools` - See Argument Reference above.
* `node_pool` - See Argument Reference above.
* `openstack` - See Argument Reference above.
* `version` - See Argument Reference above.
* `phase` - The Kubernikus cluster current status. Can either be `Pending`,
  `Creating`, `Running`, `Terminating` or `Upgrading`.
* `wormhole` - The Wormhole tunnel server endpoint.
* `kube_config` - Contains the credentials block to the Kubernikus cluster.
* `kube_config_raw` - Contains the kubeconfig with credentials to the Kubernikus
  cluster.

The `kube_config` block exports the following:

* `host` - The Kubernetes cluster server host.

* `client_key` - Base64 encoded private key used by clients to authenticate to
  the Kubernetes cluster.

* `client_certificate` - Base64 encoded public certificate used by clients to
  authenticate to the Kubernetes cluster.

* `cluster_ca_certificate` - Base64 encoded public CA certificate used as the
  root of trust for the Kubernetes cluster.

* `username` - A username provided by the kubeconfig credentials.

* `not_after` - The credentials time validity bound, formatted as an RFC3339
  date string.

* `not_before` - The credentials time validity bound, formatted as an RFC3339
  date string.

-> **NOTE:** It is possible to use these credentials with
[the Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html)
like so:

```hcl
provider "kubernetes" {
  host                   = "${ccloud_kubernetes_v1.demo.kube_config.0.host}"
  client_certificate     = "${base64decode(ccloud_kubernetes_v1.demo.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(ccloud_kubernetes_v1.demo.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(ccloud_kubernetes_v1.demo.kube_config.0.cluster_ca_certificate)}"
}
```

## Timeouts

`ccloud_kubernetes_v1` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts)
configuration options:

* `create` - (Default `30 minutes`) How long to wait for the Kubernikus Cluster
  to be created.
* `update` - (Default `30 minutes`) How long to wait for the Kubernikus Cluster
  to be updated.
* `delete` - (Default `10 minutes`) How long to wait for the Kubernikus Cluster
  to be deleted.

## Import

Kubernikus Clusters can be imported using the `name` and `is_admin` flag
(`<name>/<is_admin>`), e.g.

```
$ terraform import ccloud_kubernetes_v1.demo demo/true
```

If the `is_admin` flag is omitted, it defaults to `false`.
