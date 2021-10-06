---
date: 2021-10-06T10:20:08+02:00
title: "acloud-toolkit snapshot create"
displayName: "snapshot create"
slug: acloud-toolkit_snapshot_create
url: /references/acloud-toolkit/acloud-toolkit_snapshot_create/
description: ""
lead: ""
draft: false
images: []
menu:
  docs:
    parent: "acloud-toolkit-ref"
weight: 759
toc: true
---
## acloud-toolkit snapshot create

create creates a snapshot for a pvc

### Synopsis

create creates a snapshot for a pvc

```
acloud-toolkit snapshot create <snapshot> [flags]
```

### Examples

```

acloud-toolkit snapshot create my-snapshot --pvc my-pvc
		
```

### Options

```
  -h, --help                    help for create
  -n, --namespace string        Namespace of the PVC. Snapshot will be created within this namespace as well (default "default")
  -p, --pvc string              Name of the persistent volume to snapshot
  -s, --snapshot-class string   CSI snapshot class (default "csi-aws-vsc")
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

