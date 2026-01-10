package pages

import "github.com/gin-gonic/gin"

func RegisterPages(router *gin.RouterGroup) {
	router.GET("/", indexHandler)
	router.GET("task_templates/:task_id", taskTemplateHandler)
}
