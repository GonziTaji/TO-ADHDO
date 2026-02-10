package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RenderArticleViewOptions struct {
	ArticleId string `uri:"article_id" binding:"required"`
	ViewName  string `uri:"view_name" binding:"required"`
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
