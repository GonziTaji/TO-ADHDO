package tags

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	store := CreateStore(db)
	controller := CreateController(store)

	group := router.Group("tags")

	group.GET("/", controller.GetListHandler)
}
