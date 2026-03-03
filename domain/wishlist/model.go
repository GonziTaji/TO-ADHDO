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

type WishitemFormData struct {
	Id            string `form:"id"`
	Name          string `form:"name" binding:"required"`
	Description   string `form:"description" binding:"required"`
	ExternalUrl   string `form:"external_url" binding:"required"`
	TagsNames     string `form:"tags_names"`
	TagsIds       string `form:"tags_ids"`
	ObservedPrice string `form:"observed_price" binding:"required"`
	ImageUrl      string `form:"image_url"`
}

type Wishitem struct {
	Id            string
	Name          string
	Description   string
	ExternalUrl   string
	CratedAt      string
	Tags          []tags.Tag
	ObservedPrice string
	ImgaeUrl      string
}

type WishitemFormTemplateData struct {
	Record Wishitem
	Tags   []tags.Tag
}
