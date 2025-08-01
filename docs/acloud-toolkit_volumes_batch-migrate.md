## acloud-toolkit volumes batch-migrate

Batch migrate all volumes within a namespace to another storage class

### Synopsis

Batch migrate all volumes from a source storage class within a namespace to another storage class.
For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume(s) will remain available within the cluster.

Batch migrate supports both rclone and rsync migration modes. The default mode is rsync.
- When using rsync, by default it uses the --archive flag. It will preserve all file permissions, timestamps, and ownerships.
- When using rclone a copy command is used. Use --metadata flag to preserve metadata.

It is recommended to utilize the migration-flag option to pass additional flags to the migration tool, such as --checksum or others and optmize the migration job for your specific use case.
		

```
acloud-toolkit volumes batch-migrate [flags]
```

### Examples

```
# Migrate all rbd volumes to rbd-new using rclone:
acloud-toolkit storage batch-migrate -s rbd -t rbd-new -m rclone

# Add additional flags to rclone:
acloud-toolkit storage batch-migrate -s rbd -t rbd-new -m rclone -f '--multi-thread-streams=8 --metadata'

```

### Options

```
      --dry-run                       Perform a dry run of the batch migrate
  -h, --help                          help for batch-migrate
  -f, --migration-flags string        Additional flags to pass to the migration tool
  -m, --migration-mode string         Migration mode to use. Options: rsync, rclone (default "rsync")
      --node-selector strings         comma separated list of node labels used for nodeSelector of the migration job
      --preserve-metadata             Preserve the original metadata of the PVC
      --rclone-image string           Image used for the rclone migration tool (default "rclone/rclone:1.66.0")
      --rsync-image string            Image used for the rsync migration tool (default "registry.avisi.cloud/library/rsync:v1")
  -s, --source-storage-class string   name of the source storageclass
  -n, --target-namespace string       Namespace where the migrate job will be executed
  -t, --target-storage-class string   name of the target storageclass
      --timeout int32                 Timeout of the context in minutes (default 300)
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

