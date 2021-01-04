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

func TestDownloadFile(t *testing.T) {
	filePath := "../../../static/test.xlsx"
	url := "http://127.0.0.1:80/files/1.xlsx"
	err := DownloadFile(filePath, url)
	assert.NoError(t, err)
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
			Available: "false",
		},
		models.RowString{
			OfferID:   "-1",
			Name:      "iphone5",
			Price:     "12000",
			Quantity:  "10",
			Available: "true",
		},
		models.RowString{
			OfferID:   "3",
			Name:      "iphone5",
			Price:     "12000",
			Quantity:  "10",
			Available: "true",
		},
		models.RowString{
			OfferID:   "4",
			Name:      "iphone5",
			Price:     "12000",
			Quantity:  "10",
			Available: "true",
		},
	}

	ins, upd, del, err := Validate(rows, []int{2})
	assert.Len(t, ins, 2)
	assert.Len(t, upd, 1)
	assert.Len(t, del, 1)
	assert.NotZero(t, err)
}
