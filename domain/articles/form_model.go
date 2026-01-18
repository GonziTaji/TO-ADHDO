package articles

import "github.com/yogusita/to-adhdo/domain/tags"

// FormData represents the data structure needed for the article form template
type FormData struct {
	Article
	TagOptions []tags.Tag
}

// NewFormData creates form data with article and available tag options
func NewFormData(article Article, tagOptions []tags.Tag) FormData {
	return FormData{
		Article:    article,
		TagOptions: tagOptions,
	}
}
