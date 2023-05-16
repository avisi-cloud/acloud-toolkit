---
date: 2023-05-23T11:48:19+02:00
title: "acloud-toolkit nodes drain"
displayName: "nodes drain"
slug: acloud-toolkit_nodes_drain
url: /references/acloud-toolkit/acloud-toolkit_nodes_drain/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 759
toc: true
---
## acloud-toolkit nodes drain

drain a kubernetes node with additional options not supported by kubectl

### Synopsis

The acloud-toolkit nodes drain command is a CLI tool that allows you to gracefully remove a Kubernetes node from service, ensuring that all workloads running on the node are rescheduled to other nodes in the cluster before the node is taken offline for maintenance or other purposes. This command provides additional options that are not supported by the standard kubectl drain command.

```
acloud-toolkit nodes drain <node> [flags]
```

### Examples

```
# Drain a node without uncordoning it afterwards:
acloud-toolkit nodes drain mynode

# Drain a node and uncordon it afterwards:
acloud-toolkit nodes drain mynode --uncordon

# Drain a node and only evict pods from a specific namespace:
acloud-toolkit nodes drain mynode --namespace mynamespace

# Drain a node and only evict stateless workloads:
acloud-toolkit nodes drain mynode --ignore-statefulset-pods

# Drain a node and set the grace period to 120 seconds:
acloud-toolkit nodes drain mynode --grace-period 120

# Drain a node and set the timeout to 10 minutes:
acloud-toolkit nodes drain mynode --timeout 10m

```

### Options

```
      --grace-period int          Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used (default 60)
  -h, --help                      help for drain
      --ignore-statefulset-pods   do not drain statefulset pods
  -n, --namespace string          drain pods from a specific namespace only. Default is the configured namespace in your kubecontext.
      --timeout duration          The length of time to wait before giving up, zero means infinite
      --uncordon                  uncordon nodes after running the drain command
```

### SEE ALSO

* [acloud-toolkit nodes](/references/acloud-toolkit/acloud-toolkit_nodes/)	 - Perform actions on Kubernetes cluster nodes

