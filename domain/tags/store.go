package tags

import "github.com/yogusita/to-adhdo/database"

type Store struct {
}

func (Store) List(limit int8, include_deleted bool) ([]Category, error) {
	if limit == 0 {
		limit = 10
	}

	db, err := database.GetDatabase()

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	var tags []Category

	rows, err := db.Query(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM task_tags
		WHERE deleted_at IS NULL
		LIMIT ?
	`, limit)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	for rows.Next() {
		var tag Category

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
