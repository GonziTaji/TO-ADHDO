package articles

import (
	"github.com/yogusita/to-adhdo/domain/tags"
)

type Article struct {
	Id string

	Name        string
	Description string
	Tags        []tags.Tag
	Images      []ArticleImage
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

type ArticleImage struct {
	Id        string
	ArticleId string
	Path      string
	CreatedAt string
}
