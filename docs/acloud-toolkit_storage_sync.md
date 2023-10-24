---
date: 2023-10-08T20:36:57+02:00
title: "acloud-toolkit storage sync"
displayName: "storage sync"
slug: acloud-toolkit_storage_sync
url: /references/acloud-toolkit/acloud-toolkit_storage_sync/
description: ""
lead: ""
draft: false
images: []
menu:
    references:
        parent: "acloud-toolkit-ref"
weight: 747
toc: true
---

## acloud-toolkit storage sync

Sync a volume to another existing volume, or create a new one

### Synopsis

The `acloud-toolkit sync` command uses rsync under the hood to copy/sync all file contents from one volume to another. By default `acloud-toolkit sync` requires two existing persistentvolumeclaims to be specified. Optionally, a new PVC can be created using the `--create-pvc` and `--storageclass` flags.

```
acloud-toolkit storage sync [flags]
```

### Examples

#### Sync two existing volumes

```sh
acloud-toolkit storage sync --source-pvc pvc-a --target-pvc pvc-c
```

#### Sync and create a new volume

When specifying the `--create-pvc` flag, a new PVC will be created. The `--storageclass` flag must also be specified when creating a new PVC. The name of the PVC will be the same as the `--target-pvc` flag.

```sh
acloud-toolkit storage sync --source-pvc pvc-a --target-pvc pvc-c --create-pvc --storageclass nfs
```

#### Passing rsync options

After the `--` separator, any arguments will be passed to the rsync command. For example, to perform a dry-run:

```sh
acloud-toolkit storage sync --source-pvc pvc-a --target-pvc pvc-c -- --dry-run
```

### Including and excluding files using rsync options

Rsync can include or exclude files from the sync. For example, to exclude all files in the `/tmp` directory:

```sh
acloud-toolkit storage sync --source-pvc pvc-a --target-pvc pvc-c -- --exclude /tmp
```

Refer to the rsync documentation for more information on how to use its options.

### Debugging

Use the `--retain-job` flag to retain the job after completion. This will allow you to inspect the logs of the job before it is deleted.

```sh
acloud-toolkit storage sync --source-pvc pvc-a --target-pvc pvc-c  --retain-job -- --dry-run
```

### Options

```
      --create-pvc            create a new PVC if the target PVC does not exist
  -h, --help                  help for sync
  -n, --namespace string      namespace where the sync job will be executed
      --retain-job            retain the job after completion (default true)
      --source-pvc string     name of the source persitentvolumeclaim
      --storageclass string   name of the storageclass to use for the new PVC
      --target-pvc string     name of the target persitentvolumeclaim
  -t, --timeout int32         timeout of the context in minutes (default 60)
```

### SEE ALSO

-   [acloud-toolkit storage](/references/acloud-toolkit/acloud-toolkit_storage/) - storage for working with Kubernetes CSI
