package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app"
)

func main() {
	cmd := app.NewACloudToolKitCmd(nil, nil, nil)
	err := doc.GenMarkdownTree(cmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
