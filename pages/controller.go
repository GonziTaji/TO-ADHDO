package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles"
	"github.com/yogusita/to-adhdo/domain/tags"
)

func ArticleForm(c *gin.Context) {
	task_templates_store := articles.Store{}
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

	tags := []tags.Category{}

	c.HTML(http.StatusOK, "pages/home", gin.H{
		"tags":  tags_list,
		"tasks": tasks,
		"form_values": articles.Article{
			Id:          "",
			Name:        "",
			Description: "",
			Tags:        tags,
			CreatedAt:   "",
			UpdatedAt:   "",
			DeletedAt:   "",
		},
	})
}

func ArticleView(c *gin.Context) {
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

	c.HTML(http.StatusOK, "pages/task_template", gin.H{
		"tags": tags,
		"task": task,
	})
}
