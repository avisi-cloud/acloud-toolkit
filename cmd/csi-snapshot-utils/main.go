package main

import (
	"gitlab.avisi.cloud/ame/csi-snapshot-utils/cmd/csi-snapshot-utils/app"
	versionpkg "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/version"
)

var (
	version string
	commit  string
	branch  string
)

func main() {
	versionpkg.SetCommit(commit)
	versionpkg.SetBranch(branch)
	versionpkg.SetVersion(version)

	app.Execute()
}
