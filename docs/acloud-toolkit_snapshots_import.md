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

* [acloud-toolkit snapshots](acloud-toolkit_snapshots.md)	 - snapshot for working with Kubernetes CSI snapshot

