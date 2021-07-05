// Copyright 2019 Thomas Kooi

package app

import (
	"os"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/cmd/volume-migrator/app/cmd"
)

// Execute runs the cikcli application
func Execute() error {
	cmd := cmd.NewCikCmd(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
