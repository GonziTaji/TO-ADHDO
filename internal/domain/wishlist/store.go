package wishlist

import (
	"database/sql"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/yogusita/to-adhdo/domain/shared"
)

type Store struct {
	db *sql.DB
}

func CreateStore(db *sql.DB) Store {
	return Store{db}
}

func (s *Store) GetWishitem(item_id string) (Wishitem, error) {
	query := `
		SELECT
			w.id,
			w.name,
			w.description,
			w.observed_price,
			w.external_url
		FROM wishitems w
		WHERE w.id = ?;
	`

	row := s.db.QueryRow(query, item_id)

	var w Wishitem

	if err := row.Scan(
		&w.Id,
		&w.Name,
		&w.Description,
		&w.ObservedPrice,
		&w.ExternalUrl,
	); err != nil {
		log.Printf("error querying wishlists: %s\n", err.Error())
		return Wishitem{}, err
	}

	query = `
		SELECT
			wt.id,
			t.id,
			t.name
		FROM wishitems_tags wt
		JOIN tags t
			ON t.id = wt.tag_id
		WHERE wt.wishitem_id = ?
		ORDER BY t.name;
	`

	rows, err := s.db.Query(query, item_id)

	if err != nil {
		log.Printf("error querying wishitem tags: %s\n", err.Error())
		return Wishitem{}, err
	}

	for rows.Next() {
		var tag WishitemTag
		if err := rows.Scan(
			&tag.Id,
			&tag.TagId,
			&tag.TagName,
		); err != nil {
			log.Println("error scanning wishitem tag")
			rows.Close()
			return Wishitem{}, err
		}

		w.Tags = append(w.Tags, tag)
	}

	query = `
		SELECT
			id,
			filename
		FROM wishitems_images
		WHERE wishitem_id = ?
		ORDER BY created_at;
	`

	rows, err = s.db.Query(query, item_id)

	if err != nil {
		log.Println("error querying wishitem images")
		return Wishitem{}, nil
	}

	for rows.Next() {
		var img WishitemImage
		if err := rows.Scan(
			&img.Id,
			&img.Filepath,
		); err != nil {
			log.Println("error scanning wishitem image")
			rows.Close()
			return Wishitem{}, err
		}

		w.Images = append(w.Images, img)
	}

	return w, nil
}

func (s *Store) GetWishlist(options WishlistFilterParams) (WishlistData, error) {
	var data WishlistData

	var wishitems_qa shared.QueryArgs

	tags_count := len(options.TagsIds)

	var where_clauses []string

	if options.PriceRangeStart > 0 {
		where_clauses = append(where_clauses, "wi.observed_price >= ?")
		wishitems_qa = append(wishitems_qa, options.PriceRangeStart)
		data.PriceSelectedRange.Start = options.PriceRangeStart
	}

	if options.PriceRangeEnd > 0 {
		where_clauses = append(where_clauses, "wi.observed_price <= ?")
		wishitems_qa = append(wishitems_qa, options.PriceRangeEnd)
		data.PriceSelectedRange.End = options.PriceRangeEnd
	}

	if options.SearchTerm != "" {
		where_clauses = append(where_clauses, "wi.name LIKE ?")
		wishitems_qa = append(wishitems_qa, "%"+options.SearchTerm+"%")
		data.SearchTerm = options.SearchTerm
	}

	tags_options_q := fmt.Sprintf(`
		SELECT
			t.id,
			t.name,
			count(wi.id)
		FROM tags t
		LEFT JOIN wishitems_tags wt
			ON wt.tag_id = t.id
		LEFT JOIN wishitems wi
			ON wi.id = wt.wishitem_id
			%s
		GROUP BY t.id, t.name;
	`, mergeExtraJoinWhereClauses(where_clauses))

	// log.Printf("tags query: %s\n", tags_options_q)
	// log.Printf("query args: %v\n", wishitems_qa)

	rows, err := s.db.Query(tags_options_q, wishitems_qa...)

	if err != nil {
		log.Println("error getting tags options for wishlist")
		return WishlistData{}, err
	}

	for rows.Next() {
		var tag_option TagSelectOption

		if err := rows.Scan(
			&tag_option.Id,
			&tag_option.Name,
			&tag_option.Count,
		); err != nil {
			rows.Close()
			log.Println("error scanning tag option for wishlist")
			return WishlistData{}, err
		}

		if slices.Contains(options.TagsIds, tag_option.Id) {
			tag_option.Selected = true
		}

		data.TagsSelectOptions = append(data.TagsSelectOptions, tag_option)
	}

	if tags_count > 0 {
		ids_templates := strings.Join(slices.Repeat([]string{"?"}, tags_count), ",")
		clause := fmt.Sprintf("tag.id IN (%s)", ids_templates)
		where_clauses = append(where_clauses, clause)

		for _, id := range options.TagsIds {
			wishitems_qa = append(wishitems_qa, id)
		}
	}

	sortColumn := "id"
	switch options.SortBy {
	case WishlistSortByPrice:
		sortColumn = "observed_price"

	case WishlistSortByCratedAt:
		sortColumn = "created_at"
	}

	// SQL injection guard - we make sure the value for the query is either DESC or ASC
	sortDirection := "DESC"
	if options.SortDirection == SortDirectionAsc {
		sortColumn = "ASC"
	}

	wishitems_q := fmt.Sprintf(`
		SELECT
			wi.id,
			wi.name,
			wi.observed_price,
			wt.id,
			tag.id,
			tag.name
		FROM tags tag
		JOIN wishitems_tags wt ON wt.tag_id = tag.id
		JOIN wishitems wi ON wi.id = wt.wishitem_id
		%s
		ORDER BY wi.%s %s;
	`, mergeExtraJoinWhereClauses(where_clauses), sortColumn, sortDirection)

	// log.Printf("wishitems query:\n%s\n", wishitems_q)
	// log.Printf("wishitems args:\n%v\n", wishitems_qa)

	rows, err = s.db.Query(wishitems_q, wishitems_qa...)

	if err != nil {
		log.Printf("error querying wishitems: %s\n", err.Error())
		return WishlistData{}, err
	}

	wishitems_by_id := make(map[string]Wishitem)

	for rows.Next() {
		var scanned Wishitem
		var tag WishitemTag

		if err := rows.Scan(
			&scanned.Id,
			&scanned.Name,
			&scanned.ObservedPrice,
			&tag.Id,
			&tag.TagId,
			&tag.TagName,
		); err != nil {
			log.Printf("error scanning wishlist row: %s\n", err.Error())
			rows.Close()
			return WishlistData{}, err
		}

		wi := wishitems_by_id[scanned.Id]
		wi.Id = scanned.Id
		wi.Name = scanned.Name
		wi.ObservedPrice = scanned.ObservedPrice
		wi.Tags = append(wi.Tags, tag)

		wishitems_by_id[scanned.Id] = wi

		if data.PriceRange.End < scanned.ObservedPrice {
			data.PriceRange.End = scanned.ObservedPrice
		} else if data.PriceRange.Start > scanned.ObservedPrice {
			data.PriceRange.Start = scanned.ObservedPrice
		}
	}

	data.Items = slices.Collect(maps.Values(wishitems_by_id))

	return data, nil
}

func (s *Store) GetAdminList(options WishlistFilterParams) (WishlistData, error) {
	// TODO:
	return s.GetWishlist(options)
}

func (s *Store) SaveWishitem(formdata WishitemFormData) (string, error) {
	var q string
	var qargs shared.QueryArgs

	tx, err := s.db.Begin()

	if err != nil {
		tx.Rollback()
		return "", err
	}

	if len(formdata.Id) > 0 {
		q = `
			UPDATE wishitems SET
				name = ?,
				description = ?,
				external_url = ?,
				observed_price = ?,
				updated_at = current_timestamp
			WHERE id = ?
			RETURNING id;
		`
		qargs = shared.QueryArgs{
			formdata.Name,
			formdata.Description,
			formdata.ExternalUrl,
			formdata.ObservedPrice,
			formdata.Id,
		}
	} else {
		q = `
			INSERT INTO wishitems (name, description, external_url, observed_price)
			VALUES (?, ?, ?, ?)
			RETURNING id;
		`
		qargs = shared.QueryArgs{
			formdata.Name,
			formdata.Description,
			formdata.ExternalUrl,
			formdata.ObservedPrice,
		}
	}

	row := tx.QueryRow(q, qargs...)

	if err := row.Scan(&formdata.Id); err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", err
	}

	// TODO: insert tags
	// TODO: insert images

	return formdata.Id, nil
}

func (s *Store) DeleteWishitem(id string) error {
	q := "UPDATE wishitems SET deleted_at = current_timestamp WHERE id = ?"

	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}

	return nil
}

// given a list of strings S, it's converted into "AND S[0] AND S[1] AND S[2]..."
func mergeExtraJoinWhereClauses(where_clauses []string) string {
	var sb strings.Builder

	if len(where_clauses) > 0 {
		fmt.Fprintf(&sb, " AND %s ", where_clauses[0])

		if len(where_clauses) > 1 {
			sb.WriteString(" AND ")
			sb.WriteString(strings.Join(where_clauses[1:], " AND "))
		}
	}

	return sb.String()
}
