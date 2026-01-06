package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/database"
	"github.com/yogusita/to-adhdo/env"
)

const DEFAULT_API_PORT = "8080"

// based on https://github.com/gin-gonic/examples/blob/master/secure-web-app/main.go
func setupSecurityHeaders(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	cspPolicy := "default-src 'self'; connect-src *; font-src *; " +
		"script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';"
	c.Header("Content-Security-Policy", cspPolicy)
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	permPolicy := "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=()," +
		"magnetometer=(),gyroscope=(),fullscreen=(self),payment=()"
	c.Header("Permissions-Policy", permPolicy)

	c.Next()
}

func getTaskTemplatesHandler(c *gin.Context) {
	tasks, err := database.GetAvailableTaskTemplates(10)

	if err != nil {
		// TODO: Template para error cargando componente con reintento
		c.String(http.StatusBadRequest, "ERROR: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "task_templates_list.html", tasks)
}

func deleteTaskHandler(c *gin.Context) {
	task_id, exists := c.Params.Get("task_id")

	if !exists || task_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "malformed url",
		})

		return
	}

	err := database.DeleteTaskTemplate(task_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func postTaskTemplatesHandler(c *gin.Context) {
	name, ok := c.GetPostForm("name")

	if !ok || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "field 'name' is required",
		})

		return
	}

	description, _ := c.GetPostForm("description")

	task_tags_names, _ := c.GetPostFormArray("tags")

	task_id, err := database.CreateTaskTemplate(name, description, task_tags_names)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": task_id})
}

func registerPageRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		tasks, err := database.GetAvailableTaskTemplates(20)

		fmt.Printf("Tasks found: %v\n", tasks)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		tags, err := database.GetAvailableTags(100)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tags":  tags,
			"tasks": tasks,
			"form_values": database.TaskTemplate{
				Id:          "",
				Name:        "",
				Description: "",
				Tags:        []database.Tag{},
				CreatedAt:   0,
				UpdatedAt:   0,
				DeletedAt:   0,
			},
		})
	})

	router.GET("task_templates/:task_id", func(c *gin.Context) {
		task_id := c.Param("task_id")

		task, err := database.GetTaskTemplate(task_id)

		fmt.Printf("Task found: \"%s\n\"", task.Name)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if task.Id == "" {
			c.Status(http.StatusNotFound)
			return
		}

		tags, err := database.GetAvailableTags(100)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.HTML(http.StatusOK, "task_template.html", gin.H{
			"tags": tags,
			"task": task,
		})
	})
}

func registerApiRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/api")

	{
		api_handler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello world",
			})
		}

		api.GET("echo", api_handler)
		api.POST("tasks_templates", postTaskTemplatesHandler)
		api.GET("tasks_templates", getTaskTemplatesHandler)
		api.DELETE("tasks_templates/:task_id", deleteTaskHandler)
	}
}

func registerStaticRoutes(router *gin.RouterGroup) {
	static_path := "public/"

	router.StaticFile("/favicon.ico", static_path+"favicon.ico")
	router.Static("/public", static_path)
}

func getTemplateFunctionsMap() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...any) map[string]any {
			if len(values)%2 != 0 {
				panic("dict requires even number of arguments")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					panic("dict keys must be strings")
				}
				m[key] = values[i+1]
			}
			return m
		},
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()

	router.Use(setupSecurityHeaders)
	router.SetFuncMap(getTemplateFunctionsMap())
	router.LoadHTMLGlob("templates/*")

	registerStaticRoutes(&router.RouterGroup)
	registerApiRoutes(&router.RouterGroup)
	registerPageRoutes(router)

	return router
}

func Start() error {
	router := newRouter()

	port, _ := env.LookupEnvWithDefault("API_PORT", DEFAULT_API_PORT)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("initializing server on address %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
