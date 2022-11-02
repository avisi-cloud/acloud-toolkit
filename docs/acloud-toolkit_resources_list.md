---
date: 2022-11-02T21:41:26+01:00
title: "acloud-toolkit resources list"
displayName: "resources list"
slug: acloud-toolkit_resources_list
url: /references/acloud-toolkit/acloud-toolkit_resources_list/
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "acloud-toolkit-ref"
weight: 755
toc: true
---
## acloud-toolkit resources list

resources list displays resource requests and limits within the namespace agregrated by deployment or statefulset (experimental)

### Synopsis

resources list displays resource requests and limits within the namespace agregrated by deployment or statefulset. This is experimental functionality. The displayed CPU and Memory are a sum of all resource requests and limits for each container within a pod.

For example, a pod with two containers of each 100Mi Memory limits will show 200Mi for Memory limits.

```
acloud-toolkit resources list [flags]
```

### Examples

```

# List resource limits within a specific namespace

❯ ./bin/acloud-toolkit resources list -n nginx-ingress
NAMESPACE               TYPE            NAME                            CONTAINERS      REPLICAS        MEMORY 
nginx-ingress           Deployment      ingress-nginx-controller        1               2               150Mi 

# List resource limits for all namespaces

❯ ./bin/acloud-toolkit resources list -A
...

		
```

### Options

```
  -A, --all-namespaces     return results for all namespaces
  -h, --help               help for list
  -n, --namespace string   If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context
```

### SEE ALSO

* [acloud-toolkit resources](/references/acloud-toolkit/acloud-toolkit_resources/)	 - Gather insight into resource usage and limits within a namespace (experimental)

