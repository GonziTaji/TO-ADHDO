package database

func GetAvailableTags(limit int8) ([]Tag, error) {
	if limit == 0 {
		limit = 10
	}

	db, err := GetDatabase()

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	var tags []Tag

	rows, err := db.Query(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM task_tags
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

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

	return tags, nil
}
