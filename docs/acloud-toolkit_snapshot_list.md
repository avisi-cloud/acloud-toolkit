---
date: 2022-04-12T15:34:08+02:00
title: "acloud-toolkit snapshot list"
displayName: "snapshot list"
slug: acloud-toolkit_snapshot_list
url: /references/acloud-toolkit/acloud-toolkit_snapshot_list/
description: ""
lead: ""
draft: false
images: []
menu:
  docs:
    parent: "acloud-toolkit-ref"
weight: 756
toc: true
---
## acloud-toolkit snapshot list

List CSI snapshots within the namespace

### Synopsis

List all available CSI snapshots within the namespace

```
acloud-toolkit snapshot list [flags]
```

### Examples

```

acloud-toolkit snapshot list
		
```

### Options

```
  -A, --all-namespaces     return results for all namespaces
  -h, --help               help for list
  -n, --namespace string   return snapshots from a specific namespace. Default is the configured namespace in your kubecontext.
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

