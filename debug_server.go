package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

	return slices.Concat(layouts_templates_paths, non_layout_templates_paths)
}

func main() {
	paths := getTemplatePaths()
	fmt.Printf("Loading templates: %v\n", paths)
	
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully loaded %d templates\n", len(paths))
	
	// Try to execute the form template
	err = tmpl.ExecuteTemplate(os.Stdout, "articles/form", nil)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
}
