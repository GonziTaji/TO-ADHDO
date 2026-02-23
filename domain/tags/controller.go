package tags

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsageFilter string

const (
	Unused UsageFilter = "unused"
	Used   UsageFilter = "used"
)

type ListingTagsOptions struct {
	Limit      int8        `form:"limit"`
	Offset     int8        `form:"offset"`
	Usage      UsageFilter `form:"usage"`
	SearchTerm string      `form:"s"`
}

type TagItemList struct {
	Tag
	Usage int
}

type TagsListData struct {
	Options ListingTagsOptions
	Tags    []TagItemList
}

type Controller struct {
	store *Store
}

func CreateController(store *Store) *Controller {
	return &Controller{store}
}

func (c *Controller) GetListHandler(ctx *gin.Context) {
	var options ListingTagsOptions
	if err := ctx.Bind(&options); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("options: %v", options)

	list, err := c.store.List(options)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "tags/list", TagsListData{
		Tags:    list,
		Options: options,
	})
}

func (c *Controller) DeleteHandler(ctx *gin.Context) {
	tagid := ctx.Param("tagid")

	err := c.store.Delete(tagid)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
