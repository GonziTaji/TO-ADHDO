package model

type CatalogItem struct {
	Id           string
	Name         string
	Price        int
	Tags         []struct{ Name string }
	ThumbnailUrl string
}

type ArticleDetailTag struct {
	Id   string
	Name string
}

type ArticleDetails struct {
	Id          string
	Name        string
	Description string
	Tags        []ArticleDetailTag
	Price       int
	ImagesUrls  []string
	IsDeleted   bool
}

type CatalogFilterOptions struct {
	SearchTerm    string   `form:"s"`
	TagsIdsFilter []string `form:"tags"`
}

type CatalogData struct {
	Articles []CatalogItem
	Tags     []TagOption
	Options  CatalogFilterOptions
}
