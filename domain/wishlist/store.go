package wishlist

import (
	"database/sql"
	"log"

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
		return w, err
	}

	// TODO: populate images
	// TODO: populate tags

	return w, nil
}

func (s *Store) GetWishlist() ([]Wishitem, error) {
	query := `
		SELECT w.id, w.name
		FROM wishitems w
		WHERE deleted_at IS NULL
		ORDER BY created_at
	`

	rows, err := s.db.Query(query)

	if err != nil {
		log.Printf("error querying wishlists: %s\n", err.Error())
		return nil, err
	}

	defer rows.Close()

	var wishlist []Wishitem

	for rows.Next() {
		var w Wishitem

		if err := rows.Scan(&w.Id, &w.Name); err != nil {
			log.Printf("error scanning wishlist row: %s\n", err.Error())
			return nil, err
		}

		wishlist = append(wishlist, w)
	}

	return wishlist, nil
}

func (s *Store) GetAdminList() ([]Wishitem, error) {
	// TODO:
	return s.GetWishlist()
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
