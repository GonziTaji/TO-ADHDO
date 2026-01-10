package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
	"github.com/yogusita/to-adhdo/domain/task_templates"
)

func indexHandler(c *gin.Context) {
	task_templates_store := task_templates.Store{}
	task_tags_store := tags.Store{}

	tasks, err := task_templates_store.List(10, false)

	fmt.Printf("Tasks found: %v\n", tasks)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tags_list, err := task_tags_store.List(100, false)

	fmt.Printf("Tags found: %v\n", tags_list)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	task_tags := []tags.Tag{}

	c.HTML(http.StatusOK, "template.html", gin.H{
		"tags":  tags_list,
		"tasks": tasks,
		"form_values": task_templates.TaskTemplate{
			Id:          "",
			Name:        "",
			Description: "",
			Tags:        task_tags,
			CreatedAt:   "",
			UpdatedAt:   "",
			DeletedAt:   "",
		},
	})
}

func taskTemplateHandler(c *gin.Context) {
	task_id := c.Param("task_id")

	task_templates_store := task_templates.Store{}
	task_tags_store := tags.Store{}

	task, err := task_templates_store.Get(task_id)

	fmt.Printf("Task found: \"%s\"\n", task.Name)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if task.Id == "" {
		c.Status(http.StatusNotFound)
		return
	}

	tags, err := task_tags_store.List(100, false)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "task_templates/view.html", gin.H{
		"tags": tags,
		"task": task,
	})
}
