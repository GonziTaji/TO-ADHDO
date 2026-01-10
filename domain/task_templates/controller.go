package task_templates

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

var store Store = Store{}

func (Controller) GetListHandler(c *gin.Context) {
	limit := 10

	if limit_query := c.Query("limit"); limit_query != "" {
		l, err := strconv.Atoi(limit_query)

		if err == nil {
			limit = l
		} else {
			fmt.Printf("error parsing queryparam `limit`: %s\n", err.Error())
		}
	}

	tasks, err := store.List(int8(limit), false)

	if err != nil {
		// TODO: Template para error cargando componente con reintento
		c.String(http.StatusBadRequest, "ERROR: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "task_templates_list.html", tasks)
}

func (Controller) DeleteHandler(c *gin.Context) {
	task_id, exists := c.Params.Get("task_id")

	if !exists || task_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "malformed url",
		})

		return
	}

	err := store.Delete(task_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (Controller) CreateHandler(c *gin.Context) {
	name, ok := c.GetPostForm("name")

	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "field 'name' is required",
		})

		return
	}

	description, _ := c.GetPostForm("description")

	task_tags_names, _ := c.GetPostFormArray("tags")

	task_id, err := store.Create(name, description, task_tags_names)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": task_id})
}

func (Controller) GetTaskAsListItem(c *gin.Context) {
	task_id := c.Param("task_id")

	task, err := store.Get(task_id)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if task.Id == "" {
		c.Status(http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/task_template_list.html")

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tmpl.ExecuteTemplate(c.Writer, "task_template_list_item", task)
}
