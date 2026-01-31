package articles

import (
	"github.com/yogusita/to-adhdo/domain/tags"
)

type Image struct {
	Id string

	Path string
	Alt  string

	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type Article struct {
	Id string

	Name        string
	Description string
	Tags        []tags.Tag
	Images      []Image
	Prices      []ArticlePrice

	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type ArticlePrice struct {
	Id          string
	ArticleId   string
	Price       int
	Description string
	CreatedAt   string
}
