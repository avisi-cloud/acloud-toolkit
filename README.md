# CSI snapshot utils

CLI tooling for working with CSI snapshots.

This is based on the [manual documentation](https://insight.avisi.nl/confluence/display/AME/how-to+restore+a+snapshot+to+a+new+namespace) for restoring CSI snapshots.

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

### Create snapshot

```bash
./bin/csi-snapshot-utils create-snapshot --pvc jira-jira-0 --snapshot-name jira-$(date "+%F-%H%M%S") -n jira
```

### Examples

```
./bin/csi-snapshot-utils restore --snapshot-name=jira-2021-04-14-110521 --target-name=jira-jira-0 --target-namespace=jira-acc --source-namespace=jira
./bin/csi-snapshot-utils restore --snapshot-name=jira-db-2021-04-14-110635 --target-name=data-jira-database-postgresql-0 --target-namespace=jira-acc --source-namespace=jira
#Migrate a volume to a larger volume, as resizing a volume is not support yet. To do this, the pod using this volume should not be running!
./bin/csi-utils migrate -p mysql-volume-sugarcrm-0 -s rbd -n brn --new-size 15360
```
