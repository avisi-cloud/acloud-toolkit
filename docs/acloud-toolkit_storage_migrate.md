---
date: 2022-11-02T21:41:26+01:00
title: "acloud-toolkit storage migrate"
displayName: "storage migrate"
slug: acloud-toolkit_storage_migrate
url: /references/acloud-toolkit/acloud-toolkit_storage_migrate/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 748
toc: true
---
## acloud-toolkit storage migrate

Migrate a volume to another storage class

### Synopsis

Migrate a volume to another storage class. This will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.

```
acloud-toolkit storage migrate [flags]
```

### Options

```
  -h, --help                      help for migrate
      --new-size int              Use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC
  -p, --pvc string                name of the persitentvolumeclaim
  -s, --storageClass string       name of the new storageclass
  -n, --target-namespace string   Namespace where de migrate job will be executed (default "default")
  -t, --timeout int32             Timeout of the context in minutes (default 60)
```

### SEE ALSO

* [acloud-toolkit storage](/references/acloud-toolkit/acloud-toolkit_storage/)	 - storage for working with Kubernetes CSI

