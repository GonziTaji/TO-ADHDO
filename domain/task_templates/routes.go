package task_templates

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	c := Controller{}

	router.GET("task_templates", c.GetListHandler)
	router.POST("task_templates", c.CreateHandler)
	router.GET("task_templates/:task_id/list-item", c.GetTaskAsListItem)
	router.DELETE("task_templates/:task_id", c.DeleteHandler)
}
