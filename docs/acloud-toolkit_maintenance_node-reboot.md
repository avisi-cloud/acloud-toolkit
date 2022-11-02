---
date: 2022-11-02T21:41:26+01:00
title: "acloud-toolkit maintenance node-reboot"
displayName: "maintenance node-reboot"
slug: acloud-toolkit_maintenance_node-reboot
url: /references/acloud-toolkit/acloud-toolkit_maintenance_node-reboot/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 758
toc: true
---
## acloud-toolkit maintenance node-reboot

reboot a kubernetes node within an Avisi Cloud Kubernetes cluster if required

### Synopsis

reboot a kubernetes node within an Avisi Cloud Kubernetes cluster. This is only performed if the file '/var/run/reboot-required' is present on the host machine.


```
acloud-toolkit maintenance node-reboot <node> [flags]
```

### Examples

```

acloud-toolkit maintenance node-reboot mynode
		
```

### Options

```
  -A, --all                reboot all nodes within the cluster
  -h, --help               help for node-reboot
      --timeout duration   The length of time to wait before giving up, zero means infinite
```

### SEE ALSO

* [acloud-toolkit maintenance](/references/acloud-toolkit/acloud-toolkit_maintenance/)	 - Perform maintenance actions on Kubernetes clusters

