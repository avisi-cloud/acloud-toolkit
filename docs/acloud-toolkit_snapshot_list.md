---
date: 2022-05-12T16:20:23+02:00
title: "acloud-toolkit snapshot list"
displayName: "snapshot list"
slug: acloud-toolkit_snapshot_list
url: /references/acloud-toolkit/acloud-toolkit_snapshot_list/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
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
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

