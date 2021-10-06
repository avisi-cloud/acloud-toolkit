---
date: 2021-10-06T10:20:08+02:00
title: "acloud-toolkit snapshot restore"
displayName: "snapshot restore"
slug: acloud-toolkit_snapshot_restore
url: /references/acloud-toolkit/acloud-toolkit_snapshot_restore/
description: ""
lead: ""
draft: false
images: []
menu:
  docs:
    parent: "acloud-toolkit-ref"
weight: 757
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
      --source-namespace string         (default "default")
      --target-name string             
      --target-namespace string         (default "default")
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

