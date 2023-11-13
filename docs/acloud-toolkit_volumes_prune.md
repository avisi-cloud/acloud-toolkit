---
date: 2023-11-13T16:20:39+01:00
title: "acloud-toolkit volumes prune"
displayName: "volumes prune"
slug: acloud-toolkit_volumes_prune
url: /references/acloud-toolkit/acloud-toolkit_volumes_prune/
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
## acloud-toolkit volumes prune

Prune removes any unused and released persistent volumes

### Synopsis

The 'prune' command removes any released persistent volumes.

```
acloud-toolkit volumes prune <persistent-volume-claim> [flags]
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

* [acloud-toolkit volumes](/references/acloud-toolkit/acloud-toolkit_volumes/)	 - Various commands for working with Kubernetes CSI volumes

