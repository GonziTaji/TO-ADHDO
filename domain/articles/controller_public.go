package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	catalog_items, err := c.store.Catalog()

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "articles/catalog", catalog_items)
}
