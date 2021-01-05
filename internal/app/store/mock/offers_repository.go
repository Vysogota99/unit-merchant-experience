package mock

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store/postgres"
)

// OfferRepository ...
type offerRepository struct {
	store *StoreMock
}

// GetOffersIDBySalerID - возвращает массив id товаров определенного продавца
func (o *offerRepository) GetOffersIDSBySalerID(id int) ([]int, error) {
	rows := o.store.mock.NewRows(
		[]string{"id"},
	).AddRow("1").AddRow("2").AddRow("3").AddRow("4")

	query := `SELECT id FROM offers WHERE saler_id = $1 ORDER BY saler_id asc`

	o.store.mock.ExpectBegin()
	stmt := o.store.mock.ExpectPrepare(regexp.QuoteMeta(query))
	stmt.ExpectQuery().WithArgs(id).WillReturnRows(rows)
	o.store.mock.ExpectCommit()

	storePostgres := postgres.New(o.store.db)

	return storePostgres.Offer().GetOffersIDSBySalerID(id)
}

// InsertOffers ...
func (o *offerRepository) InsertOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error) {
	return 0, nil
}

// UpdateOffers ...
func (o *offerRepository) UpdateOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error) {
	return 0, nil
}

// DeleteOffers ...
func (o *offerRepository) DeleteOffers(ctx context.Context, tx *sql.Tx, ids []int) (int, error) {
	return 0, nil
}

func (o *offerRepository) WorkerPipeline(rowsToInsert []models.Row, rowsToUpdate []models.Row, idsToDelete []int, salerID int) (*models.WorkerResult, error) {
	return nil, nil
}

func (o *offerRepository) GetOffers(ctx context.Context, offerID, salerID, offer string) ([]models.Row, error) {
	return []models.Row{
		models.Row{
			OfferID: 1,
			SalerID: 2,
			Name: "name",
			Price: 123,
			Quantity: 1,
		},
	}, nil
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
