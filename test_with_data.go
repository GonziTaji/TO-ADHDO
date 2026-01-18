package main

import (
	"fmt"
	"html/template"
	"os"
)

// Mock structures
type Tag struct {
	Id   string
	Name string
}

type Article struct {
	Id          string
	Name        string
	Description string
	Tags        []Tag
	TagOptions  []Tag  // Add the missing field
}

func main() {
	// Create test data with TagOptions
	article := Article{
		Id:          "",
		Name:        "",
		Description: "",
		Tags:        []Tag{},
		TagOptions: []Tag{
			{Id: "1", Name: "Technology"},
			{Id: "2", Name: "Programming"},
			{Id: "3", Name: "Go"},
		},
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		"domain/articles/templates/form.html",
		"domain/shared/templates/meta_tags.component.html",
	)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}
	
	err = tmpl.ExecuteTemplate(os.Stdout, "articles/form", article)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
}
