package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
