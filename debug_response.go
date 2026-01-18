package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	
	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles"
	"github.com/yogusita/to-adhdo/domain/tags"
	_ "modernc.org/sqlite"
)

func main() {
	// Setup database
	db, err := sql.Open("sqlite", "main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Setup template
	router.LoadHTMLFiles(
		"domain/articles/templates/form.html",
		"domain/shared/templates/meta_tags.component.html",
	)

	// Setup stores and controller
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
