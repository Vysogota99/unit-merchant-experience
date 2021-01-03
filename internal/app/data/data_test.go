package data

import (
	"testing"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestReadXLSX(t *testing.T) {
	fileName := "../../../build/nginx/files/1.xlsx"
	data, err := ReadXLSX(fileName)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestValidateDataFromXLSX(t *testing.T) {
	rows := []models.RowString{
		models.RowString{
			OfferID:   "1",
			Name:      "iphone4",
			Price:     "10000",
			Quantity:  "110",
			Available: "true",
		},
		models.RowString{
			OfferID:   "2",
			Name:      "iphone5",
			Price:     "12000",
			Quantity:  "10",
			Available: "true",
		},
	}

	res, err := Validate(rows)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
