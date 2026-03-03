package wishlist

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	store := CreateStore(db)
	controller := CreateController(&store)

	public_group := router.Group("wishlist")

	public_group.GET("/", controller.GetListHandler)

	api_group := router.Group("api/" + public_group.BasePath())

	api_group.GET("/preview", controller.GetPreview)

	admin_group := router.Group("admin/" + public_group.BasePath())

	admin_group.GET("/", controller.GetAdminListHandler)

	admin_group.GET("/new", controller.GetFormHandler)
	admin_group.POST("/new", controller.CreateHandler)

	admin_group.GET("/:id", controller.GetFormHandler)
	admin_group.PUT("/:id", controller.UpdateHandler)
	admin_group.DELETE("/:id", controller.DeleteHandler)
}
