package tags

import (
	"database/sql"
	"log"
	"strings"
)

type Store struct {
	db *sql.DB
}

func CreateStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) List(options ListingTagsOptions) ([]Tag, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}

	var sb_query strings.Builder

	sb_query.WriteString(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM tags
	`)

	log.Printf("\n\ninclude deleted: %v\n\n", options.IncludeDeleted)

	if !options.IncludeDeleted {
		sb_query.WriteString(" WHERE deleted_at IS NULL ")
	}

	sb_query.WriteString(" LIMIT ? OFFSET ? ;")

	log.Printf("query: %s", sb_query.String())

	rows, err := s.db.Query(sb_query.String(), options.Limit, options.Offset)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	var tags []Tag
	for rows.Next() {
		var tag Tag

		rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.DeletedAt,
		)

		tags = append(tags, tag)
	}

	if err := rows.Close(); err != nil {
		return tags, err
	}

	return tags, nil
}
