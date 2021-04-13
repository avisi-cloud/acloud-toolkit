# CSI snapshot utils

CLI tooling for working with CSI snapshots

## Installation

```bash
{
    make build;
    cp bin/csi-snapshot-utils /usr/local/bin/csi-snapshot-utils;
    csi-snapshot-utils -h;
}
```

## Usage

```bash
$ ./bin/csi-snapshot-utils restore -h
restore a snapshot

Usage:
  csi-snapshot-utils restore [flags]

Flags:
  -h, --help                           help for restore
      --restore-storage-class string    (default "ebs-restore")
      --snapshot-name string           name of the snapshot
      --source-namespace string         (default "default")
      --target-name string             
      --target-namespace string         (default "default")
```

### Restore example

```bash
$ ./bin/csi-snapshot-utils restore --snapshot-name=test-snapshotgroup-1614584579 --target-name=test3 --target-namespace=default
using snapshot test-snapshotgroup-1614584579 for restoring
created PVC test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
PVC has volume pvc-dc23367b-f5b2-4fb7-9ddc-da0d864a7147...
deleted the PVC test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
removed the PV pvc-dc23367b-f5b2-4fb7-9ddc-da0d864a7147 claim ref to test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
created a new PVC test3 in namespace default...
```
