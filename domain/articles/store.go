package articles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/yogusita/to-adhdo/domain/tags"
	"github.com/yogusita/to-adhdo/domain/uploads"
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

	var articles_ids QueryArgs
	var articles_ids_placeholders []string
	var articles []Article

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
		log.Printf("failed to get articles: %s\n", err.Error())
		return nil, err
	}

	for rows.Next() {
		var article Article

		if err := rows.Scan(
			&article.Id,
			&article.Name,
			&article.Description,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.DeletedAt,
		); err != nil {
			log.Printf("failed to scan articles row: %s\n", err.Error())
			return nil, err
		}

		articles = append(articles, article)
		articles_ids = append(articles_ids, article.Id)
		articles_ids_placeholders = append(articles_ids_placeholders, "?")
	}

	if err := rows.Close(); err != nil {
		log.Printf("failed to close articles rows: %s\n", err.Error())
		return nil, err
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
		AND pivot.deleted_at IS NULL;
	`, strings.Join(articles_ids_placeholders, ", "))

	rows, err = s.db.Query(query, articles_ids...)

	if err != nil {
		log.Printf("failed to get articles_tags: %s\n", err.Error())
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
			log.Printf("failed to scan articles_tags row: %s\n", err.Error())
			return nil, err
		}

		article_idx := slices.IndexFunc(articles, func(a Article) bool { return a.Id == article_id })
		article := &articles[article_idx]
		article.Tags = append(article.Tags, tag)
	}

	if err = rows.Close(); err != nil {
		log.Printf("failed to close articles_tags rows: %s\n", err.Error())
		return nil, err
	}

	query = fmt.Sprintf(`
		SELECT
			id,
			article_id,
			price,
			description,
			created_at
		FROM articles_prices
		WHERE article_id IN (%s)
		ORDER BY created_at DESC;
		`, strings.Join(articles_ids_placeholders, ", "))

	rows, err = s.db.Query(query, articles_ids...)

	if err != nil {
		log.Printf("failed to get articles_prices: %s\n", err.Error())
		return nil, err
	}

	for rows.Next() {
		var price ArticlePrice

		if err := rows.Scan(
			&price.Id,
			&price.ArticleId,
			&price.Price,
			&price.Description,
			&price.CreatedAt,
		); err != nil {
			log.Printf("failed to scan articles_prices row: %s\n", err.Error())
			return nil, err
		}

		article_idx := slices.IndexFunc(articles, func(a Article) bool { return a.Id == price.ArticleId })
		article := &articles[article_idx]
		article.Prices = append(article.Prices, price)
	}

	if err = rows.Close(); err != nil {
		log.Printf("failed to close articles_prices rows: %s\n", err.Error())
		return nil, err
	}

	query = fmt.Sprintf(`
		SELECT id, article_id, filename
		FROM articles_images
		WHERE article_id IN (%s)
	`, strings.Join(articles_ids_placeholders, ", "))

	log.Printf("query: \n%s\n", query)

	rows, err = s.db.Query(query, articles_ids...)

	if err != nil {
		log.Printf("failed to get articles_images: %s\n", err.Error())
		return nil, err
	}

	for rows.Next() {
		image := ArticleImage{}

		if err := rows.Scan(
			&image.Id,
			&image.ArticleId,
			&image.Filename,
		); err != nil {
			log.Printf("failed to scan articles_images row: %s\n", err.Error())
			return nil, err
		}

		image.Url = uploads.GetFilePublicUrl(articles_images_bucket, image.Filename)

		article_idx := slices.IndexFunc(articles, func(a Article) bool { return a.Id == image.ArticleId })
		article := &articles[article_idx]
		article.Images = append(article.Images, image)
	}

	if err := rows.Close(); err != nil {
		log.Printf("failed to close articles_images rows: %s\n", err.Error())
		return nil, err
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
		AND pivot.deleted_at IS NULL
		AND articles.deleted_at IS NULL;
	`

	log.Printf("query: %s\n", query)

	rows, err := s.db.Query(query, article_id)

	if err != nil {
		return article, err
	}

	for rows.Next() {
		var tag tags.Tag

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

	if err := rows.Close(); err != nil {
		return article, err
	}

	prices_rows, err := s.db.Query(`
		SELECT id, article_id, price, description, created_at
		FROM articles_prices
		WHERE article_id = ?
		ORDER BY created_at DESC;
	`, article.Id)

	if err != nil {
		return article, err
	}

	for prices_rows.Next() {
		var price ArticlePrice

		if err := prices_rows.Scan(
			&price.Id,
			&price.ArticleId,
			&price.Price,
			&price.Description,
			&price.CreatedAt,
		); err != nil {
			return article, err
		}

		article.Prices = append(article.Prices, price)
	}

	if err := prices_rows.Close(); err != nil {
		return article, err
	}

	images_rows, err := s.db.Query(`
		SELECT id, article_id, filename, created_at
		FROM articles_images
		WHERE article_id = ?
		ORDER BY created_at DESC;
	`, article.Id)

	if err != nil {
		return article, err
	}

	for images_rows.Next() {
		var image ArticleImage

		if err := images_rows.Scan(
			&image.Id,
			&image.ArticleId,
			&image.Filename,
			&image.CreatedAt,
		); err != nil {
			return article, err
		}

		image.Url = uploads.GetFilePublicUrl(articles_images_bucket, image.Filename)

		article.Images = append(article.Images, image)
	}

	if err := images_rows.Close(); err != nil {
		return article, err
	}

	return article, nil
}

func (s *Store) Create(article *Article) (string, error) {
	tx, err := s.db.Begin()

	if err != nil {
		log.Printf("error starting transaction: %s\n", err.Error())
		return "", err
	}

	res, err := tx.Exec(`
		INSERT INTO articles (name, description)
		VALUES (?, ?)
	`, article.Name, article.Description)

	if err != nil {
		log.Printf("error inserting article: %s\n", err.Error())
		return "", err
	}

	i_article_id, _ := res.LastInsertId()
	article_id := dbIdToString(i_article_id)
	article.Id = article_id

	if len(article.Prices) > 0 {
		err := createPrices(tx, *article)

		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	var tags_ids []string
	var tags_names_to_create []string

	for _, tag := range article.Tags {
		if tag.Id != "" {
			tags_ids = append(tags_ids, tag.Id)
		} else {
			tags_names_to_create = append(tags_names_to_create, tag.Name)
		}
	}

	new_tags_ids, err := createTags(tx, tags_names_to_create)

	if err != nil {
		log.Printf("could not create tags: %s\n", err.Error())
		return "", err
	}

	tags_ids = slices.Concat(tags_ids, new_tags_ids)

	if err := createArticleTags(tx, article_id, tags_ids); err != nil {
		log.Printf("could not create articles_tags of new tags for article id %s: %s\n", article.Id, err.Error())
		tx.Rollback()
		return "", err
	}

	if err := persistArticleImages(tx, article); err != nil {
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("error commiting transaction: %s\n", err.Error())
		tx.Rollback()
		return "", err
	}

	return article_id, nil
}

// Updates the corresponding db tables with the data inside the article
//
// If a tag that is related to the article in the db, is not present in the struct, the relationship is deleted.
//
// Tags without id are treated like new tags. They are created and then related to the article
func (s *Store) Update(article *Article) error {
	tx, err := s.db.Begin()

	if err != nil {
		log.Printf("error starting transaction: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec(`
		UPDATE articles SET
			name = ?,
			description = ?,
			updated_at = current_timestamp
		WHERE id = ?;
	`, article.Name, article.Description, article.Id)

	if err != nil {
		log.Printf("error updating article: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	if len(article.Prices) > 0 {
		err := createPrices(tx, *article)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	input_tags_ids := make(map[string]string)
	var new_tags_names []string

	for _, tag := range article.Tags {
		if tag.Id == "" {
			// TODO: check if names are not repeated? or deal with this from the DB and handle the error?
			new_tags_names = append(new_tags_names, tag.Name)
		} else {
			input_tags_ids[tag.Id] = tag.Id
		}
	}

	tags_ids_for_new_relationships, err := createTags(tx, new_tags_names)

	if err != nil {
		log.Printf("error creating new tags for article: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	tags_ids_in_relationships_query := `
		SELECT tag_id
		FROM articles_tags
		WHERE article_id = ?
		AND deleted_at IS NULL
	`

	rows, err := tx.Query(tags_ids_in_relationships_query, article.Id)

	if err != nil {
		log.Printf("could not get articles_tags: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	tags_ids_in_relationships := make(map[string]string)

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			log.Printf("could not scan rows of articles_tags: %s\n", err.Error())
			tx.Rollback()
			return err
		}

		tags_ids_in_relationships[id] = id
	}

	if err := rows.Close(); err != nil {
		log.Printf("could not close rows of articles_tags of article: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	for _, tag_id := range input_tags_ids {
		if tags_ids_in_relationships[tag_id] == "" {
			tags_ids_for_new_relationships = append(tags_ids_for_new_relationships, tag_id)
		}
	}

	for _, tag_id := range tags_ids_in_relationships {
		if input_tags_ids[tag_id] == "" {
			_, err := tx.Exec(`
					UPDATE articles_tags
					SET
						updated_at = current_timestamp,
						deleted_at = current_timestamp
					WHERE article_id = ?
					and deleted_at IS NULL;
					`,
				article.Id, tag_id,
			)

			if err != nil {
				log.Printf("could not soft delete articles_tags for tag id: %s\n", tag_id)
				tx.Rollback()
				return err
			}
		}
	}

	if err := createArticleTags(tx, article.Id, tags_ids_for_new_relationships); err != nil {
		log.Printf("could not create articles_tags of new tags for article id %s: %s\n", article.Id, err.Error())
		tx.Rollback()
		return err
	}

	if err := persistArticleImages(tx, article); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("error commiting transaction: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (s *Store) Delete(article_id string) error {
	tx, err := s.db.Begin()

	if err != nil {
		log.Printf("Error starting transaction: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec(`
		UPDATE articles
		SET deleted_at = current_timestamp 
		WHERE id = ?
	`, article_id)

	if err != nil {
		log.Printf("Error soft-deleting article: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	relationships_tables := []string{
		"articles_images",
		"articles_tags",
	}

	for _, table := range relationships_tables {
		_, err = tx.Exec(`
			UPDATE articles_tags
			SET deleted_at = current_timestamp 
			WHERE article_id = ?
		`, article_id)

		if err != nil {
			log.Printf("Error soft-deleting article's %s: %s\n", table, err.Error())

			if e := tx.Rollback(); e != nil {
				return err
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting transaction: %s\n", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

// -------
// @review: move this to the tags store and call it from a service?

// @returns an array with the ids of the new tags
func createTags(tx *sql.Tx, tags_names []string) ([]string, error) {
	if len(tags_names) == 0 {
		return []string{}, nil
	}

	var insert_tags_query strings.Builder

	insert_tags_query.WriteString("INSERT INTO tags (name) VALUES")

	var new_tags_names QueryArgs

	for i, tag_name := range tags_names {
		new_tags_names = append(new_tags_names, tag_name)

		if i < len(tags_names)-1 {
			insert_tags_query.WriteString(" (?),")
		} else {
			insert_tags_query.WriteString(" (?) RETURNING id")
		}
	}

	var tags_ids []string

	rows, err := tx.Query(insert_tags_query.String(), new_tags_names...)

	if err != nil {
		return tags_ids, err
	}

	for rows.Next() {
		var tag_id string

		if err := rows.Scan(&tag_id); err != nil {
			return tags_ids, err
		}

		tags_ids = append(tags_ids, tag_id)
	}

	if e := rows.Close(); e != nil {
		log.Printf("Error closing rows: %s", e.Error())
		return tags_ids, e
	}

	return tags_ids, nil
}

func createArticleTags(tx *sql.Tx, article_id string, tags_ids []string) error {
	if len(tags_ids) == 0 {
		return nil
	}

	var query_sb strings.Builder
	var query_args QueryArgs

	query_sb.WriteString(
		"INSERT INTO articles_tags (article_id, tag_id) VALUES ",
	)

	for i, tag_id := range tags_ids {
		query_args = append(query_args, article_id, tag_id)

		if i < len(tags_ids)-1 {
			query_sb.WriteString(" (?, ?),")
		} else {
			query_sb.WriteString(" (?, ?);")
		}
	}

	res, err := tx.Exec(query_sb.String(), query_args...)

	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()

	if int(count) != len(query_args)/2 {
		return errors.New("One or more rows could not be inserted")
	}

	return nil
}

func dbIdToString(id int64) string {
	return strconv.Itoa(int(id))
}

func createPrices(tx *sql.Tx, article Article) error {
	var query_sb strings.Builder
	query_args := QueryArgs{}

	query_sb.WriteString(`
			INSERT INTO articles_prices (article_id, price, description)
			VALUES 
		`)

	for i, price := range article.Prices {
		query_args = append(query_args, article.Id, price.Price, price.Description)

		if i == 0 {
			query_sb.WriteString(" (?, ?, ?)")
		} else {
			query_sb.WriteString(", (?, ?, ?)")
		}
	}

	query_sb.WriteString(";")

	_, err := tx.Exec(query_sb.String(), query_args...)

	if err != nil {
		log.Printf("error creating prices: %s\n", err.Error())
		return err
	}

	return nil
}

func persistArticleImages(tx *sql.Tx, article *Article) error {
	if len(article.Images) == 0 {
		rows, err := tx.Query("SELECT filename FROM articles_images WHERE article_id = ?;", article.Id)

		for rows.Next() {
			var filename string

			if err := rows.Scan(&filename); err != nil {
				return err
			}

			if err := uploads.DeleteFile(articles_images_bucket, filename); err != nil {
				return err
			}
		}

		if err := rows.Close(); err != nil {
			return err
		}

		_, err = tx.Exec("DELETE FROM articles_images WHERE article_id = ?;", article.Id)

		return err
	}

	filenames_qargs := QueryArgs{}
	filenames_qargs_tmpls := []string{}

	for _, img := range article.Images {
		filenames_qargs = append(filenames_qargs, img.Filename)
		filenames_qargs_tmpls = append(filenames_qargs_tmpls, "?")
	}

	check_new_images_query := fmt.Sprintf(`
		WITH input_filenames(name) AS (
			SELECT %s
		)
		SELECT name
		FROM input_filenames
		WHERE name NOT IN (
			SELECT filename FROM articles_images WHERE article_id = ?
		);
	`, strings.Join(filenames_qargs_tmpls, " UNION ALL SELECT "))

	log.Print(check_new_images_query + "\n")

	rows, err := tx.Query(
		check_new_images_query,
		append(filenames_qargs, article.Id)...,
	)

	if err != nil {
		return err
	}

	inserted_count := 0

	for rows.Next() {
		filename := ""

		if err := rows.Scan(&filename); err != nil {
			return err
		}

		res, err := tx.Exec(
			"INSERT INTO articles_images (article_id, filename) VALUES (?, ?)",
			article.Id, filename,
		)

		if err != nil {
			return err
		}

		count, _ := res.RowsAffected()

		if count != 1 {
			return fmt.Errorf("Article image \"%s\" could not be inserted", filename)
		}

		inserted_count++
	}

	if err := rows.Close(); err != nil {
		return err
	}

	if inserted_count == len(article.Images) {
		// Every image has already been processed
		return nil
	}

	check_missing_images_query := fmt.Sprintf(`
		SELECT id, filename
		FROM articles_images
		WHERE article_id = ?
		AND filename NOT IN (%s);
	`, strings.Join(filenames_qargs_tmpls, ","))

	log.Print(check_missing_images_query + "\n")

	rows, err = tx.Query(
		check_missing_images_query,
		slices.Concat(QueryArgs{article.Id}, filenames_qargs)...,
	)

	if err != nil {
		return err
	}

	for rows.Next() {
		id := ""
		filename := ""

		if err := rows.Scan(&id, &filename); err != nil {
			return err
		}

		res, err := tx.Exec("DELETE FROM articles_images WHERE id = ?", id)

		if err != nil {
			return err
		}

		count, _ := res.RowsAffected()

		if count != 1 {
			return fmt.Errorf("Article image of id %s could not be deleted", id)
		}

		if err := uploads.DeleteFile(articles_images_bucket, filename); err != nil {
			return err
		}
	}

	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}
