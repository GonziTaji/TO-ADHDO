package articles

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
)

type ArticleFormData struct {
	Article    Article
	TagOptions []tags.Tag
}

type RenderArticleViewOptions struct {
	ArticleId string `uri:"article_id" binding:"required"`
	ViewName  string `uri:"view_name" binding:"required"`
}

type ListingArticlesOptions struct {
	Limit          int8 `query:"limit"`
	Offset         int8 `query:"offset"`
	IncludeDeleted bool `query:"include_deleted"`
}

type CreateArticleData struct {
	Name        string   `form:"name" binding:"required"`
	Description string   `form:"description"`
	TagNames    []string `form:"tags_names"`
	TagIds      []string `form:"tags_ids"`
}

type UpdateArticleData struct {
	Id string `form:"id"`
	CreateArticleData
}

type CreateArticleResponse struct {
	Id        string
	CreatedAt string
}

type Controller struct {
	store     *Store
	views     *Views
	tagsStore *tags.Store
}

func CreateController(store *Store, views *Views, tagsStore *tags.Store) *Controller {
	return &Controller{store, views, tagsStore}
}

func (c *Controller) GetHandler(ctx *gin.Context) {
	var options RenderArticleViewOptions

	if err := ctx.ShouldBindUri(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	article, err := c.store.Get(options.ArticleId)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if article.Id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	switch options.ViewName {
	case "list-item":
		err = c.views.AsListItem(ctx.Writer, article)
	default:
		ctx.String(http.StatusBadRequest, "invalid view name: \"%s\"", options.ViewName)
		return
	}

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func (c *Controller) GetFormHandler(ctx *gin.Context) {
	article_id := ctx.Param("article_id")

	var article Article

	if article_id != "" {
		var err error
		article, err = c.store.Get(article_id)

		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	tagOptions, err := c.tagsStore.List(tags.ListingTagsOptions{})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("options: %v\n", tagOptions)

	formData := ArticleFormData{
		Article:    article,
		TagOptions: tagOptions,
	}

	ctx.HTML(http.StatusOK, "articles/form", formData)
}

func (c *Controller) GetListHandler(ctx *gin.Context) {
	var options ListingArticlesOptions

	if err := ctx.ShouldBindQuery(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	articles, err := c.store.List(&options)

	log.Printf("articles in list: %v", articles)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/list", articles)
}

func (c *Controller) DeleteHandler(ctx *gin.Context) {
	article_id := ctx.Param("article_id")

	if article_id == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := c.store.Delete(article_id); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *Controller) CreateHandler(ctx *gin.Context) {
	var form CreateArticleData

	if err := ctx.Bind(&form); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	new_article := Article{
		Name:        form.Name,
		Description: form.Description,
		Tags:        []tags.Tag{},
	}

	// should this be in a "service" layer?
	for i, name := range form.TagNames {
		tag := tags.Tag{
			Id:   form.TagIds[i],
			Name: name,
		}

		new_article.Tags = append(new_article.Tags, tag)
	}

	article_id, err := c.store.Create(&new_article)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": article_id})
}

func (c *Controller) UpdateHandler(ctx *gin.Context) {
	var form UpdateArticleData

	if err := ctx.Bind(&form); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	new_article := Article{
		Id:          form.Id,
		Name:        form.Name,
		Description: form.Description,
		Tags:        []tags.Tag{},
	}

	// should this be in a "service" layer?
	for i, name := range form.TagNames {
		tag := tags.Tag{
			Id:   form.TagIds[i],
			Name: name,
		}

		new_article.Tags = append(new_article.Tags, tag)
	}

	article_id, err := c.store.Create(&new_article)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": article_id})
}
