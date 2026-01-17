package articles

import "github.com/yogusita/to-adhdo/domain/tags"

type BaseRecord struct {
	Id        string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type Image struct {
	BaseRecord
	Path string
	Alt  string
}

type Article struct {
	BaseRecord
	Name        string
	Description string
	Tags        []tags.Tag
	Images      []Image
}
