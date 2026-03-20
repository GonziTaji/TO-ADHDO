package server

import (
	"github.com/gin-gonic/gin"
)

func registerStaticRoutes(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "public/favicon.ico")
	router.Static("/public/media/uploads", "public/media/uploads")

	// Serve Vite hashed JS/CSS assets for the built web app
	router.Static("/assets", webDistDir+"/assets")

	// Home page
	router.GET("/", ServePage("index.html"))
}
