---
date: 2024-05-06T11:06:16+02:00
title: "acloud-toolkit volumes batch-migrate"
displayName: "volumes batch-migrate"
slug: acloud-toolkit_volumes_batch-migrate
url: /references/acloud-toolkit/acloud-toolkit_volumes_batch-migrate/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 750
toc: true
---
## acloud-toolkit volumes batch-migrate

Batch migrate all volumes within a namespace to another storage class

### Synopsis

Batch migrate all volumes from a source storage class within a namespace to another storage class. For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.

```
acloud-toolkit volumes batch-migrate [flags]
```

### Options

```
      --dry-run                       Perform a dry run of the batch migrate
  -h, --help                          help for batch-migrate
  -s, --source-storage-class string   name of the source storageclass
  -n, --target-namespace string       Namespace where the migrate job will be executed (default "default")
  -t, --target-storage-class string   name of the target storageclass
      --timeout int32                 Timeout of the context in minutes (default 300)
```

### SEE ALSO

* [acloud-toolkit volumes](/references/acloud-toolkit/acloud-toolkit_volumes/)	 - Various commands for working with Kubernetes CSI volumes

