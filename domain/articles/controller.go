package articles

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
	"github.com/yogusita/to-adhdo/domain/uploads"
)

const articles_images_bucket = "articles_images"

type TagOption struct {
	Id       string
	Name     string
	Disabled bool
}

type ArticleFormTemplateData struct {
	Article    Article
	TagOptions []TagOption
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

type ArticleFormData struct {
	Id                    string   `form:"id"`
	Name                  string   `form:"name" binding:"required"`
	Description           string   `form:"description"`
	TagNames              []string `form:"tags_names"`
	TagIds                []string `form:"tags_ids"`
	NewPrice              int      `form:"new_price"`
	NewPriceDescription   string   `form:"new_price_description"`
	ArticleImageFilenames []string `form:"article_images_filenames"`
	ArticleImageIds       []string `form:"article_images_ids"`
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

	tags, err := c.tagsStore.List(tags.ListingTagsOptions{})
	tag_options := make([]TagOption, len(tags))

	tags_ids_in_article := make(map[string]bool)

	for _, tag := range article.Tags {
		tags_ids_in_article[tag.Id] = true
	}

	for i, tag := range tags {
		tag_options[i] = TagOption{
			Name:     tag.Name,
			Id:       tag.Id,
			Disabled: tags_ids_in_article[tag.Id],
		}
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	formData := ArticleFormTemplateData{
		Article:    article,
		TagOptions: tag_options,
	}

	ctx.HTML(http.StatusOK, "articles/form", formData)
}

func (c *Controller) GetCatalogHandler(ctx *gin.Context) {
	var options ListingArticlesOptions

	if err := ctx.ShouldBindQuery(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	articles, err := c.store.List(&options)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/catalog", articles)
}

func (c *Controller) GetListHandler(ctx *gin.Context) {
	var options ListingArticlesOptions

	if err := ctx.ShouldBindQuery(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	articles, err := c.store.List(&options)

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
	var form ArticleFormData

	if err := ctx.Bind(&form); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	article := Article{
		Name:        form.Name,
		Description: form.Description,
		Tags:        []tags.Tag{},
		Images:      []ArticleImage{},
	}

	// should this be in a "service" layer?
	// Same logic as in the UpdateHandler
	for i, name := range form.TagNames {
		tag := tags.Tag{
			Id:   form.TagIds[i],
			Name: name,
		}

		article.Tags = append(article.Tags, tag)
	}

	for i, id := range form.ArticleImageIds {
		image := ArticleImage{
			Id:       id,
			Filename: form.ArticleImageFilenames[i],
		}

		article.Images = append(article.Images, image)
	}

	if form.NewPrice != 0 {
		article.Prices = append(article.Prices, ArticlePrice{
			Price:       form.NewPrice,
			Description: form.NewPriceDescription,
		})
	}

	article_id, err := c.store.Create(&article)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("Location", "/articles/"+article_id)
	ctx.Status(http.StatusCreated)
}

func (c *Controller) UpdateHandler(ctx *gin.Context) {
	var form ArticleFormData

	if err := ctx.Bind(&form); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	article := Article{
		Id:          form.Id,
		Name:        form.Name,
		Description: form.Description,
		Tags:        []tags.Tag{},
	}

	// should this be in a "service" layer?
	// Same logic as in the CreateHandler
	for i, name := range form.TagNames {
		tag := tags.Tag{
			Id:   form.TagIds[i],
			Name: name,
		}

		article.Tags = append(article.Tags, tag)
	}

	for i, id := range form.ArticleImageIds {
		image := ArticleImage{
			Id:       id,
			Filename: form.ArticleImageFilenames[i],
		}

		article.Images = append(article.Images, image)
	}

	if form.NewPrice != 0 {
		article.Prices = append(article.Prices, ArticlePrice{
			Price:       form.NewPrice,
			Description: form.NewPriceDescription,
		})
	}

	err := c.store.Update(&article)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("Location", "/articles/"+article.Id)
	ctx.Status(http.StatusOK)
}

func (c *Controller) UploadImageHandler(ctx *gin.Context) {
	fileheader, err := ctx.FormFile("file")

	file, err := fileheader.Open()

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	file_ext := filepath.Ext(fileheader.Filename)

	saved_filename, err := uploads.SaveFile(articles_images_bucket, file_ext, file)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := file.Close(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/form/image-miniature", ArticleImage{
		Id:       "",
		Filename: saved_filename,
		Url:      uploads.GetFilePublicUrl(articles_images_bucket, saved_filename),
	})
}
