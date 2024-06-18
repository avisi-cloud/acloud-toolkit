---
date: 2024-06-18T12:42:26+02:00
title: "acloud-toolkit volumes list"
displayName: "volumes list"
slug: acloud-toolkit_volumes_list
url: /references/acloud-toolkit/acloud-toolkit_volumes_list/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 749
toc: true
---
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

* [acloud-toolkit volumes](/references/acloud-toolkit/acloud-toolkit_volumes/)	 - Various commands for working with Kubernetes CSI volumes

