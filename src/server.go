package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/env"
)

type Task struct {
	Id          string
	Name        string
	Description string
	Tags        []Tag
	Deleted     bool
}

type Tag struct {
	Id   string
	Name string
}

var TAGS = []Tag{
	{Id: "1", Name: "aseo"},
	{Id: "2", Name: "orden"},
	{Id: "3", Name: "living"},
	{Id: "4", Name: "baño"},
}

var TASKS = []Task{
	{Id: "1", Name: "test task", Description: "hardcoded task"},
}

const DEFAULT_API_PORT = "8080"

// based on https://github.com/gin-gonic/examples/blob/master/secure-web-app/main.go
func SetupSecurityHeaders(c *gin.Context) {
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

func GetTaskListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "tasks_list.html", TASKS)
}

func DeleteTaskHandler(c *gin.Context) {
	task_id := c.Param("task_id")

	var found = false

	for i := range TASKS {
		if TASKS[i].Id == task_id {
			TASKS[i].Deleted = true
			found = true
			break
		}
	}

	if found == false {
		c.JSON(http.StatusNotFound, "Task not found")
	} else {
		c.JSON(http.StatusOK, "Task not found")
	}
}

func PostTaskHandler(c *gin.Context) {
	name, ok := c.GetPostForm("name")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "field 'name' is required",
		})

		return
	}

	description, _ := c.GetPostForm("description")

	task_tags_names, _ := c.GetPostFormArray("tags")

	new_task := Task{
		Name:        name,
		Description: description,
		Tags:        []Tag{},
		Deleted:     false,
	}

	new_tags_names := []string{}

	for _, task_tag := range task_tags_names {
		var tag_exists = false

		for _, existing_tag := range TAGS {
			if task_tag == existing_tag.Name {
				tag_exists = true

				new_task.Tags = append(new_task.Tags, existing_tag)

				break
			}
		}

		if !tag_exists {
			new_tag := Tag{
				Id:   strconv.Itoa(len(TAGS)),
				Name: task_tag,
			}

			new_task.Tags = append(new_task.Tags, new_tag)

			TAGS = append(TAGS, new_tag)
			new_tags_names = append(new_tags_names, new_tag.Name)
		}
	}

	TASKS = append(TASKS, new_task)

	c.JSON(http.StatusCreated, gin.H{
		"new_tags": new_tags_names,
	})
}

func RegisterPageRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"tags":  TAGS,
			"tasks": TASKS,
		})
	})
}

func RegisterApiRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/api")

	{
		api_handler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Uniform API",
			})
		}

		api.GET("", api_handler)
		api.POST("tasks", PostTaskHandler)
		api.GET("tasks", GetTaskListHandler)
		api.DELETE("tasks/:task_id", DeleteTaskHandler)
	}
}

func RegisterStaticRoutes(router *gin.RouterGroup) {
	static_path := "src/public/"

	router.StaticFile("/favicon.ico", static_path+"favicon.ico")
	router.Static("/public", static_path)
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.Use(SetupSecurityHeaders)

	router.LoadHTMLGlob("src/templates/*")

	RegisterStaticRoutes(&router.RouterGroup)
	RegisterApiRoutes(&router.RouterGroup)
	RegisterPageRoutes(router)

	return router
}

func main() {
	router := NewRouter()

	port, _ := env.LookupEnvWithDefault("API_PORT", DEFAULT_API_PORT)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("initializing server on address %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
