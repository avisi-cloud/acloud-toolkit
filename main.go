package main

import (
	"math/rand"
	"time"

	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app"
	versionpkg "gitlab.avisi.cloud/ame/acloud-toolkit/pkg/version"
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
	// make sure we have seed the rand package
	rand.Seed(time.Now().UnixNano())
	app.Execute()
}

func init() {
	versionpkg.Init(version, commit, date, builtBy)
}
