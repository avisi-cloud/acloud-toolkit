---
date: 2022-10-17T15:00:02+02:00
title: "acloud-toolkit storage batch-migrate"
displayName: "storage batch-migrate"
slug: acloud-toolkit_storage_batch-migrate
url: /references/acloud-toolkit/acloud-toolkit_storage_batch-migrate/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 753
toc: true
---
## acloud-toolkit storage batch-migrate

Batch migrate all volumes within a namespace to another storage class

### Synopsis

Batch migrate all volumes from a source storage class within a namespace to another storage class. For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.

```
acloud-toolkit storage batch-migrate [flags]
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

* [acloud-toolkit storage](/references/acloud-toolkit/acloud-toolkit_storage/)	 - storage for working with Kubernetes CSI

