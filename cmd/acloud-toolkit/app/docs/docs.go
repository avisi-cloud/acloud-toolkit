package docs

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	browserURL = "https://docs.avisi.cloud/docs/cli/acloud-toolkit/overview"
)

// NewOpenDocs returns the Cobra version sub command
func NewOpenDocs() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "online-docs",
		Short: "Open the online documentation for acloud-toolkit",
		Long:  `Open the online documentation for acloud-toolkit. This will open a new tab in your default browser for our docs.avisi.cloud website.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Opening %s ...\n", browserURL)
			return openBrowser(browserURL)
		},
	}

	return versionCmd
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
