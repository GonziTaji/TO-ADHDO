package tags

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	store := CreateStore(db)
	controller := CreateController(store)

	group := router.Group("tags")
	admin := router.Group("admin/" + group.BasePath())

	admin.GET("/", controller.GetListHandler)
	admin.DELETE("/:tagid", controller.DeleteHandler)
}
