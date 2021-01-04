package postgres

import (
	"context"
	"log"
)

// OfferRepository ...
type offerRepository struct {
	store *StorePSQL
}

// GetOffersIDBySalerID - возвращает массив id товаров определенного продавца
func (o *offerRepository) GetOffersIDSBySalerID(id int) ([]int, error) {
	query := `
		SELECT id 
		FROM offers
		WHERE saler_id = $1
		ORDER BY saler_id asc 
	`

	ctx := context.Background()
	tx, err := o.store.db.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer stmt.Close()

	result := make([]int, 0)
	rows, err := stmt.QueryContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var lotID int
		if err := rows.Scan(&lotID); err != nil {
			return nil, err
		}

		result = append(result, lotID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}
