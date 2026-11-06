## acloud-toolkit snapshots delete

Delete all CSI snapshots created by a single snapshot create run

### Synopsis

This command deletes all CSI snapshots of a batch. Every `acloud-toolkit snapshot create` run labels the snapshots it creates with `acloud-toolkit.k8s.avisi.cloud/snapshot-batch=<batch-id>`, and prints the batch id when it finishes. The batch id of existing snapshots is shown in the Batch column of `acloud-toolkit snapshot list`.

Whether the underlying snapshot data is removed as well depends on the deletion policy of the VolumeSnapshotClass.

```
acloud-toolkit snapshots delete --batch-id <batch-id> [flags]
```

### Examples

```

# Show which snapshots of batch "20261106-142501-x7k2p" in the current namespace would be deleted:
acloud-toolkit snapshot delete --batch-id=20261106-142501-x7k2p --dry-run

# Delete all snapshots of batch "before-upgrade" in the namespace "my-namespace":
acloud-toolkit snapshot delete --batch-id=before-upgrade --namespace=my-namespace

# Delete all snapshots of batch "before-upgrade" in all namespaces:
acloud-toolkit snapshot delete --batch-id=before-upgrade -A
		
```

### Options

```
  -A, --all-namespaces             Delete the snapshots of the batch in all namespaces
      --batch-id snapshot create   Delete all snapshots of this batch (created by a single snapshot create run). (required)
      --dry-run                    Only show the snapshots that would be deleted
  -h, --help                       help for delete
  -n, --namespace string           If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit snapshots](acloud-toolkit_snapshots.md)	 - snapshot for working with Kubernetes CSI snapshot

