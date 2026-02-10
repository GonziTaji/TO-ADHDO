package articles

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	store := CreateStore(db)
	tagsStore := tags.CreateStore(db)
	controller := CreateController(store, &Views{}, tagsStore)

	// Catalog routes

	group := router.Group("articles")

	group.GET("/", controller.GetCatalogHandler)
	group.GET("/:article_id", controller.GetHandler)

	// Admin routes

	admin := router.Group("admin/" + group.BasePath())

	admin.GET("/", controller.GetListHandler)
	admin.GET("/new", controller.GetFormHandler)
	admin.GET("/:article_id/edit", controller.GetFormHandler)

	admin.POST("/", controller.CreateHandler)
	admin.POST("/uploads", controller.UploadImageHandler)
	admin.PUT("/:article_id", controller.UpdateHandler)
	admin.DELETE("/:article_id", controller.DeleteHandler)
}
