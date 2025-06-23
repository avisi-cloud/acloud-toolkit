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

# Create snapshots for all PVCs in the namespace "my-namespace":
acloud-toolkit snapshot create --all --namespace=my-namespace

# Create snapshots for all PVCs in the namespace "my-namespace" with a prefix "backup":
acloud-toolkit snapshot create --all --namespace=my-namespace --prefix=backup
		
```

### Options

```
      --all                     Create snapshots for all PVCs in the namespace, and use pvc name as snapshot name
      --concurrent-limit int    Maximum number of concurrent snapshot creation operations (default 10)
  -h, --help                    help for create
  -n, --namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
      --prefix -                Add a prefix seperated by - to the snapshot name when using --all
  -p, --pvc string              Name of the PVC to snapshot. (required)
  -s, --snapshot-class string   Name of the CSI volume snapshot class to use. Uses the default VolumeSnapshotClass by default
  -t, --timeout duration        Duration to wait for the created snapshot to be ready for use (default 1h0m0s)
```

### SEE ALSO

* [acloud-toolkit snapshots](acloud-toolkit_snapshots.md)	 - snapshot for working with Kubernetes CSI snapshot

