---
date: 2022-10-17T15:00:02+02:00
title: "acloud-toolkit snapshot create"
displayName: "snapshot create"
slug: acloud-toolkit_snapshot_create
url: /references/acloud-toolkit/acloud-toolkit_snapshot_create/
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
  -n, --namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
  -p, --pvc string              Name of the persistent volume to snapshot
  -s, --snapshot-class string   CSI volume snapshot class. If empty, use deafult volume snapshot class
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

