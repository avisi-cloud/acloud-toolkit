# CSI snapshot utils

CLI tooling for working with CSI snapshots.

This is based on the [manual documentation](https://insight.avisi.nl/confluence/display/AME/how-to+restore+a+snapshot+to+a+new+namespace) for restoring CSI snapshots.

## Installation

```bash
{
    make build;
    cp bin/acloud-toolkit /usr/local/bin/acloud-toolkit;
    acloud-toolkit -h;
}
```

## Usage

```bash
$ ./bin/acloud-toolkit restore -h
restore a snapshot

Usage:
  acloud-toolkit restore [flags]

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
$ ./bin/acloud-toolkit restore --snapshot-name=test-snapshotgroup-1614584579 --target-name=test3 --target-namespace=default
using snapshot test-snapshotgroup-1614584579 for restoring
created PVC test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
PVC has volume pvc-dc23367b-f5b2-4fb7-9ddc-da0d864a7147...
deleted the PVC test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
removed the PV pvc-dc23367b-f5b2-4fb7-9ddc-da0d864a7147 claim ref to test3-721bd6a7-6f75-471f-ad24-a99e063a6cf0...
created a new PVC test3 in namespace default...
```

### Create snapshot

```bash
./bin/acloud-toolkit create-snapshot --pvc jira-jira-0 --snapshot-name jira-$(date "+%F-%H%M%S") -n jira
```

### Examples

```
./bin/acloud-toolkit restore --snapshot-name=jira-2021-04-14-110521 --target-name=jira-jira-0 --target-namespace=jira-acc --source-namespace=jira
./bin/acloud-toolkit restore --snapshot-name=jira-db-2021-04-14-110635 --target-name=data-jira-database-postgresql-0 --target-namespace=jira-acc --source-namespace=jira
```


```
kubectl get pv --no-headers | grep Released|awk '{print $1}'|xargs kubectl  patch pv -p '{"spec":{"persistentVolumeReclaimPolicy":"Delete"}}'
```

kubectl get pv --no-headers | grep Available|awk '{print $1}'|xargs kubectl  patch pv -p '{"spec":{"persistentVolumeReclaimPolicy":"Delete"}}'

### Documentation

- [documentation](docs/acloud-toolkit.md)

Documentation is auto generated through running:
```bash
go run tools/docs.go
```

Please check out `docs/`.
