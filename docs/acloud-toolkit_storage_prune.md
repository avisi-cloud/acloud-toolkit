---
date: 2023-05-23T11:48:19+02:00
title: "acloud-toolkit storage prune"
displayName: "storage prune"
slug: acloud-toolkit_storage_prune
url: /references/acloud-toolkit/acloud-toolkit_storage_prune/
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
## acloud-toolkit storage prune

prune removes any unused and released persistent volumes

### Synopsis

The 'prune' command removes any released persistent volumes.

```
acloud-toolkit storage prune <persistent-volume-claim> [flags]
```

### Examples

```

# Prune all persistent volumes that are set to Released
acloud-toolkit storage prune

```

### Options

```
      --dry-run   Perform a dry run of volume prune (default true)
  -h, --help      help for prune
```

### SEE ALSO

* [acloud-toolkit storage](/references/acloud-toolkit/acloud-toolkit_storage/)	 - storage for working with Kubernetes CSI

