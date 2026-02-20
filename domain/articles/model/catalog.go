package model

type CatalogItem struct {
	Id           string
	Name         string
	Price        int
	Tags         []struct{ Name string }
	ThumbnailUrl string
}
