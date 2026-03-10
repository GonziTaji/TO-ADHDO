package wishlist

import "github.com/yogusita/to-adhdo/domain/tags"

type WishitemMeta struct {
	URL         string
	Title       string
	Name        string
	Description string
	Image       string
	Price       int
}

type WishlistSortValue string
type SortDirection string

const (
	WishlistSortByPrice    WishlistSortValue = "price"
	WishlistSortByCratedAt WishlistSortValue = "created"
)

const (
	SortDirectionDesc SortDirection = "desc"
	SortDirectionAsc  SortDirection = "asc"
)

type WishitemFormData struct {
	Id            string   `form:"id"`
	Name          string   `form:"name" binding:"required"`
	Description   string   `form:"description" binding:"required"`
	ExternalUrl   string   `form:"external_url" binding:"required"`
	TagsNames     []string `form:"tags_names"`
	TagsIds       []string `form:"tags_ids"`
	ObservedPrice string   `form:"observed_price" binding:"required"`
	ImageUrl      string   `form:"image_url"`
}

type Wishitem struct {
	Id            string
	Name          string
	Description   string
	ExternalUrl   string
	CratedAt      string
	Tags          []tags.Tag
	ObservedPrice int
	ImgaeUrl      string
}

type WishlistFilterParams struct {
	SearchTerm      string            `form:"search"`
	TagsIds         []string          `form:"tags"`
	PriceRangeStart int               `form:"price_start"`
	PriceRangeEnd   int               `form:"price_end"`
	SortBy          WishlistSortValue `form:"sort,default=desc"`
	SortDirection   SortDirection     `form:"dir,default=created"`
}

type TagWithCount struct {
	tags.Tag
	Count int
}

type TagSelectOption struct {
	TagWithCount
	Selected bool
}

type PriceRangeData struct {
	Start int
	End   int
}

type WishlistData struct {
	Items              []Wishitem
	SearchTerm         string
	TagsSelectOptions  []TagSelectOption
	PriceRange         PriceRangeData
	PriceSelectedRange PriceRangeData

	Sort struct {
		Column    WishlistSortValue
		Direction SortDirection
	}
}

type WishitemFormTemplateData struct {
	Record Wishitem
	Tags   []tags.Tag
}
