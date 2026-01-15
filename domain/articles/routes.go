package articles

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	c := Controller{}

	group := router.Group("articles")

	group.GET("/", c.GetListHandler)
	group.GET("/:article_id/:view_id", c.GetHandler)
	group.POST("/", c.CreateHandler)
	group.DELETE("/:article_id", c.DeleteHandler)
}
