package migratevolume

type MigrationMode string

const (
	MigrationModeRsync  MigrationMode = "rsync"
	MigrationModeRclone MigrationMode = "rclone"
)

const (
	DefaultRSyncContainerImage  = "registry.avisi.cloud/library/rsync:v1"
	DefaultRCloneContainerImage = "rclone/rclone:1.66.0"
)
