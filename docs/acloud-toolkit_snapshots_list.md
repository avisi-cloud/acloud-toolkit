## acloud-toolkit snapshots list

List all available CSI snapshots within the current namespace

### Synopsis

This command lists all available CSI snapshots within the current namespace. CSI snapshots are used to capture a point-in-time copy of a Kubernetes PVC, allowing you to preserve the data stored in the PVC for backup, disaster recovery, or other purposes.

By default, this command lists all snapshots in the current namespace. You can also specify a different namespace if needed.

```
acloud-toolkit snapshots list [flags]
```

### Examples

```

# List all available CSI snapshots within the current namespace:
acloud-toolkit snapshot list

# List all available CSI snapshots within the "my-namespace" namespace:
acloud-toolkit snapshot list --namespace=my-namespace

# List all available CSI snapshots within all namespaces:
acloud-toolkit snapshot list -A
		
```

### Options

```
  -A, --all-namespaces     return results for all namespaces
  -S, --handles            show snapshot content handle (default true)
  -h, --help               help for list
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit snapshots](acloud-toolkit_snapshots.md)	 - snapshot for working with Kubernetes CSI snapshot

