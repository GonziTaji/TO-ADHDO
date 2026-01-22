package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListingTagsOptions struct {
	Limit          int8 `query:"limit"`
	Offset         int8 `query:"offset"`
	IncludeDeleted bool `query:"include_deleted"`
}

type TagItemList struct {
	Tag
	Usage int
}

type TagsListData []TagItemList

type Controller struct {
	store *Store
}

func CreateController(store *Store) *Controller {
	return &Controller{store}
}

func (c *Controller) GetListHandler(ctx *gin.Context) {
	var options ListingTagsOptions
	if err := ctx.ShouldBindQuery(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	list, err := c.store.List(options)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "tags/list", list)
}
