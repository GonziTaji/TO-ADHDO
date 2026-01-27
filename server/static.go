package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerStaticRoutes(router *gin.Engine) {
	static_path := "public/"

	router.Use(blockExtensions(".html"))

	router.StaticFile("/favicon.ico", static_path+"favicon.ico")
	router.Static("/public", static_path)
}

func blockExtensions(exts ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		for _, ext := range exts {
			if strings.HasSuffix(path, ext) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		c.Next()
	}
}
