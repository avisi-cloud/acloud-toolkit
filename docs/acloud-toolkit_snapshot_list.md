---
date: 2023-05-23T11:48:19+02:00
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
weight: 755
toc: true
---
## acloud-toolkit snapshot list

List all available CSI snapshots within the current namespace

### Synopsis

This command lists all available CSI snapshots within the current namespace. CSI snapshots are used to capture a point-in-time copy of a Kubernetes PVC, allowing you to preserve the data stored in the PVC for backup, disaster recovery, or other purposes.

By default, this command lists all snapshots in the current namespace. You can also specify a different namespace if needed.

```
acloud-toolkit snapshot list [flags]
```

### Examples

```

# List all available CSI snapshots within the current namespace:
acloud-toolkit snapshot list

# List all available CSI snapshots within the "my-namespace" namespace:
acloud-toolkit snapshot list --namespace=my-namespace

# List all available CSI snapshots within all namespaces:
acloud-toolkit snapshot list -A
		
```

### Options

```
  -A, --all-namespaces     return results for all namespaces
  -h, --help               help for list
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

