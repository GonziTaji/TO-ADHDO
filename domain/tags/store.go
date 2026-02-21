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

func (s *Store) List(options ListingTagsOptions) ([]TagItemList, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}

	var sb_query strings.Builder

	sb_query.WriteString(`
		SELECT t.id, t.name, t.created_at, COUNT(at.article_id) as usage
		FROM tags t
		LEFT JOIN articles_tags at ON t.id = at.tag_id
	`)

	sb_query.WriteString(`
		GROUP BY t.id, t.name, t.created_at
		ORDER BY t.created_at DESC
		LIMIT 100 OFFSET 0;
	`)

	log.Printf("query: %s", sb_query.String())

	rows, err := s.db.Query(sb_query.String(), options.Limit, options.Offset)

	if err != nil {
		log.Printf("could not get task list: %s\n", err.Error())
		return nil, err
	}

	var list []TagItemList
	for rows.Next() {
		var tag_data TagItemList

		if err := rows.Scan(
			&tag_data.Id,
			&tag_data.Name,
			&tag_data.CreatedAt,
			&tag_data.Usage,
		); err != nil {
			log.Printf("could not scan task row: %s\n", err.Error())
			return nil, err
		}

		log.Printf("list item > %v\n", tag_data)

		list = append(list, tag_data)
	}

	if err := rows.Close(); err != nil {
		log.Printf("could not close rows: %s\n", err.Error())
		return nil, err
	}

	return list, nil
}
