package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TagDTO is the JSON shape for a single tag in API responses.
type TagDTO struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

func (c *Controller) ApiListHandler(ctx *gin.Context) {
	var options ListingTagsOptions
	if err := ctx.Bind(&options); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := c.store.List(options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]TagDTO, 0, len(list))
	for _, t := range list {
		dtos = append(dtos, TagDTO{Id: t.Id, Name: t.Name})
	}

	ctx.JSON(http.StatusOK, dtos)
}
