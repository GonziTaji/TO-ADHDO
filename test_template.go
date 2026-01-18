package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	tmpl, err := template.ParseFiles("domain/articles/templates/form.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}
	
	err = tmpl.ExecuteTemplate(os.Stdout, "articles/form", nil)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
}
