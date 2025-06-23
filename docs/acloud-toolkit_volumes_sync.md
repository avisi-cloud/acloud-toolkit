## acloud-toolkit volumes sync

Sync a volume to another existing volume, or create a new volume

### Synopsis

Sync a volume to another existing volume, or create a new volume. This will create a new PVC using the target storage class or use an existing one, and copy all file contents over to the new volume using rsync. The existing persistent volume and persistent volume claim will remain available in the cluster.

```
acloud-toolkit volumes sync [flags]
```

### Examples

```
# Sync two existing volumes
acloud-toolkit volumes sync --source-pvc pvc-a --target-pvc pvc-b

# Passing rsync options
# After the `--` separator, any arguments will be passed to the rsync command. For example, to perform a dry-run:
acloud-toolkit volumes sync --source-pvc pvc-a --target-pvc pvc-b -- --dry-run

# Including and excluding files using rsync options
# rsync can include or exclude files from the sync. For example, to exclude all files in the `/tmp` directory:
acloud-toolkit volumes sync --source-pvc pvc-a --target-pvc pvc-b -- --exclude /tmp

# Sync and create a new volume
# When specifying the `--create-pvc` flag, a new PVC will be created. The `--storageclass` flag must also be specified when creating a new PVC. The name of the PVC will be the same as the `--target-pvc` flag.
acloud-toolkit volumes sync --source-pvc pvc-a --target-pvc pvc-c --create-pvc --storageclass nfs

# Debug job by retaining it after completion and setting the ttl to 3600 seconds
acloud-toolkit volumes sync --source-pvc pvc-a --target-pvc pvc-b --ttl 3600 --retain-job -- --dry-run

```

### Options

```
      --create-pvc            create a new PVC if the target PVC does not exist
  -h, --help                  help for sync
  -n, --namespace string      namespace where the sync job will be executed
      --new-size int          use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC
      --retain-job            retain the job after completion
      --source-pvc string     name of the source persitentvolumeclaim
      --storageclass string   name of the storageclass to use for the new PVC (default "ebs-restore")
      --target-pvc string     name of the target persitentvolumeclaim
  -t, --timeout int32         timeout of the context in minutes (default 60)
      --ttl int32             time to live in seconds after the job has finished, requires --retain-job to be true (default 3600)
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

