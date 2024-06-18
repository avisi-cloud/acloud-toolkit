---
date: 2024-06-18T12:42:26+02:00
title: "acloud-toolkit snapshots create"
displayName: "snapshots create"
slug: acloud-toolkit_snapshots_create
url: /references/acloud-toolkit/acloud-toolkit_snapshots_create/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 756
toc: true
---
## acloud-toolkit snapshots create

Create a snapshot of a Kubernetes PVC (persistent volume claim).

### Synopsis

This command creates a snapshot of a Kubernetes PVC, allowing you to capture a point-in-time copy of the data stored in the PVC. Snapshots can be used for data backup, disaster recovery, and other purposes.

To create a snapshot, you need to provide the name of the PVC to snapshot, as well as a name for the snapshot. You can also specify a namespace if the PVC is not in the current namespace context. If no snapshot class is specified, the default snapshot class will be used.

```
acloud-toolkit snapshots create <snapshot> [flags]
```

### Examples

```

# Create a snapshot of the PVC "my-pvc" with the name "my-snapshot":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc

#Create a snapshot of the PVC "my-pvc" with the name "my-snapshot" in the namespace "my-namespace":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc --namespace=my-namespace
		
```

### Options

```
  -h, --help                    help for create
  -n, --namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
  -p, --pvc string              Name of the PVC to snapshot. (required)
  -s, --snapshot-class string   Name of the CSI volume snapshot class to use. Uses the default VolumeSnapshotClass by default
  -t, --timeout duration        Duration to wait for the created snapshot to be ready for use (default 1h0m0s)
```

### SEE ALSO

* [acloud-toolkit snapshots](/references/acloud-toolkit/acloud-toolkit_snapshots/)	 - snapshot for working with Kubernetes CSI snapshot

