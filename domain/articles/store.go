package articles

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yogusita/to-adhdo/database"
	categories "github.com/yogusita/to-adhdo/domain/tags"
)

type QueryArgs []any

type Store struct {
}

func (Store) List(limit int8, include_deleted bool) ([]Article, error) {
	if limit == 0 {
		limit = 10
	}

	articles_by_id := make(map[string]Article)
	var articles_ids []any // type any to use it as query parameter

	db, err := database.GetDatabase()

	if err != nil {
		// TODO: handle error
		return nil, err
	}

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

	if include_deleted == false {
		query_sb.WriteString(" WHERE deleted_at IS NULL")
	}

	query_sb.WriteString(" ORDER BY updated_at DESC;")

	rows, err := db.Query(query_sb.String())

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
			categories.id,
			categories.name,
			categories.created_at,
			categories.updated_at,
			COALESCE(categories.deleted_at, ''),
			pivot.article_id as article_id
		FROM articles_categories as pivot
		JOIN categories on categories.id = pivot.category_id
		WHERE pivot.article_id IN (%s)
	`, strings.Join(query_placeholders, ","))

	rows, err = db.Query(query, articles_ids...)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	for rows.Next() {
		var category categories.Category
		var article_id string

		if err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
			&article_id,
		); err != nil {
			return nil, err
		}

		article := articles_by_id[article_id]
		article.Tags = append(article.Tags, category)

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

func (Store) Get(article_id string) (Article, error) {
	article := Article{}

	db, err := database.GetDatabase()

	if err != nil {
		return article, err
	}

	query := `
		SELECT
			COALESCE(categories.id, '') as category_id,
			COALESCE(categories.name, '') as category_id,
			articles.id,
			articles.name,
			COALESCE(articles.description, ''),
			articles.created_at,
			articles.updated_at,
			COALESCE(articles.deleted_at, '')
		FROM articles 
		LEFT JOIN articles_categories as pivot on articles.id = pivot.article_id
		LEFT JOIN categories on pivot.category_id = categories.id
		WHERE articles.id = ?
	`

	rows, err := db.Query(query, article_id)

	if err != nil {
		return article, err
	}

	for rows.Next() {
		var category categories.Category

		fmt.Println("in if")
		err = rows.Scan(
			&category.Id,
			&category.Name,
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

		if category.Id != "" {
			article.Tags = append(article.Tags, category)
		}
	}

	return article, nil
}

func (Store) Create(article *Article) (string, error) {
	var new_categories_names QueryArgs

	db, err := database.GetDatabase()

	if err != nil {
		return "", err
	}

	tx, err := db.Begin()

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

	var insert_categories_query strings.Builder

	insert_categories_query.WriteString("INSERT INTO categories (name) VALUES")

	for i, category := range article.Categories {
		if category.Id == "" {
			new_categories_names = append(new_categories_names, category.Name)

			if i < len(article.Categories)-1 {
				insert_categories_query.WriteString(" (?),")
			} else {
				insert_categories_query.WriteString(" (?) RETURNING id")
			}
		}
	}

	rows, err := db.Query(insert_categories_query.String(), new_categories_names...)

	for rows.Next() {
		var new_category Category

		if err := rows.Scan(
			&new_category.Id,
			&new_category.Name,
			&new_category.CreatedAt,
			&new_category.UpdatedAt,
		); err != nil {
			return "", err
		}

		for i, c := range article.Categories {
			if c.Name == new_category.Name {
				article.Categories[i] = new_category
				break
			}
		}
	}

	if e := rows.Close(); e != nil {
		fmt.Printf("Error closing rows: ", e.Error())
		return "", e
	}

	var query_sb strings.Builder
	var insert_pivot_args QueryArgs

	query_sb.WriteString("INSERT INTO articles_categories (category_id, article_id) VALUES")

	for i, category_id := range article.Categories {
		insert_pivot_args = append(insert_pivot_args, category_id, article_id)

		if i < len(article.Categories)-1 {
			query_sb.WriteString(" (?, ?),")
		} else {
			query_sb.WriteString(" (?, ?)")
		}
	}

	res, err = tx.Exec(query_sb.String(), insert_pivot_args...)

	if err != nil {
		tx.Rollback()
		fmt.Printf("error inserting pivot table: %s\n", err.Error())
		return "", err
	}

	count, _ := res.RowsAffected()

	if int(count) != len(article.Categories) {
		tx.Rollback()
		return "", errors.New("One or more rows could not be inserted")
	}

	if e := tx.Commit(); e != nil {
		fmt.Printf("error commiting transaction: %s\n", err.Error())
		return "", err
	}

	return article_id, nil
}

func (Store) Delete(article_id string) error {
	db, err := database.GetDatabase()

	if err != nil {
		return err
	}

	tx, err := db.Begin()

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
		"articles_categories",
	}

	for _, table := range relationships_tables {
		_, err = tx.Exec(`
			UPDATE articles_categories
			SET deleted_at = current_timestamp 
			WHERE article_id = ?
		`, article_id)

		if err != nil {
			fmt.Printf("Error soft-deleting article's %s: %s\n", table, err.Error())
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
