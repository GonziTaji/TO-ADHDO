package tags

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB, servePage func(string) gin.HandlerFunc) {
	store := CreateStore(db)
	controller := CreateController(store)

	// -----------------------------------------------------------------------
	// Page route
	// -----------------------------------------------------------------------
	router.GET("/tags", servePage("tags/list.html"))

	// Backward-compatibility redirect
	router.GET("/tags/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/tags")
	})

	// -----------------------------------------------------------------------
	// JSON API routes
	// -----------------------------------------------------------------------
	api := router.Group("/api/tags")

	api.GET("", controller.ApiListHandler)
	api.DELETE("/:tagid", controller.DeleteHandler)
}
