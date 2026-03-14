package articles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/yogusita/to-adhdo/domain/articles/model"
	"github.com/yogusita/to-adhdo/domain/shared"
	"github.com/yogusita/to-adhdo/domain/tags"
	"github.com/yogusita/to-adhdo/domain/uploads"
)

type Store struct {
	db *sql.DB
}

func CreateStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) WithTx(fn func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()

	if err != nil {
		log.Printf("error starting transaction: %s\n", err.Error())
		return err
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback failed: %w", rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("commit failed: %w", err)
		}

		return err
	}

	return nil
}

func (s *Store) Catalog(options model.CatalogFilterOptions) ([]model.CatalogItem, error) {
	var queryargs shared.QueryArgs
	var query_sb strings.Builder

	query_sb.WriteString(`
		SELECT
			a.id,
			a.name,
			COALESCE(t.name, '') as tag_name,
			COALESCE(p.price, 0),
			COALESCE(tn.filename, 'no-thumbnail.png'),
			a.available_for_trade,
			a.reference_price

		FROM articles a

		LEFT JOIN articles_tags at
			ON at.article_id = a.id

		LEFT JOIN tags t
			ON t.id = at.tag_id

		LEFT JOIN (
			SELECT article_id, filename, MAX(created_at)
			FROM articles_images
			GROUP BY article_id
		) tn ON tn.article_id = a.id

		LEFT JOIN (
			SELECT article_id, price, MAX(created_at)
			FROM articles_prices
			GROUP BY article_id
		) p ON p.article_id = a.id

		WHERE a.deleted_at IS NULL
	`)

	if len(options.SearchTerm) > 0 {
		queryargs = append(queryargs, "%"+options.SearchTerm+"%")
		query_sb.WriteString("\nAND a.name LIKE ?")
	}

	if len(options.TagsIdsFilter) > 0 {
		query_sb.WriteString("\nAND t.id IN (")
		for i, tagid := range options.TagsIdsFilter {
			if i == 0 {
				query_sb.WriteString("?")
			} else {
				query_sb.WriteString(",?")
			}

			log.Printf("found tag id for filter: %s\n", tagid)

			queryargs = append(queryargs, tagid)
		}
		query_sb.WriteString(")")
	}

	queryargs = append(queryargs, 100)
	query_sb.WriteString("\nLIMIT ?;")

	items_by_id := make(map[string]model.CatalogItem)

	log.Printf("query: \n%s\n", query_sb.String())
	log.Printf("\nargs: \n")
	for _, arg := range queryargs {
		log.Printf("- %v\n", arg)
	}

	rows, err := s.db.Query(query_sb.String(), queryargs...)

	if err != nil {
		log.Printf("error getting catalog items: %s\n", err.Error())
		return nil, err
	}

	for rows.Next() {
		log.Println("in article row for catalog")
		var scanned_item model.CatalogItem
		var item_tag struct{ Name string }

		if err := rows.Scan(
			&scanned_item.Id,
			&scanned_item.Name,
			&item_tag.Name,
			&scanned_item.Price,
			&scanned_item.ThumbnailUrl,
			&scanned_item.AvailableForTrade,
			&scanned_item.ReferencePrice,
		); err != nil {
			log.Printf("error scanning catalog item data from db %s\n", err.Error())
			return nil, err
		}

		if items_by_id[scanned_item.Id].Id == "" {
			scanned_item.ThumbnailUrl = uploads.GetFilePublicUrl(articles_images_bucket, scanned_item.ThumbnailUrl)
			items_by_id[scanned_item.Id] = scanned_item
		}

		if item_tag.Name != "" {
			item := items_by_id[scanned_item.Id]
			item.Tags = append(item.Tags, item_tag)
			items_by_id[scanned_item.Id] = item
		}
	}

	items := slices.Collect(maps.Values(items_by_id))
	log.Printf("\n\nitems: %v\n\n", items)

	return items, nil
}

func (s *Store) List(options *ListingArticlesOptions) ([]model.Article, error) {
	if options.Limit == 0 {
		options.Limit = 100
	}

	var articles_ids shared.QueryArgs
	var articles_ids_placeholders []string
	var articles []model.Article

	var query_sb strings.Builder

	query_sb.WriteString(`
		SELECT id,
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
		var article model.Article

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
			at.article_id as article_id
		FROM articles_tags as at
		JOIN tags on tags.id = at.tag_id
		WHERE at.article_id IN (%s);
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
			&article_id,
		); err != nil {
			log.Printf("failed to scan articles_tags row: %s\n", err.Error())
			return nil, err
		}

		article_idx := slices.IndexFunc(articles, func(a model.Article) bool { return a.Id == article_id })
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
		var price model.ArticlePrice

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

		article_idx := slices.IndexFunc(articles, func(a model.Article) bool { return a.Id == price.ArticleId })
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
		image := model.ArticleImage{}

		if err := rows.Scan(
			&image.Id,
			&image.ArticleId,
			&image.Filename,
		); err != nil {
			log.Printf("failed to scan articles_images row: %s\n", err.Error())
			return nil, err
		}

		image.Url = uploads.GetFilePublicUrl(articles_images_bucket, image.Filename)

		article_idx := slices.IndexFunc(articles, func(a model.Article) bool { return a.Id == image.ArticleId })
		article := &articles[article_idx]
		article.Images = append(article.Images, image)
	}

	if err := rows.Close(); err != nil {
		log.Printf("failed to close articles_images rows: %s\n", err.Error())
		return nil, err
	}

	return articles, nil
}

func (s *Store) GetDetails(article_id string) (model.ArticleDetails, error) {
	article_details := model.ArticleDetails{}

	query := `
		SELECT
			id,
			name,
			COALESCE(articles.description, ''),
			available_for_trade,
			(
				CASE WHEN deleted_at IS NULL
				THEN FALSE ELSE TRUE END
			) as available
		FROM articles
		WHERE id = ?;
	`

	article_row := s.db.QueryRow(query, article_id)

	if err := article_row.Scan(
		&article_details.Id,
		&article_details.Name,
		&article_details.Description,
		&article_details.AvailableForTrade,
		&article_details.IsDeleted,
	); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error scanning article details from db: %s\n", err.Error())
		return article_details, err
	}

	query = `
		SELECT
			t.id,
			t.name
		FROM articles a
		JOIN articles_tags at on a.id = at.article_id
		JOIN tags t on at.tag_id = t.id
		WHERE a.id = ?
	`

	rows, err := s.db.Query(query, article_id)

	if err != nil {
		return article_details, err
	}

	for rows.Next() {
		var tag model.ArticleDetailTag

		if err = rows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Printf("error scanning tags for article details from db: %s\n", err.Error())
			return article_details, err
		}

		article_details.Tags = append(article_details.Tags, tag)
	}

	if err := rows.Close(); err != nil {
		return article_details, err
	}

	price_row := s.db.QueryRow(`
		SELECT price
		FROM articles_prices
		WHERE article_id = ?
		ORDER BY created_at DESC;
	`, article_details.Id)

	if err := price_row.Scan(&article_details.Price); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error scanning price of article details: %s\n", err.Error())
		return article_details, err
	}

	images_rows, err := s.db.Query(`
		SELECT filename
		FROM articles_images
		WHERE article_id = ?
		ORDER BY created_at DESC;
	`, article_details.Id)

	if err != nil {
		return article_details, err
	}

	for images_rows.Next() {
		var filename string

		if err := images_rows.Scan(&filename); err != nil {
			return article_details, err
		}

		filename = uploads.GetFilePublicUrl(articles_images_bucket, filename)
		article_details.ImagesUrls = append(article_details.ImagesUrls, filename)
	}

	if err := images_rows.Close(); err != nil {
		return article_details, err
	}

	return article_details, nil
}

func (s *Store) Get(article_id string) (model.Article, error) {
	article := model.Article{}

	query := `
		SELECT
			COALESCE(tags.id, '') as tag_id,
			COALESCE(tags.name, '') as tag_name,
			articles.id,
			articles.name,
			COALESCE(articles.description, ''),
			articles.created_at,
			articles.updated_at,
			COALESCE(articles.deleted_at, ''),
			articles.available_for_trade
		FROM articles 
		LEFT JOIN articles_tags as pivot on articles.id = pivot.article_id
		LEFT JOIN tags on pivot.tag_id = tags.id
		WHERE articles.id = ?
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
			&article.AvailableForTrade,
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
		var price model.ArticlePrice

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
		var image model.ArticleImage

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

func (s *Store) Create(article *model.Article) (string, error) {
	articleID := ""

	err := s.WithTx(func(tx *sql.Tx) error {
		res, err := tx.Exec(`
			INSERT INTO articles (name, description, available_for_trade)
			VALUES (?, ?, ?)
		`, article.Name, article.Description, article.AvailableForTrade)

		if err != nil {
			log.Printf("error inserting article: %s\n", err.Error())
			return err
		}

		iArticleID, _ := res.LastInsertId()
		articleID = dbIdToString(iArticleID)
		article.Id = articleID

		if len(article.Prices) > 0 {
			if err := createPrices(tx, *article); err != nil {
				return err
			}
		}

		if err := persistArticleTags(tx, article); err != nil {
			log.Printf("could not create articles_tags of new tags for article id %s: %s\n", article.Id, err.Error())
			return err
		}

		if err := persistArticleImages(tx, article); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return articleID, nil
}

// Updates the corresponding db tables with the data inside the article
//
// If a tag that is related to the article in the db, is not present in the struct, the relationship is deleted.
//
// Tags without id are treated like new tags. They are created and then related to the article
func (s *Store) Update(article *model.Article) error {
	return s.WithTx(func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			UPDATE articles SET
				name = ?,
				description = ?,
				available_for_trade = ?,
				updated_at = current_timestamp
			WHERE id = ?;
		`, article.Name, article.Description, article.AvailableForTrade, article.Id)

		if err != nil {
			log.Printf("error updating article: %s\n", err.Error())
			return err
		}

		if len(article.Prices) > 0 {
			if err := createPrices(tx, *article); err != nil {
				return err
			}
		}

		if err := persistArticleTags(tx, article); err != nil {
			return err
		}

		if err := persistArticleImages(tx, article); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) Delete(article_id string) error {
	return s.WithTx(func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			UPDATE articles
			SET deleted_at = current_timestamp 
			WHERE id = ?
		`, article_id)

		if err != nil {
			log.Printf("Error soft-deleting article: %s\n", err.Error())
			return err
		}

		relationships_tables := []string{
			"articles_images",
			"articles_tags",
		}
		// TODO: remove images from disk

		for _, table := range relationships_tables {
			query := fmt.Sprintf("DELETE FROM %s WHERE article_id = ?;", table)

			if _, err := tx.Exec(query, article_id); err != nil {
				log.Printf("Error soft-deleting article's %s: %s\n", table, err.Error())
				return err
			}
		}

		return nil
	})
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

	var new_tags_names shared.QueryArgs

	for i, tag_name := range tags_names {
		new_tags_names = append(new_tags_names, tag_name)

		if i < len(tags_names)-1 {
			insert_tags_query.WriteString(" (?),")
		} else {
			insert_tags_query.WriteString(" (?) RETURNING id")
		}
	}

	rows, err := tx.Query(insert_tags_query.String(), new_tags_names...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags_ids := []string{}

	for rows.Next() {
		tags_ids = append(tags_ids, "")
		target := &tags_ids[len(tags_ids)-1]

		if err := rows.Scan(target); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags_ids, err
}

func createArticleTags(tx *sql.Tx, article_id string, tags_ids []string) error {
	if len(tags_ids) == 0 {
		return nil
	}

	var query_sb strings.Builder
	query_sb.WriteString("INSERT INTO articles_tags (article_id, tag_id) VALUES ")
	var queryargs shared.QueryArgs

	for i, tag_id := range tags_ids {
		queryargs = append(queryargs, article_id, tag_id)

		if i == 0 {
			query_sb.WriteString("(?, ?)")
		} else {
			query_sb.WriteString(", (?, ?)")
		}
	}

	res, err := tx.Exec(query_sb.String(), queryargs...)

	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()

	if int(count) != len(queryargs)/2 {
		return errors.New("One or more rows could not be inserted")
	}

	return nil
}

func dbIdToString(id int64) string {
	return strconv.Itoa(int(id))
}

func createPrices(tx *sql.Tx, article model.Article) error {
	var query_sb strings.Builder
	query_args := shared.QueryArgs{}

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

func persistArticleImages(tx *sql.Tx, article *model.Article) error {
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
		// TODO: delete images from disk

		return err
	}

	filenames_qargs := shared.QueryArgs{}
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
		slices.Concat(shared.QueryArgs{article.Id}, filenames_qargs)...,
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
		// TODO: delete images from disk

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

func persistArticleTags(tx *sql.Tx, article *model.Article) error {
	tags_ids_in_article := make(map[string]string)
	var new_tags_names []string

	for _, tag := range article.Tags {
		if tag.Id == "" {
			// TODO: check if names are not repeated? or deal with this from the DB and handle the error?
			new_tags_names = append(new_tags_names, tag.Name)
		} else {
			tags_ids_in_article[tag.Id] = tag.Id
		}
	}

	tags_ids_for_new_relationships, err := createTags(tx, new_tags_names)

	if err != nil {
		log.Printf("error creating new tags for article: %s\n", err.Error())
		return err
	}

	tags_ids_in_relationships_query := `
		SELECT tag_id
		FROM articles_tags
		WHERE article_id = ?;
	`

	rows, err := tx.Query(tags_ids_in_relationships_query, article.Id)

	if err != nil {
		log.Printf("could not get articles_tags: %s\n", err.Error())
		return err
	}

	defer rows.Close()

	tags_ids_in_relationships := make(map[string]string)

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			log.Printf("could not scan rows of articles_tags: %s\n", err.Error())
			return err
		}

		tags_ids_in_relationships[id] = id
	}

	if rows.Err() != nil {
		log.Printf("could not close rows of articles_tags of article: %s\n", rows.Err().Error())
		return err
	}

	for _, tag_id := range tags_ids_in_article {
		if tags_ids_in_relationships[tag_id] == "" {
			tags_ids_for_new_relationships = append(tags_ids_for_new_relationships, tag_id)
		}
	}

	log.Printf("\ntags_ids_in_relationships: %v \n", tags_ids_in_relationships)

	for _, tag_id := range tags_ids_in_relationships {
		log.Printf("tags_ids_in_article[tag_id] where tag_id = %s: %v \n", tag_id, tags_ids_in_article[tag_id])

		if tags_ids_in_article[tag_id] == "" {
			log.Println("executing delete on articles_tags")

			_, err := tx.Exec("DELETE FROM articles_tags WHERE tag_id = ? AND article_id = ?;", tag_id, article.Id)

			if err != nil {
				log.Printf("could not delete articles_tags for tag id: %s\n", tag_id)
				return err
			}
		}
	}

	err = createArticleTags(tx, article.Id, tags_ids_for_new_relationships)

	if err != nil {
		log.Printf("could not create articles_tags of new tags for article id %s: %s\n", article.Id, err.Error())
		return err
	}

	return nil
}
