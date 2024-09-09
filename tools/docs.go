package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"

	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app"
)

const mdxTemplate = `---
title: "%s"
description: ""
---
`

const (
	docsPath = "./docs"
)

func main() {
	cmd := app.NewACloudToolKitCmd(nil, nil, nil)

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		finalName := strings.TrimPrefix(base, cmd.Name()+"_")
		return fmt.Sprintf(mdxTemplate, strings.Replace(finalName, "_", " ", -1))
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/docs/cli/" + cmd.Name() + "/commands/" + strings.ToLower(base) + "/"
	}

	if err := os.MkdirAll(docsPath, 0755); err != nil {
		log.Fatal(err)
	}
	cmd.DisableAutoGenTag = true
	if err := doc.GenMarkdownTreeCustom(cmd, docsPath, filePrepender, linkHandler); err != nil {
		log.Fatal(err)
	}

	// Rename .md files to .mdx
	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			newPath := strings.TrimSuffix(path, ".md") + ".mdx"
			err := os.Rename(path, newPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add `bash` after the first triple backtick in each code block in .mdx files
	err = filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mdx") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			updatedContent := addBashToCodeBlocks(string(content))
			updatedContent = removeFirstMarkdownLine(updatedContent)
			err = os.WriteFile(path, []byte(updatedContent), info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// addBashToCodeBlocks adds `bash` after the first triple backtick in each code block
func addBashToCodeBlocks(content string) string {
	var result strings.Builder
	inCodeBlock := false
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				inCodeBlock = false
			} else {
				inCodeBlock = true
				if line == "```" {
					line = "```bash"
				}
			}
		}
		result.WriteString(line + "\n")
	}

	return result.String()
}

// removeFirstMarkdownLine removes the first line starting with `##`
func removeFirstMarkdownLine(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	skip := true

	for _, line := range lines {
		if skip && strings.HasPrefix(line, "##") {
			skip = false
			continue
		}
		result.WriteString(line + "\n")
	}

	return result.String()
}
