// Copyright 2019 Thomas Kooi

package app

import (
	"os"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/cmd/csi-utils/app/cmd"
)

// Execute runs the csi-utils application
func Execute() error {
	cmd := cmd.NewCSIUtilCmd(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
