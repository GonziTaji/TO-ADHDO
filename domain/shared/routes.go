package shared

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	group := router.Group("shared")

	group.Static("/static", "domain/shared/static")
}
