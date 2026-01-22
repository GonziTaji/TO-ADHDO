package tags

import "database/sql"

type Tag struct {
	Id        string
	Name      string
	CreatedAt string
	UpdatedAt string
	DeletedAt sql.NullString
}
