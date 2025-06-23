## acloud-toolkit snapshots restore

Restore a Kubernetes PVC from a CSI snapshot.

### Synopsis

This command restores a Kubernetes PVC from a CSI snapshot. To restore a PVC, you need to provide the name of the snapshot, the name of the PVC to restore to, and the namespace of the target PVC. You can also specify a different namespace for the snapshot if needed.

By default, this command restores the PVC to the default storage class installed within the cluster. You can specify a different storage class if needed by using the --restore-storage-class option. Please note that this command requires the volume mode to be set to "Immediate".
		

```
acloud-toolkit snapshots restore <snapshot> [flags]
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
      --restore-storage-class string   
      --source-namespace string        If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
  -t, --timeout duration               Duration to wait for the restored snapshot to complete (default 10m0s)
```

### SEE ALSO

* [acloud-toolkit snapshots](acloud-toolkit_snapshots.md)	 - snapshot for working with Kubernetes CSI snapshot

