---
date: 2022-11-02T21:41:26+01:00
title: "acloud-toolkit maintenance drain"
displayName: "maintenance drain"
slug: acloud-toolkit_maintenance_drain
url: /references/acloud-toolkit/acloud-toolkit_maintenance_drain/
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
## acloud-toolkit maintenance drain

drain a kubernetes node

### Synopsis

drain a kubernetes node with additional options not supported by kubectl

```
acloud-toolkit maintenance drain <node> [flags]
```

### Examples

```

acloud-toolkit maintenance drain mynode --namespace-only default
		
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

* [acloud-toolkit maintenance](/references/acloud-toolkit/acloud-toolkit_maintenance/)	 - Perform maintenance actions on Kubernetes clusters

