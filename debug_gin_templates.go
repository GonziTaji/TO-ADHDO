package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"io/fs"
	
	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles"
	"github.com/yogusita/to-adhdo/domain/tags"
	_ "modernc.org/sqlite"
)

// Copy of the server's template loading function
func getTemplatePaths() []string {
	var layouts_templates_paths []string
	var non_layout_templates_paths []string

	// Simplified walk function
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
	// Setup database
	db, err := sql.Open("sqlite", "main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup Gin exactly like the server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Load templates exactly like the server
	templatePaths := getTemplatePaths()
	fmt.Printf("Loading templates: %v\n", templatePaths)
	router.LoadHTMLFiles(templatePaths...)

	// Setup stores and controller exactly like the server
	articlesStore := articles.CreateStore(db)
	tagsStore := tags.CreateStore(db)
	views := &articles.Views{}
	controller := articles.CreateController(articlesStore, views, tagsStore)

	router.GET("/test", controller.GetFormHandler)

	// Create test request
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Content-Length: %d\n", len(w.Body.String()))
	fmt.Printf("Response body:\n%s\n", w.Body.String())
}
