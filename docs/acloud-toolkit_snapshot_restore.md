---
date: 2022-10-17T15:00:02+02:00
title: "acloud-toolkit snapshot restore"
displayName: "snapshot restore"
slug: acloud-toolkit_snapshot_restore
url: /references/acloud-toolkit/acloud-toolkit_snapshot_restore/
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
## acloud-toolkit snapshot restore

Restore a snapshot

### Synopsis

restore a snapshot

```
acloud-toolkit snapshot restore <snapshot> [flags]
```

### Examples

```

acloud-toolkit snapshot restore my-snapshot --target-name my-pvc --restore-storage-class ebs-restore
		
```

### Options

```
  -h, --help                           help for restore
      --restore-storage-class string    (default "ebs-restore")
      --source-namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
      --target-name string             
      --target-namespace string        
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

