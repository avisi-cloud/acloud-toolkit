## acloud-toolkit volumes migrate

Migrate the filesystem on a persistent volume to another storage class

### Synopsis

Migrate the filesystem on a persistent volume to another storage class.
This will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.

Migrate supports both rclone and rsync migration modes. The default mode is rsync.
- When using rsync, by default it uses the --archive flag. It will preserve all file permissions, timestamps, and ownerships.
- When using rclone a copy command is used. Use --metadata flag to preserve metadata.

It is recommended to utilize the migration-flag option to pass additional flags to the migration tool, such as --checksum or others and optmize the migration job for your specific use case.


```
acloud-toolkit volumes migrate [flags]
```

### Examples

```
# Migrate volumes from one storage class to another, in this example migrate pvc `app-data` to `gp2`
acloud-toolkit volumes migrate -s gp2 --pvc app-data -n default
```

### Options

```
  -h, --help                      help for migrate
  -f, --migration-flags string    Additional flags to pass to the migration tool
  -m, --migration-mode string     Migration mode to use. Options: rsync, rclone. Default is rsync with rclone being newly introduced (default "rsync")
      --new-size int              Use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC
      --node-selector strings     comma separated list of node labels used for nodeSelector of the migration job
  -p, --pvc string                name of the persitentvolumeclaim
      --rclone-image string       Image used for the rclone migration tool (default "rclone/rclone:1.66.0")
      --rsync-image string        Image used for the rsync migration tool (default "registry.avisi.cloud/library/rsync:v1")
  -s, --storageClass string       name of the new storageclass
  -n, --target-namespace string   Namespace where the volume migrate job will be executed
  -t, --timeout int32             Timeout of the context in minutes (default 300)
```

### SEE ALSO

* [acloud-toolkit volumes](acloud-toolkit_volumes.md)	 - Various commands for working with Kubernetes CSI volumes

