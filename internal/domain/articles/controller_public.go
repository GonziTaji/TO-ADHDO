package articles

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles/model"
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

	article, err := c.service.GetDetails(options.ArticleId)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if article.Id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "articles-view", article)
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

	catalogData, err := c.service.Catalog(options)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/www/catalog/pages", catalogData)
}

func (c *Controller) GetCatalogListHandler(ctx *gin.Context) {
	var options model.CatalogFilterOptions
	if err := ctx.Bind(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.service.CatalogList(options)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("asdfasdfas\n")

	ctx.HTML(http.StatusOK, "catalog-list", items)
}
