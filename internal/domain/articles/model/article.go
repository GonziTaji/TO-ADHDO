package model

import (
	"github.com/yogusita/to-adhdo/domain/tags"
)

type Article struct {
	Id string

	Name              string
	Description       string
	Tags              []tags.Tag
	Images            []ArticleImage
	Prices            []ArticlePrice
	AvailableForTrade bool

	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type TagOption struct {
	Id       string
	Name     string
	Disabled bool
	Selected bool
}

type ArticleFormTemplateData struct {
	Article    Article
	TagOptions []TagOption
}
