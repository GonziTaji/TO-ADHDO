package server

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const webDistDir = "web/dist"

// ServePage returns a gin handler that serves a specific built HTML file from web/dist.
// path is relative to web/dist, e.g. "articles/catalog.html".
func ServePage(path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullPath := filepath.Join(webDistDir, path)
		log.Printf("%s\n", fullPath)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.File(fullPath)
	}
}
