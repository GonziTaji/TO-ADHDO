package server

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/pages"
	"github.com/yogusita/to-adhdo/domain/task_templates"
	"github.com/yogusita/to-adhdo/env"
)

const DEFAULT_API_PORT = "8080"

// based on https://github.com/gin-gonic/examples/blob/master/secure-web-app/main.go
func setupSecurityHeaders(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	cspPolicy := "default-src 'self'; connect-src *; font-src *; " +
		"script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';"
	c.Header("Content-Security-Policy", cspPolicy)
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	permPolicy := "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=()," +
		"magnetometer=(),gyroscope=(),fullscreen=(self),payment=()"
	c.Header("Permissions-Policy", permPolicy)

	c.Next()
}

func registerStaticRoutes(router *gin.RouterGroup) {
	static_path := "public/"

	router.Use(blockExtensions(".html"))

	router.StaticFile("/favicon.ico", static_path+"favicon.ico")
	router.Static("/public", static_path)
}

func newRouter() *gin.Engine {
	router := gin.Default()

	router.Use(setupSecurityHeaders)

	router.SetFuncMap(template.FuncMap{
		"dict": dictFuncMap,
	})

	// router.LoadHTMLGlob("public/lib/components/*/template.html")
	// router.LoadHTMLGlob("public/lib/pages/*/template.html")

	router.LoadHTMLGlob("public/**/template.html")

	rootRouterGroup := &router.RouterGroup

	registerStaticRoutes(rootRouterGroup)

	pages.RegisterPages(rootRouterGroup)

	apiRouterGroup := router.Group("/api")

	task_templates.RegisterRoutes(apiRouterGroup)

	return router
}

func Start() error {
	router := newRouter()

	port, _ := env.LookupEnvWithDefault("API_PORT", DEFAULT_API_PORT)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("initializing server on address %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
