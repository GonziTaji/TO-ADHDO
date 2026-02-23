package server

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles"
	"github.com/yogusita/to-adhdo/domain/tags"
	"github.com/yogusita/to-adhdo/env"
	"github.com/yogusita/to-adhdo/server/funcmap"
	_ "modernc.org/sqlite"
)

// Based on https://github.com/gin-gonic/examples/blob/master/secure-web-app/main.go
func setupSecurityHeaders(c *gin.Context) {
	cspPolicy := "default-src 'self'; connect-src *; font-src *; " +
		"script-src-elem * 'unsafe-inline'; img-src * data: blob:; style-src * 'unsafe-inline';"

	permPolicy := "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=()," +
		"magnetometer=(),gyroscope=(),fullscreen=(self),payment=()"

	header_value_pairs := [][2]string{
		{"Referrer-Policy", "strict-origin"},
		{"Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"},
		{"Content-Security-Policy", cspPolicy},
		{"Permissions-Policy", permPolicy},
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
		{"X-Content-Type-Options", "nosniff"},
	}

	for _, pair := range header_value_pairs {
		var (
			key   = pair[0]
			value = pair[1]
		)

		c.Header(key, value)
	}

	c.Next()
}

func newRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.Use(setupSecurityHeaders)

	router.SetFuncMap(funcmap.GetFuncMap())

	registerStaticRoutes(router)

	router.LoadHTMLFiles(getTemplatePaths()...)

	articles.RegisterRoutes(router, db)
	tags.RegisterRoutes(router, db)

	router.Use(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "Nothing here uwu")
	})

	return router
}

// TODO: most of this is not longer required. review and rewrite this func
func getTemplatePaths() []string {
	// Layouts need to be loaded first to define their blocks before the rest of the templates define them
	var layouts_templates_paths []string
	var non_layout_templates_paths []string

	// TODO: take this out of here. Loading files should be done outside the initialization of the router)
	// Walk function never errors
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

func Start() error {
	db, err := newDB()

	if err != nil {
		return err
	}

	router := newRouter(db)

	port, _ := env.LookupEnvWithDefault("API_PORT", "8080")

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second, // set timeout due CWE-400 - Potential Slowloris Attack
	}

	log.Printf("Initializing server on port %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
