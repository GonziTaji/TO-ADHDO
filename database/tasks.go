package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Tag struct {
	Id        string
	Name      string
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

type TaskTemplate struct {
	Id          string
	Name        string
	Description string
	Tags        []Tag
	CreatedAt   int
	UpdatedAt   int
	DeletedAt   int
}

type TaskId string
type TagId string
type TagName string

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
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM task_templates
		WHERE deleted_at IS NULL
	`, nil)

	if err != nil {
		// TODO: handle error
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
			tags.deleted_at,
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

		rows.Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt, &tag.DeletedAt, &task_id)

		task := tasks_by_id[task_id]
		task.Tags = append(task.Tags, tag)

		tasks_by_id[task_id] = task
	}

	tasks := make([]TaskTemplate, 0, len(tasks_by_id))
	for _, task := range tasks_by_id {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func CreateTaskTemplate(task_name string, task_description string, tag_names []string) (TaskId, error) {
	db, err := GetDatabase()

	if err != nil {
		return "", err
	}

	tx, err := db.Begin()

	if err != nil {
		fmt.Printf("error starting transaction: %s\n", err.Error())
		return "", err
	}

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
			SELECT name
			FROM input_tags
			WHERE name NOT IN (SELECT name from task_tags);
		`)

		log.Printf("select missing tags query: %s\n", query.String())

		rows, err := db.Query(query.String(), query_values...)

		if err != nil {
			fmt.Printf("error selecting missing tags: %s\n", err.Error())
			return "", err
		}

		var tags_to_create []any

		for rows.Next() {
			var tag_name string
			rows.Scan(&tag_name)
			tags_to_create = append(tags_to_create, tag_name)
		}

		if len(tags_to_create) > 0 {
			query := fmt.Sprintf(`
				INSERT INTO task_tags (name)
				VALUES %s
			`,
				strings.Repeat("(?),\n", len(tags_to_create)-1)+"(?)",
			)

			fmt.Printf("insert task_tags query: %s", query)

			res, err := tx.Exec(query, tags_to_create...)

			if err != nil {
				fmt.Printf("error inserting tags: %s\n", err.Error())
				return "", err
			}

			rows_affected, _ := res.RowsAffected()
			if rows_affected == 0 {
				return "", fmt.Errorf("Tags not inserted")
			}
		}
	}

	res, err := tx.Exec(`
		INSERT INTO task_templates (name, description)
		VALUES (?, ?)
	`, task_name, task_description)

	if err != nil {
		fmt.Printf("error inserting task template: %s\n", err.Error())
		return "", err
	}

	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		return "", fmt.Errorf("Task not inserted")
	}

	id, _ := res.LastInsertId()

	err = tx.Commit()
	if err != nil {
		fmt.Printf("error commiting transaction: %s\n", err.Error())
		return "", err
	}

	return TaskId(strconv.Itoa(int(id))), nil
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
