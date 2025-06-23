## acloud-toolkit volumes list

List all persistent volumes in a Kubernetes cluster

### Synopsis

This command lists all CSI persistent volumes within the cluster. This command allows you to list and filter persistent volumes based on various criteria, making it easier to inspect and manage your storage resources.

```
acloud-toolkit volumes list [flags]
```

### Examples

```
# List all persistent volumes of a specific storage class within the cluster:
acloud-toolkit volumes list -s my-storage-class

# List all unattached persistent volumes:
acloud-toolkit volumes list --unattached-only

# List all unattached CSI persistent volumes:
acloud-toolkit volumes list --unattached-only --csi-only

```

### Options

```
      --csi-only               show CSI persistent volumes only
  -h, --help                   help for list
  -s, --storage-class string   run for storage class. Will use default storage class if left empty
      --unattached-only        show unattached persistent volumes only
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

