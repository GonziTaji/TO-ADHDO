package wishlist

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB, servePage func(string) gin.HandlerFunc) {
	store := CreateStore(db)
	controller := CreateController(&store)

	// -----------------------------------------------------------------------
	// Page route
	// -----------------------------------------------------------------------
	router.GET("/wishlist", servePage("wishlist/wishlist.html"))

	// Backward-compatibility redirect
	router.GET("/wishlist/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/wishlist")
	})

	// -----------------------------------------------------------------------
	// JSON API routes
	// -----------------------------------------------------------------------
	api := router.Group("/api/wishlist")

	api.GET("", controller.ApiListHandler)
	api.GET("/preview", controller.GetPreview)

	// -----------------------------------------------------------------------
	// Admin routes (kept for compatibility, HTML-rendered)
	// -----------------------------------------------------------------------
	admin := router.Group("admin/wishlist")

	admin.GET("/", controller.GetAdminListHandler)
	admin.GET("/new", controller.GetFormHandler)
	admin.POST("/new", controller.CreateHandler)
	admin.GET("/:id", controller.GetFormHandler)
	admin.PUT("/:id", controller.UpdateHandler)
	admin.DELETE("/:id", controller.DeleteHandler)
}
