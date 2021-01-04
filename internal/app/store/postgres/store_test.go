package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestGetOffersIDSBySalerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	store := New(db)

	uID := 1

	rows := mock.NewRows(
		[]string{"id"},
	).AddRow("1").AddRow("2").AddRow("3").AddRow("4")

	query := `SELECT id FROM offers WHERE saler_id = $1 ORDER BY saler_id asc`

	mock.ExpectBegin()
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query))
	stmt.ExpectQuery().WithArgs(uID).WillReturnRows(rows)
	mock.ExpectCommit()

	res, err := store.Offer().GetOffersIDSBySalerID(uID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestWorkerPipeline(t *testing.T) {
	db, err := sql.Open("postgres", "user=user1 password=password dbname=app sslmode=disable")
	assert.NoError(t, err)
	store := New(db)

	rowsIns := []models.Row{
		models.Row{
			OfferID:   1,
			Name:      "iphone4",
			Price:     10000,
			Quantity:  110,
			Available: true,
		},
		models.Row{
			OfferID:   2,
			Name:      "iphone4",
			Price:     10000,
			Quantity:  110,
			Available: true,
		},
	}

	rowsUpd := []models.Row{
		models.Row{
			OfferID:   1,
			Name:      "iphone4",
			Price:     10000,
			Quantity:  110,
			Available: true,
		},
		models.Row{
			OfferID:   2,
			Name:      "iphone4",
			Price:     10000,
			Quantity:  110,
			Available: true,
		},
	}

	idsDel := []int{}

	_, err = store.Offer().WorkerPipeline(rowsIns, rowsUpd, idsDel, 1)
	assert.NoError(t, err)
}

func TestGetOffer(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	store := New(db)

	offerID := ""
	salerID := "1"
	offer := "ip"

	rows := mock.NewRows(
		[]string{"id", "saler_id", "name", "price", "quantity"},
	).AddRow("1", "1", "iphone_X", "40000", "10").
		AddRow("2", "1", "iphone_XR", "48000", "5")

	query := `
		SELECT id, saler_id, name, price, quantity
		FROM offers
		WHERE saler_id = 1 AND name LIKE '%ip%'
	`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	mock.ExpectCommit()

	res, err := store.Offer().GetOffers(context.Background(), offerID, salerID, offer)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}
