package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetOffersIDSBySalerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	store := New(db)

	uID := 1

	rows := mock.NewRows(
		[]string{"id", "saler_id", "name", "price", "quantity"},
	).AddRow("10", "1", "iphone_X", "40000", "10").
		AddRow("11", "1", "iphone_XR", "42000", "0").
		AddRow("12", "1", "iphone_11", "51000", "90")

	query := `
		SELECT id 
		FROM offers
		WHERE saler_id = $
		ORDER BY saler_id asc 
	`

	mock.ExpectBegin()
	stmt := mock.ExpectPrepare(query)
	stmt.ExpectQuery().WithArgs(uID).WillReturnRows(rows)
	mock.ExpectCommit()

	res, err := store.Offer().GetOffersIDSBySalerID(uID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
