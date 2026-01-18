package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func getTemplatePaths() []string {
	var layouts_templates_paths []string
	var non_layout_templates_paths []string

	_ = filepath.Walk("domain", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".layout.html") {
			layouts_templates_paths = append(layouts_templates_paths, path)
		} else if strings.HasSuffix(path, ".html") {
			non_layout_templates_paths = append(non_layout_templates_paths, path)
		}
		return nil
	})

	return slices.Concat(layouts_templates_paths, non_layout_templates_paths)
}

func main() {
	paths := getTemplatePaths()
	fmt.Println("Templates that would be loaded:")
	for i, path := range paths {
		fmt.Printf("%d: %s\n", i+1, path)
	}
}
