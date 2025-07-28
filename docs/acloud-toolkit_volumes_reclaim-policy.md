## acloud-toolkit volumes reclaim-policy

Change the reclaim policy of a persistent volume

### Synopsis

Change the reclaim policy of a persistent volume. You can specify either a persistent volume name directly or a persistent volume claim name.

The reclaim policy determines what happens to the underlying storage when the persistent volume is released:
- Retain: The volume will be retained and must be manually reclaimed
- Delete: The volume will be automatically deleted when released
- Recycle: The volume will be scrubbed and made available for new claims (deprecated)

```
acloud-toolkit volumes reclaim-policy [flags]
```

### Examples

```

# Set reclaim policy to Retain for a specific PV
acloud-toolkit volumes reclaim-policy --pv my-pv --policy Retain

# Set reclaim policy to Delete for a PV via PVC using current namespace from kubeconfig
acloud-toolkit volumes reclaim-policy --pvc my-pvc --policy Delete

# Set reclaim policy to Retain for a PV via PVC in a specific namespace
acloud-toolkit volumes reclaim-policy --pvc data-pvc --namespace production --policy Retain

```

### Options

```
  -h, --help               help for reclaim-policy
  -n, --namespace string   namespace of the persistent volume claim (optional when using --pvc, defaults to current kubeconfig context)
  -p, --policy string      reclaim policy to set (Retain, Delete, Recycle)
      --pv string          name of the persistent volume
      --pvc string         name of the persistent volume claim
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

