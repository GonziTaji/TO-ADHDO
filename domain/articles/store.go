package articles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/yogusita/to-adhdo/domain/tags"
)

type QueryArgs []any

type Store struct {
	db *sql.DB
}

func CreateStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) List(options *ListingArticlesOptions) ([]Article, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}

	articles_by_id := make(map[string]Article)
	var articles_ids []any // type any to use it as query parameter

	var query_sb strings.Builder

	query_sb.WriteString(`
		SELECT
			id,
			name,
			COALESCE(description, ''),
			created_at,
			updated_at,
			COALESCE(deleted_at, '')
		FROM articles
	`)

	if !options.IncludeDeleted {
		query_sb.WriteString(" WHERE deleted_at IS NULL")
	}

	query_sb.WriteString(" ORDER BY updated_at DESC")
	query_sb.WriteString(" LIMIT ? OFFSET ?;")

	rows, err := s.db.Query(query_sb.String(), options.Limit, options.Offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var article Article

		rows.Scan(
			&article.Id,
			&article.Name,
			&article.Description,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.DeletedAt,
		)

		articles_by_id[article.Id] = article
		articles_ids = append(articles_ids, article.Id)
	}

	articles_count := len(articles_ids)

	query_placeholders := make([]string, articles_count)

	for i := range articles_count {
		query_placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
		SELECT
			tags.id,
			tags.name,
			tags.created_at,
			tags.updated_at,
			COALESCE(tags.deleted_at, ''),
			pivot.article_id as article_id
		FROM articles_tags as pivot
		JOIN tags on tags.id = pivot.tag_id
		WHERE pivot.article_id IN (%s)
	`, strings.Join(query_placeholders, ","))

	rows, err = s.db.Query(query, articles_ids...)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	for rows.Next() {
		var tag tags.Tag
		var article_id string

		if err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.DeletedAt,
			&article_id,
		); err != nil {
			return nil, err
		}

		article := articles_by_id[article_id]
		article.Tags = append(article.Tags, tag)

		articles_by_id[article_id] = article
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	articles := make([]Article, 0, len(articles_by_id))
	for _, article := range articles_by_id {
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *Store) Get(article_id string) (Article, error) {
	article := Article{}

	query := `
		SELECT
			COALESCE(tags.id, '') as tag_id,
			COALESCE(tags.name, '') as tag_id,
			articles.id,
			articles.name,
			COALESCE(articles.description, ''),
			articles.created_at,
			articles.updated_at,
			COALESCE(articles.deleted_at, '')
		FROM articles 
		LEFT JOIN articles_tags as pivot on articles.id = pivot.article_id
		LEFT JOIN tags on pivot.tag_id = tags.id
		WHERE articles.id = ?
	`

	rows, err := s.db.Query(query, article_id)

	if err != nil {
		return article, err
	}

	for rows.Next() {
		var tag tags.Tag

		fmt.Println("in if")
		err = rows.Scan(
			&tag.Id,
			&tag.Name,
			&article.Id,
			&article.Name,
			&article.Description,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.DeletedAt,
		)

		if err != nil {
			return article, err
		}

		if tag.Id != "" {
			article.Tags = append(article.Tags, tag)
		}
	}

	return article, nil
}

func (s *Store) Create(article *Article) (string, error) {
	log.Printf("creating article: %v\n", article)

	var new_tags_names QueryArgs

	tx, err := s.db.Begin()

	if err != nil {
		fmt.Printf("error starting transaction: %s\n", err.Error())
		return "", err
	}

	res, err := tx.Exec(`
		INSERT INTO articles (name, description)
		VALUES (?, ?)
	`, article.Name, article.Description)

	if err != nil {
		fmt.Printf("error inserting article: %s\n", err.Error())
		return "", err
	}

	i_article_id, _ := res.LastInsertId()
	article_id := strconv.Itoa(int(i_article_id))

	if len(article.Tags) > 0 {
		var insert_tags_query strings.Builder

		insert_tags_query.WriteString("INSERT INTO tags (name) VALUES")

		var tag_ids []string

		for i, tag := range article.Tags {
			if tag.Id == "" {
				new_tags_names = append(new_tags_names, tag.Name)

				if i < len(article.Tags)-1 {
					insert_tags_query.WriteString(" (?),")
				} else {
					insert_tags_query.WriteString(" (?) RETURNING id")
				}
			} else {
				fmt.Printf("inside else\n")
				tag_ids = append(tag_ids, tag.Id)
			}
		}

		if len(new_tags_names) > 0 {
			log.Printf("query: \n%s\n", insert_tags_query.String())

			rows, err := s.db.Query(insert_tags_query.String(), new_tags_names...)

			if err != nil {
				return "", err
			}

			for rows.Next() {
				var id string

				if err := rows.Scan(&id); err != nil {
					return "", err
				}

				tag_ids = append(tag_ids, id)
			}

			if e := rows.Close(); e != nil {
				fmt.Printf("Error closing rows: %s", e.Error())
				return "", e
			}
		}

		if len(tag_ids) > 0 {
			var query_sb strings.Builder
			var insert_pivot_args QueryArgs

			query_sb.WriteString("INSERT INTO articles_tags (tag_id, article_id) VALUES")

			for i, tag_id := range tag_ids {
				fmt.Printf("in for tag_ids. tag_id: %s, article_id: %s", tag_id, article_id)
				insert_pivot_args = append(insert_pivot_args, tag_id, article_id)

				if i < len(tag_ids)-1 {
					query_sb.WriteString(" (?, ?),")
				} else {
					query_sb.WriteString(" (?, ?)")
				}
			}

			log.Printf("query: \n%s\n", query_sb.String())

			res, err = tx.Exec(query_sb.String(), insert_pivot_args...)

			if err != nil {
				tx.Rollback()
				fmt.Printf("error inserting pivot table: %s\n", err.Error())
				return "", err
			}

			count, _ := res.RowsAffected()

			log.Printf("count: %d, pivot_args: %d", count, len(insert_pivot_args))

			if int(count) != len(insert_pivot_args)/2 {
				tx.Rollback()
				return "", errors.New("One or more rows could not be inserted")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("error commiting transaction: %s\n", err.Error())

		return "", err
	}

	return article_id, nil
}

func (s *Store) Update(article *Article) error {
	log.Printf("updating article: %v\n", article)
	var new_tags_names QueryArgs

	tx, err := s.db.Begin()

	if err != nil {
		fmt.Printf("error starting transaction: %s\n", err.Error())
		return err
	}

	res, err := tx.Exec(`
		UPDATE articles SET
			name = ?,
			description = ?,
			updated_at = current_timestamp
		WHERE id = ?;
	`, article.Name, article.Description, article.Id)

	if err != nil {
		fmt.Printf("error updating article: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	var get_tags_query strings.Builder
	var insert_tags_query strings.Builder

	get_tags_query.WriteString("SELECT name FROM (")

	insert_tags_query.WriteString("INSERT INTO tags (name) VALUES")

	input_tags_ids := make(map[string]string)

	var tag_ids_of_new_relationships []string
	var tag_ids_of_dead_relationships []string

	for _, tag := range article.Tags {
		if tag.Id == "" {
			// TODO: check if names are not repeated? or deal with this from the DB and handle the error?
			new_tags_names = append(new_tags_names, tag.Name)

			if len(new_tags_names) == 1 {
				insert_tags_query.WriteString(" (?)")
			} else {
				insert_tags_query.WriteString(", (?)")
			}
		} else {
			input_tags_ids[tag.Id] = tag.Id

			if len(input_tags_ids) == 1 {
				get_tags_query.WriteString(" SELECT ? AS name\n")
			} else {
				get_tags_query.WriteString(" UNION ALL SELECT ?\n")
			}
		}
	}

	if len(new_tags_names) > 0 {
		insert_tags_query.WriteString(" RETURNING id;")

		log.Printf("insert_tags_query: %s\n", insert_tags_query.String())

		rows, err := s.db.Query(insert_tags_query.String(), new_tags_names...)

		if err != nil {
			log.Printf("could not insert tags: %s\n", err.Error())
			return err
		}

		for rows.Next() {
			var id string

			if err := rows.Scan(&id); err != nil {
				log.Printf("error scanning row of inserted tag: %s\n", err.Error())
				return err
			}

			tag_ids_of_new_relationships = append(tag_ids_of_new_relationships, id)
		}

		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows of inserted tags: %s", err.Error())
			return err
		}
	}

	if len(input_tags_ids) > 0 {
		get_tags_query.WriteString(
			") input WHERE id NOT IN (select tag_id FROM articles_tags WHERE article_id = ? and deleted_at IS NULL)",
		)

		args := make(QueryArgs, len(input_tags_ids)+1)
		for _, id := range input_tags_ids {
			args = append(args, id)
		}
		args = append(args, article.Id)

		log.Printf("query to get tags: %s\n", get_tags_query.String())

		rows, err := s.db.Query(get_tags_query.String(), args...)

		if err != nil {
			log.Printf("error executing query to get existing related tags: %s", err.Error())
			return err
		}

		existing_related_tags_ids := make(map[string]string)

		for rows.Next() {
			var id string
			rows.Scan(&id)

			existing_related_tags_ids[id] = id
		}

		if err := rows.Close(); err != nil {
			log.Printf("error closing rows of existing related tags: %s", err.Error())
			return err
		}

		for _, existing_tag_id := range input_tags_ids {
			if existing_related_tags_ids[existing_tag_id] == "" {
				tag_ids_of_new_relationships = append(tag_ids_of_new_relationships, existing_tag_id)
			}
		}

		for _, related_tag_id := range existing_related_tags_ids {
			if input_tags_ids[related_tag_id] == "" {
				tag_ids_of_dead_relationships = append(tag_ids_of_dead_relationships, related_tag_id)
			}
		}
	}

	if len(tag_ids_of_new_relationships) > 0 {
		var query_sb strings.Builder
		var insert_pivot_args QueryArgs

		query_sb.WriteString("INSERT INTO articles_tags (tag_id, article_id) VALUES")

		for i, tag_id := range tag_ids_of_new_relationships {
			insert_pivot_args = append(insert_pivot_args, tag_id, article.Id)

			if i < len(tag_ids_of_new_relationships)-1 {
				query_sb.WriteString(" (?, ?),")
			} else {
				query_sb.WriteString(" (?, ?)")
			}
		}

		log.Printf("query: \n%s\n", query_sb.String())

		res, err = tx.Exec(query_sb.String(), insert_pivot_args...)

		if err != nil {
			tx.Rollback()
			fmt.Printf("error inserting pivot table: %s\n", err.Error())
			return err
		}

		count, _ := res.RowsAffected()

		log.Printf("count: %d, pivot_args: %d", count, len(insert_pivot_args))

		if int(count) != len(insert_pivot_args)/2 {
			tx.Rollback()
			return errors.New("One or more rows could not be inserted")
		}
	}

	if len(tag_ids_of_dead_relationships) > 0 {
		var query_sb strings.Builder
		var delete_pivot_args QueryArgs

		query_sb.WriteString("UPDATE articles_tags SET\n")
		query_sb.WriteString("deleted_at = current_timestamp, updated_at = current_timestamp\n")
		query_sb.WriteString("WHERE aticle_id = ? and tag_id in (")

		delete_pivot_args = append(delete_pivot_args, article.Id)

		for i, tag_id := range tag_ids_of_dead_relationships {
			delete_pivot_args = append(delete_pivot_args, tag_id)

			if i == 0 {
				query_sb.WriteString(" ?")
			} else if i < len(tag_ids_of_dead_relationships)-1 {
				query_sb.WriteString(", ?")
			} else {
				query_sb.WriteString(", ?)")
			}
		}

		log.Printf("query: \n%s\n", query_sb.String())

		res, err = tx.Exec(query_sb.String(), delete_pivot_args...)

		if err != nil {
			tx.Rollback()
			fmt.Printf("error deleting pivot table: %s\n", err.Error())
			return err
		}

		count, _ := res.RowsAffected()

		log.Printf("affected rows: %d, pivot_args: %d", count, len(delete_pivot_args))

		if int(count) != len(delete_pivot_args)/2 {
			tx.Rollback()
			return errors.New("One or more rows could not be deleted")
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("error commiting transaction: %s\n", err.Error())

		return err
	}

	return nil
}

func (s *Store) Delete(article_id string) error {
	tx, err := s.db.Begin()

	if err != nil {
		fmt.Printf("Error starting transaction: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec(`
		UPDATE articles
		SET deleted_at = current_timestamp 
		WHERE id = ?
	`, article_id)

	if err != nil {
		fmt.Printf("Error soft-deleting article: %s\n", err.Error())
		return err
	}

	relationships_tables := []string{
		"articles_images",
		"articles_tags",
	}

	for _, table := range relationships_tables {
		fmt.Printf("soft-deleting article's %s\n", table)
		_, err = tx.Exec(`
			UPDATE articles_tags
			SET deleted_at = current_timestamp 
			WHERE article_id = ?
		`, article_id)

		if err != nil {
			fmt.Printf("Error soft-deleting article's %s: %s\n", table, err.Error())

			if e := tx.Rollback(); e != nil {
				return err
			}

			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		fmt.Printf("Error commiting transaction: %s\n", err.Error())
		return err
	}

	return nil
}
