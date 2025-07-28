## acloud-toolkit volumes prune

Prune removes any unused and released persistent volumes

### Synopsis

The 'prune' command removes any released persistent volumes. By default it will run in dry-run mode, which will only show the volumes that would be pruned. Use the --dry-run=false flag to actually prune the volumes.

```
acloud-toolkit volumes prune <persistent-volume-claim> [flags]
```

### Examples

```

# See all persistent volumes that are set to Released that would be pruned
acloud-toolkit storage prune -A

# Prune all persistent volumes that are set to Released
acloud-toolkit storage prune -A --dry-run=false

# Prune all persistent volumes that are set to Released within a specific namespace
acloud-toolkit storage prune -n my-namespace --dry-run=false

```

### Options

```
  -A, --all                              Prune volumes from all namespaces
      --dry-run                          Perform a dry run of volume prune (default true)
  -h, --help                             help for prune
  -l, --label-selector string            Label selector to filter the volumes to prune
      --min-released-duration duration   Minimum duration since the volume was released
  -n, --namespace string                 Namespace to prune volumes from. Volume namespaces are cluster scoped, so the namespace is only used to filter the PVCs
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

