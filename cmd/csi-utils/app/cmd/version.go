// Copyright 2019 Thomas Kooi

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	branch  string
)

// NewVersionCmd returns the Cobra version sub command
func NewVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `version information`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("well this didn't do anything ... \n")
		},
	}

	return versionCmd
}
