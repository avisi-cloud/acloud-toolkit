---
date: 2023-09-18T14:42:04+02:00
title: "acloud-toolkit snapshot import"
displayName: "snapshot import"
slug: acloud-toolkit_snapshot_import
url: /references/acloud-toolkit/acloud-toolkit_snapshot_import/
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
## acloud-toolkit snapshot import

Import raw Snapshot ID into a CSI snapshot.

### Synopsis

This command creates Kubernetes CSI snapshot resources using a snapshot ID from the backend storage, for example AWS EBS, or Ceph RBD.
		

```
acloud-toolkit snapshot import <snapshot> [flags]
```

### Examples

```

acloud-toolkit snapshot import --name example snap-12345
		
```

### Options

```
  -h, --help                            help for import
      --name string                     name of the snapshot
      --namespace string                If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
      --snapshot-storage-class string   
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

