package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"io/fs"
)

func getTemplatePaths() []string {
	var layouts_templates_paths []string
	var non_layout_templates_paths []string

	_ = filepath.Walk("domain", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".layout.html") {
			layouts_templates_paths = append(layouts_templates_paths, path)
		} else if strings.HasSuffix(path, ".html") {
			non_layout_templates_paths = append(non_layout_templates_paths, path)
		}
		return nil
	})

	paths := slices.Concat(layouts_templates_paths, non_layout_templates_paths)
	
	fmt.Printf("Layout templates (%d):\n", len(layouts_templates_paths))
	for _, path := range layouts_templates_paths {
		fmt.Printf("  - %s\n", path)
	}
	
	fmt.Printf("Non-layout templates (%d):\n", len(non_layout_templates_paths))
	for _, path := range non_layout_templates_paths {
		fmt.Printf("  - %s\n", path)
	}
	
	fmt.Printf("Total templates: %d\n", len(paths))
	return paths
}

func main() {
	paths := getTemplatePaths()
	fmt.Printf("\nFinal template order:\n")
	for i, path := range paths {
		fmt.Printf("%d: %s\n", i+1, path)
	}
}
