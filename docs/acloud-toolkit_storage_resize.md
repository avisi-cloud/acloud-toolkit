---
date: 2023-09-18T14:42:04+02:00
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
weight: 748
toc: true
---
## acloud-toolkit storage resize

resize adjusts the volume size of a persistent volume claim

### Synopsis

The 'resize' command adjusts the size of a persistent volume claim (PVC). The command takes a PVC name as input along with an optional namespace parameter and a new size in gigabytes.

```
acloud-toolkit storage resize <persistent-volume-claim> [flags]
```

### Examples

```

# Resize a PVC named 'data' in the default namespace to 20 gigabytes
acloud-toolkit storage resize data --size 20G

# Resize a PVC named 'data' in the 'prod' namespace to 50 gigabytes
acloud-toolkit storage resize data --namespace prod --size 50G	  

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

