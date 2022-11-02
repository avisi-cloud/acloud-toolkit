---
date: 2022-11-02T21:41:26+01:00
title: "acloud-toolkit maintenance node-upgrade"
displayName: "maintenance node-upgrade"
slug: acloud-toolkit_maintenance_node-upgrade
url: /references/acloud-toolkit/acloud-toolkit_maintenance_node-upgrade/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 757
toc: true
---
## acloud-toolkit maintenance node-upgrade

upgrade a kubernetes node within an Avisi Cloud Kubernetes cluster (Bring Your Own Node only)

### Synopsis

Upgrade a kubernetes node within an Avisi Cloud Kubernetes cluster that is running with Bring Your Own Node enabled.

This command will upgrade both the Container Runtime and Kubelet version of a node.


```
acloud-toolkit maintenance node-upgrade <node> [flags]
```

### Examples

```

acloud-toolkit maintenance node-upgrade mynode
		
```

### Options

```
  -A, --all                upgrade all nodes within the cluster
  -h, --help               help for node-upgrade
      --timeout duration   The length of time to wait before giving up, zero means infinite
```

### SEE ALSO

* [acloud-toolkit maintenance](/references/acloud-toolkit/acloud-toolkit_maintenance/)	 - Perform maintenance actions on Kubernetes clusters

