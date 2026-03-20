package wishlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WishlistResponseDTO is the combined payload for GET /api/wishlist
type WishlistResponseDTO struct {
	Items              []WishitemDTO        `json:"Items"`
	SearchTerm         string               `json:"SearchTerm"`
	PriceRange         PriceRangeDTO        `json:"PriceRange"`
	PriceSelectedRange PriceRangeDTO        `json:"PriceSelectedRange"`
	TagsSelectOptions  []TagSelectOptionDTO `json:"TagsSelectOptions"`
}

type WishitemDTO struct {
	Id            string           `json:"Id"`
	Name          string           `json:"Name"`
	ObservedPrice int              `json:"ObservedPrice"`
	Tags          []WishitemTagDTO `json:"Tags"`
}

type WishitemTagDTO struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type PriceRangeDTO struct {
	Start int `json:"Start"`
	End   int `json:"End"`
}

type TagSelectOptionDTO struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Count    int    `json:"Count"`
	Selected bool   `json:"Selected"`
}

func (c *Controller) ApiListHandler(ctx *gin.Context) {
	var options WishlistFilterParams
	if err := ctx.BindQuery(&options); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if options.PriceRangeEnd != 0 && options.PriceRangeEnd < options.PriceRangeStart {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "the start of the price range cannot be greater than its end"})
		return
	}

	data, err := c.store.GetWishlist(options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := WishlistResponseDTO{
		SearchTerm: data.SearchTerm,
		PriceRange: PriceRangeDTO{
			Start: data.PriceRange.Start,
			End:   data.PriceRange.End,
		},
		PriceSelectedRange: PriceRangeDTO{
			Start: data.PriceSelectedRange.Start,
			End:   data.PriceSelectedRange.End,
		},
	}

	for _, item := range data.Items {
		dto := WishitemDTO{
			Id:            item.Id,
			Name:          item.Name,
			ObservedPrice: item.ObservedPrice,
		}
		for _, t := range item.Tags {
			dto.Tags = append(dto.Tags, WishitemTagDTO{Id: t.TagId, Name: t.TagName})
		}
		resp.Items = append(resp.Items, dto)
	}

	for _, t := range data.TagsSelectOptions {
		resp.TagsSelectOptions = append(resp.TagsSelectOptions, TagSelectOptionDTO{
			Id:       t.Id,
			Name:     t.Name,
			Count:    t.Count,
			Selected: t.Selected,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}
