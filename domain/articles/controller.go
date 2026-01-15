package articles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct{
	store Store
	views Views
}

type CreateArticleResponse struct {
	Id        string
	CreatedAt string
}

func NewArticlesController(store Store, views Views) *Controller {
	return &ArticlesController{
		store, views
	}
}

func (Controller) GetHandler(c *gin.Context) {
	article_id := c.Params.ByName("article_id")

	article, err := controller.Get(article_id)

	if err != nil {
		c.String(http.StatusInternalServerError, "ERROR: "+err.Error())
		return
	}

	if article.Id == "" {
		c.Status(http.StatusNotFound)
		return
	}

	view_name := c.Params.ByName("view_name")

	fmt.Printf("VIEW NAME: %s\n", view_name)

	switch view_name {
	case "list-item":
		err = views.AsListItem(c.Writer, article)
	default:
		c.String(http.StatusBadRequest, "invalid view name: \"%s\"", view_name)
		return
	}

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}

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

	articles, err := store.List(int8(limit), false)

	if err != nil {
		// TODO: Template para error cargando componente con reintento
		c.String(http.StatusBadRequest, "ERROR: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "articles_list.html", articles)
}

func (Controller) DeleteHandler(c *gin.Context) {
	article_id, exists := c.Params.Get("article_id")

	if !exists || article_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "malformed url",
		})

		return
	}

	err := store.Delete(article_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (Controller) CreateHandler(c *gin.Context) {
	categories := c.PostFormArray("categories")

	new_article := Article{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Categories:  []Category{},
	}

	if new_article.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "field 'name' is required",
		})

		return
	}

	for _, category_data := range categories {

	}

	err := store.Create(name, description, article_tags_names)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": article_id})
}
