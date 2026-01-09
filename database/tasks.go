package database

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Tag struct {
	Id        string
	Name      string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

type TaskTemplate struct {
	Id          string
	Name        string
	Description string
	Tags        []Tag
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}

func GetAvailableTaskTemplates(limit int8) ([]TaskTemplate, error) {
	if limit == 0 {
		limit = 10
	}

	tasks_by_id := make(map[string]TaskTemplate)
	var tasks_ids []any // type any to use it as query parameter

	db, err := GetDatabase()

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			id,
			name,
			COALESCE(description, ''),
			created_at,
			updated_at,
			COALESCE(deleted_at, '')
		FROM task_templates
		WHERE deleted_at IS NULL
		ORDER BY updated_at DESC
	`, nil)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task TaskTemplate

		rows.Scan(
			&task.Id,
			&task.Name,
			&task.Description,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
		)

		tasks_by_id[task.Id] = task
		tasks_ids = append(tasks_ids, task.Id)
	}

	tasks_count := len(tasks_ids)

	query_placeholders := make([]string, tasks_count)

	for i := range tasks_count {
		query_placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
		SELECT
			tags.id,
			tags.name,
			tags.created_at,
			tags.updated_at,
			COALESCE(tags.deleted_at, ''),
			pivot.task_template_id as task_id
		FROM task_template_task_tags as pivot
		JOIN task_tags as tags on tags.id = pivot.task_tag_id
		WHERE pivot.task_template_id IN (%s)
	`, strings.Join(query_placeholders, ","))

	rows, err = db.Query(query, tasks_ids...)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	for rows.Next() {
		var tag Tag
		var task_id string

		if err := rows.Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt, &tag.DeletedAt, &task_id); err != nil {
			return nil, err
		}

		task := tasks_by_id[task_id]
		task.Tags = append(task.Tags, tag)

		tasks_by_id[task_id] = task
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	tasks := make([]TaskTemplate, 0, len(tasks_by_id))
	for _, task := range tasks_by_id {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskTemplate(task_id string) (TaskTemplate, error) {
	task_template := TaskTemplate{}

	db, err := GetDatabase()

	if err != nil {
		return task_template, err
	}

	query := `
		SELECT
			COALESCE(task_tags.id, '') as tag_id,
			COALESCE(task_tags.name, '') as tag_name,
			task_templates.id,
			task_templates.name,
			COALESCE(task_templates.description, ''),
			task_templates.created_at,
			task_templates.updated_at,
			COALESCE(task_templates.deleted_at, '')
		FROM task_templates 
		LEFT JOIN task_template_task_tags as pivot on task_templates.id = pivot.task_template_id
		LEFT JOIN task_tags on pivot.task_tag_id = task_tags.id
		WHERE task_templates.id = ?
	`

	log.Printf("Executing query to get task and tags info: \n\n%s\n\nWith task id: %s", query, task_id)

	rows, err := db.Query(query, task_id)

	if err != nil {
		return task_template, err
	}

	for rows.Next() {
		var new_tag Tag

		fmt.Println("in if")
		err = rows.Scan(
			&new_tag.Id,
			&new_tag.Name,
			&task_template.Id,
			&task_template.Name,
			&task_template.Description,
			&task_template.CreatedAt,
			&task_template.UpdatedAt,
			&task_template.DeletedAt,
		)

		if err != nil {
			return task_template, err
		}

		if new_tag.Id != "" {
			task_template.Tags = append(task_template.Tags, new_tag)
		}
	}

	return task_template, nil
}

// Returns the id of the new task_template
func CreateTaskTemplate(task_name string, task_description string, tag_names []string) (string, error) {
	db, err := GetDatabase()

	if err != nil {
		return "", err
	}

	tx, err := db.Begin()

	if err != nil {
		fmt.Printf("error starting transaction: %s\n", err.Error())
		return "", err
	}

	res, err := tx.Exec(`
		INSERT INTO task_templates (name, description)
		VALUES (?, ?)
	`, task_name, task_description)

	if err != nil {
		fmt.Printf("error inserting task template: %s\n", err.Error())
		return "", err
	}

	i_task_id, _ := res.LastInsertId()
	task_template_id := strconv.Itoa(int(i_task_id))

	if len(tag_names) > 0 {
		var query_values []any
		for _, name := range tag_names {
			query_values = append(query_values, name)
		}

		var query strings.Builder

		query.WriteString(`
			WITH input_tags(name) AS (
			SELECT ?
		`)

		for range len(query_values) - 1 {
			query.WriteString("UNION SELECT ?\n")
		}

		query.WriteString(")\n")

		query.WriteString(`
			SELECT 
				input_tags.name,
				task_tags.id
			FROM input_tags
			LEFT JOIN task_tags on input_tags.name = task_tags.name;
		`)

		log.Printf("select missing tags query: %s\n", query.String())

		rows, err := db.Query(query.String(), query_values...)

		if err != nil {
			fmt.Printf("Error getting tags info: %s\n", err.Error())
			tx.Rollback()

			return "", err
		}

		var tags_to_create []any
		var task_tags_ids []string

		for rows.Next() {
			var tag_name string
			var tag_id string
			rows.Scan(&tag_name, &tag_id)

			fmt.Printf("tag found: id='%s' name='%s'\n", tag_id, tag_name)

			if tag_id == "" {
				tags_to_create = append(tags_to_create, tag_name)
			} else {
				task_tags_ids = append(task_tags_ids, tag_id)
			}
		}

		if len(tags_to_create) > 0 {
			for _, new_tag_name := range tags_to_create {
				res, err := tx.Exec("INSERT INTO task_tags (name) VALUES (?)", new_tag_name)

				if err != nil {
					tx.Rollback()
					fmt.Printf("error inserting named %s: %s\n", new_tag_name, err.Error())
					return "", err
				}

				new_tag_id, _ := res.LastInsertId()

				task_tags_ids = append(task_tags_ids, fmt.Sprintf("%d", new_tag_id))
			}
		}

		var query_sb strings.Builder
		query_params := make([]any, 0, len(task_tags_ids)*2)

		query_sb.WriteString("INSERT INTO task_template_task_tags (task_tag_id, task_template_id) VALUES")

		for i, task_tag_id := range task_tags_ids {
			query_sb.WriteString(" (?, ?)")
			query_params = append(query_params, task_tag_id, task_template_id)

			if i < len(task_tags_ids)-1 {
				query_sb.WriteString(",")
			}
		}

		res, err := tx.Exec(query_sb.String(), query_params...)

		if err != nil {
			tx.Rollback()
			fmt.Printf("error inserting pivot table: %s\n", err.Error())
			return "", err
		}

		count, _ := res.RowsAffected()

		if int(count) != len(task_tags_ids) {
			tx.Rollback()
			return "", errors.New("One or more rows could not be inserted")
		}
	}

	err = tx.Commit()

	if err != nil {
		fmt.Printf("error commiting transaction: %s\n", err.Error())
		return "", err
	}

	return task_template_id, nil
}

func DeleteTaskTemplate(task_id string) error {
	db, err := GetDatabase()

	if err != nil {
		return err
	}

	tx, err := db.Begin()

	if err != nil {
		fmt.Printf("Error starting transaction: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec(`
		UPDATE task_templates
		SET deleted_at = current_timestamp 
		WHERE id = ?
	`, task_id)

	if err != nil {
		fmt.Printf("Error soft-deleting task_template: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec(`
		UPDATE task_template_task_tags
		SET deleted_at = current_timestamp 
		WHERE task_template_id = ?
	`, task_id)

	err = tx.Commit()

	if err != nil {
		fmt.Printf("Error commiting transaction: %s\n", err.Error())
		return err
	}

	return nil
}
