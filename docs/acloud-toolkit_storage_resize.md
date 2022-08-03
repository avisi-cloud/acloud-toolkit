---
date: 2022-10-17T15:00:02+02:00
title: "acloud-toolkit storage resize"
displayName: "storage resize"
slug: acloud-toolkit_storage_resize
url: /references/acloud-toolkit/acloud-toolkit_storage_resize/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 751
toc: true
---
## acloud-toolkit storage resize

resize adjusts the volume size of the pvc

### Synopsis

resize adjusts the volume size of the pvc

```
acloud-toolkit storage resize <persistent-volume-claim> [flags]
```

### Options

```
  -h, --help               help for resize
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
  -p, --pvc string         Name of the persistent volume to snapshot
      --size string        New size. Example: 10G
```

### SEE ALSO

* [acloud-toolkit storage](/references/acloud-toolkit/acloud-toolkit_storage/)	 - storage for working with Kubernetes CSI

