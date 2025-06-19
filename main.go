package main

import (
	"github.com/avisi-cloud/acloud-toolkit/cmd"
	versionpkg "github.com/avisi-cloud/acloud-toolkit/pkg/version"
)

// these must be set by the compiler using LDFLAGS
// -X main.version= -X main.commit= -X main.date= -X main.builtBy=
var (
	version string
	commit  string
	date    string
	builtBy string
)

func main() {
	cmd.Execute()
}

func init() {
	versionpkg.Init(version, commit, date, builtBy)
}
