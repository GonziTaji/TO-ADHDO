package model

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
