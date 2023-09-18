---
date: 2023-09-18T14:42:04+02:00
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
weight: 753
toc: true
---
## acloud-toolkit snapshot restore

Restore a Kubernetes PVC from a CSI snapshot.

### Synopsis

This command restores a Kubernetes PVC from a CSI snapshot. To restore a PVC, you need to provide the name of the snapshot, the name of the PVC to restore to, and the namespace of the target PVC. You can also specify a different namespace for the snapshot if needed.

By default, this command restores the PVC to the default storage class installed within the cluster. You can specify a different storage class if needed by using the --restore-storage-class option. Please note that this command requires the volume mode to be set to "Immediate".
		

```
acloud-toolkit snapshot restore <snapshot> [flags]
```

### Examples

```

acloud-toolkit snapshot restore my-snapshot --restore-pvc-name my-pvc --restore-storage-class ebs-restore
		
```

### Options

```
  -h, --help                           help for restore
      --restore-pvc-name string        
      --restore-pvc-namespace string   
      --restore-storage-class string    (default "ebs-restore")
      --source-namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit snapshot](/references/acloud-toolkit/acloud-toolkit_snapshot/)	 - snapshot for working with Kubernetes CSI snapshot

