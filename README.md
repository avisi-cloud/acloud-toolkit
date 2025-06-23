# acloud-toolkit

[![Build Status](https://github.com/avisi-cloud/acloud-toolkit/actions/workflows/build.yml/badge.svg)](https://github.com/avisi-cloud/acloud-toolkit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/avisi-cloud/acloud-toolkit)](https://goreportcard.com/report/github.com/avisi-cloud/acloud-toolkit)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Our powerful CLI toolkit for Kubernetes automation and repetitive tasks, specializing in CSI snapshot management, volume migration, and storage.

## ✨ Features

- **Snapshot Management**: Create, restore, import, and list CSI snapshots
- **Volume Migration**: Migrate volumes between storage classes
- **Volume Sync**: Synchronize data between persistent volumes
- **Volume Resize**: Easily resize persistent volumes
- **Storage Cleanup**: Prune orphaned volumes and snapshots

## 🚀 Quick Start

### Installation

#### From Homebrew (macOS/Linux) (Recommended)
```bash
brew install avisi-cloud/tools/acloud-toolkit --cask
```

#### From Release

Download the latest release from the [Releases page](https://github.com/avisi-cloud/acloud-toolkit/releases) and extract it. Then copy the binary to the desired location, e.g. `/usr/local/bin`.

#### From Source
```bash
git clone https://github.com/avisi-cloud/acloud-toolkit.git
cd acloud-toolkit
make build
sudo cp bin/acloud-toolkit /usr/local/bin/acloud-toolkit
```

### Verify Installation
```bash
acloud-toolkit version
```

## 📖 Usage Examples

### Snapshot Operations

#### Create a snapshot
```bash
# Create snapshot from a PVC
acloud-toolkit snapshot create my-snapshot --pvc my-pvc

# Create snapshots for all PVCs in the namespace "my-namespace" with a prefix "backup":
acloud-toolkit snapshot create --all --namespace my-namespace --prefix backup
```

#### Restore from snapshot
```bash
# Restore to new PVC
acloud-toolkit snapshot restore my-snapshot \
  --restore-pvc-name my-pvc \
  --restore-storage-class ebs-restore
```

#### Import external snapshots
```bash
# Import AWS EBS snapshot
acloud-toolkit snapshot import \
  snap-1234567890abcdef0 \
  --name my-imported-snapshot
```

### Storage Management

#### Prune orphaned resources
```bash
acloud-toolkit volumes prune                    # Preview what will be deleted
acloud-toolkit volumes prune --dry-run=false    # Execute cleanup
```

## 🤝 Contributing

### Development Setup
```bash
git clone https://github.com/avisi-cloud/acloud-toolkit.git
cd acloud-toolkit
make test
make build
```

### Running Tests
```bash
make test          # Unit tests
make lint          # Code linting
make race          # Race condition detection
```

### Generate Documentation
```bash
go run tools/docs.go
```

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
