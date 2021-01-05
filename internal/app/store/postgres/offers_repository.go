package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"

	_ "github.com/lib/pq"
)

// OfferRepository ...
type offerRepository struct {
	store *StorePSQL
}

// GetOffersIDBySalerID - возвращает массив id товаров определенного продавца
func (o *offerRepository) GetOffersIDSBySalerID(id int) ([]int, error) {
	db, err := sql.Open("postgres", o.store.connectionString)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := `
		SELECT id 
		FROM offers
		WHERE saler_id = $1
		ORDER BY id asc 
	`

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer stmt.Close()

	result := make([]int, 0)
	rows, err := stmt.QueryContext(ctx, id)
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

// InsertOffers ...
func (o *offerRepository) InsertOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error) {
	if len(rows) == 0 {
		return 0, nil
	}

	query := `
		INSERT INTO offers(id, saler_id, name, price, quantity)
		VALUES 
	`
	for _, row := range rows {
		values := fmt.Sprintf("(%d, %d, '%s', %f, %d),", row.OfferID, salerID, row.Name, row.Price, row.Quantity)
		query = fmt.Sprintf("%s %s", query, values)
	}

	query = query[:len(query)-1]

	res, err := tx.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

// UpdateOffers ...
func (o *offerRepository) UpdateOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error) {
	db, err := sql.Open("postgres", o.store.connectionString)
	if err != nil {
		return 0, err
	}

	defer db.Close()

	if len(rows) == 0 {
		return 0, nil
	}

	counter := 0
	for _, row := range rows {
		query := `
			UPDATE offers
			SET name=$1, price=$2, quantity=$3
			WHERE id=$4
		`

		stmt, err := db.Prepare(query)
		if err != nil {
			return 0, err
		}

		res, err := stmt.ExecContext(ctx, row.Name, row.Price, row.Quantity, row.OfferID)
		if err != nil {
			return 0, nil
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, nil
		}

		counter += int(rowsAffected)
	}

	return counter, nil
}

// DeleteOffers ...
func (o *offerRepository) DeleteOffers(ctx context.Context, tx *sql.Tx, ids []int) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	query := `
		DELETE FROM offers 
		WHERE id IN (%s)
	`
	query = fmt.Sprintf(query, arrayToString(ids, ","))

	res, err := tx.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (o *offerRepository) WorkerPipeline(rowsToInsert []models.Row, rowsToUpdate []models.Row, idsToDelete []int, salerID int) (*models.WorkerResult, error) {
	db, err := sql.Open("postgres", o.store.connectionString)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	nInserted, err := o.store.offerRepository.InsertOffers(ctx, tx, rowsToInsert, salerID)
	if err != nil {
		return nil, err
	}

	nDeleted, err := o.store.offerRepository.DeleteOffers(ctx, tx, idsToDelete)
	if err != nil {
		return nil, err
	}

	nUpdated, err := o.store.offerRepository.UpdateOffers(ctx, tx, rowsToUpdate, salerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.WorkerResult{
		NInserted: nInserted,
		NDeleted:  nDeleted,
		NUpdated:  nUpdated,
	}, nil
}

func (o *offerRepository) GetOffers(ctx context.Context, offerID, salerID, offer string) ([]models.Row, error) {
	db, err := sql.Open("postgres", o.store.connectionString)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, saler_id, name, price, quantity
		FROM offers %s
	`

	offer = fmt.Sprintf("'%%%s%%'", offer)
	fields := map[string][]string{
		"id":       []string{"=", offerID},
		"saler_id": []string{"=", salerID},
		"name":     []string{"LIKE", offer},
	}

	condition := "WHERE %s"
	predicats := make([]string, 0)

	for key, value := range fields {
		if value[1] != "" {
			predicats = append(predicats, fmt.Sprintf("%s %s %s", key, value[0], value[1]))
		}
	}

	if cap(predicats) == 1 {
		query = fmt.Sprintf(query, "")
	} else {
		condition = fmt.Sprintf(condition, strings.Join(predicats, " AND "))
		query = fmt.Sprintf(query, condition)
	}

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]models.Row, 0)
	for rows.Next() {
		row := models.Row{}
		if err := rows.Scan(&row.OfferID, &row.SalerID, &row.Name, &row.Price, &row.Quantity); err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
