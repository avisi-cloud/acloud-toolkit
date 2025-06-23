## acloud-toolkit volumes resize

Resize adjusts the volume size of a persistent volume claim

### Synopsis

The 'resize' command adjusts the size of a persistent volume claim (PVC). The command takes a PVC name as input along with an optional namespace parameter and a new size in gigabytes.

```
acloud-toolkit volumes resize <persistent-volume-claim> [flags]
```

### Examples

```

# Resize a PVC named 'data' in the default namespace to 20 gigabytes
acloud-toolkit storage resize data --size 20G

# Resize a PVC named 'data' in the 'prod' namespace to 50 gigabytes
acloud-toolkit storage resize data --namespace prod --size 50G	  

```

### Options

```
  -h, --help               help for resize
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
      --size string        New size. Example: 10G
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

