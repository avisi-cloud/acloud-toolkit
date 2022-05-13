package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra/doc"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app"
)

const fmTemplate = `---
date: %s
title: "%s"
displayName: "%s"
slug: %s
url: %s
description: ""
lead: ""
draft: false
images: []
menu:
  references:
    parent: "%s-ref"
weight: %d
toc: true
---
`

func main() {
	cmd := app.NewACloudToolKitCmd(nil, nil, nil)

	weight := 760
	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		displayName := strings.TrimPrefix(base, cmd.Name()+"_")
		url := "/references/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
		weight--
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), strings.Replace(displayName, "_", " ", -1), base, url, cmd.Name(), weight)
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/references/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
	}

	cmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTreeCustom(cmd, "./docs", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
