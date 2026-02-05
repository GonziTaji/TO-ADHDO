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

	group := router.Group("articles")

	group.GET("/", controller.GetListHandler)
	group.GET("/catalog", controller.GetCatalogHandler)
	group.GET("/new", controller.GetFormHandler)
	group.GET("/:article_id", controller.GetHandler)
	group.GET("/:article_id/edit", controller.GetFormHandler)

	group.POST("/", controller.CreateHandler)
	group.POST("/uploads", controller.UploadImageHandler)
	group.PUT("/:article_id", controller.UpdateHandler)
	group.DELETE("/:article_id", controller.DeleteHandler)
}
