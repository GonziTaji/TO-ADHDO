package articles

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles/model"
	"github.com/yogusita/to-adhdo/domain/tags"
)

type RenderArticleViewOptions struct {
	ArticleId string `uri:"article_id" binding:"required"`
}

func (c *Controller) GetHandler(ctx *gin.Context) {
	var options RenderArticleViewOptions

	if err := ctx.ShouldBindUri(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	article, err := c.store.GetDetails(options.ArticleId)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if article.Id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "articles/view", article)
}

func (c *Controller) GetCatalogHandler(ctx *gin.Context) {
	var options model.CatalogFilterOptions
	if err := ctx.Bind(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("options: %v\n", options)
	log.Printf("options.tags: %v\n", options.TagsIdsFilter)
	log.Printf("len.tags: %d\n", len(options.TagsIdsFilter))

	catalog_items, err := c.store.Catalog(options)

	tags, err := c.tagsStore.List(tags.ListingTagsOptions{})
	tags_options := []model.TagOption{}
	for _, tag := range tags {
		to := model.TagOption{Id: tag.Id, Name: tag.Name}

		to.Selected = slices.ContainsFunc(options.TagsIdsFilter, func(tagid_filter string) bool {
			return tagid_filter == tag.Id
		})

		tags_options = append(tags_options, to)
	}

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/catalog", model.CatalogData{
		Articles: catalog_items,
		Tags:     tags_options,
		Options:  options,
	})
}
