package articles

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/internal/domain/tags"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	store := CreateStore(db)
	tagsStore := tags.CreateStore(db)
	service := CreateService(store, &Views{}, tagsStore)
	controller := CreateController(service)

	// Catalog routes

	group := router.Group("catalog")

	group.GET("/", controller.GetCatalogHandler)
	group.GET("/list", controller.GetCatalogListHandler)
	group.GET("/:article_id", controller.GetHandler)

	// Admin routes

	admin := router.Group("admin/articles")

	admin.GET("/", controller.GetListHandler)
	admin.GET("/new", controller.GetFormHandler)
	admin.GET("/:article_id/edit", controller.GetFormHandler)

	admin.POST("/", controller.CreateHandler)
	admin.POST("/uploads", controller.UploadImageHandler)
	admin.PUT("/:article_id", controller.UpdateHandler)
	admin.DELETE("/:article_id", controller.DeleteHandler)
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
