package server

import (
	"github.com/gin-gonic/gin"
)

func registerStaticRoutes(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "public/favicon.ico")
	router.Static("/public/media/uploads", "public/media/uploads")
}
