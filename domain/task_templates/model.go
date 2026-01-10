package task_templates

import "github.com/yogusita/to-adhdo/domain/tags"

type TaskTemplate struct {
	Id          string
	Name        string
	Description string
	Tags        []tags.Tag
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
