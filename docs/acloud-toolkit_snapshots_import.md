---
date: 2024-06-18T12:42:26+02:00
title: "acloud-toolkit snapshots import"
displayName: "snapshots import"
slug: acloud-toolkit_snapshots_import
url: /references/acloud-toolkit/acloud-toolkit_snapshots_import/
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
## acloud-toolkit snapshots import

Import a raw snapshot ID into a CSI snapshot.

### Synopsis

This command creates Kubernetes CSI snapshot resources using a snapshot ID from the backend storage, for example AWS EBS, or Ceph RBD.
		

```
acloud-toolkit snapshots import <snapshot> [flags]
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

* [acloud-toolkit snapshots](/references/acloud-toolkit/acloud-toolkit_snapshots/)	 - snapshot for working with Kubernetes CSI snapshot

