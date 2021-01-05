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
	type testCase struct {
		name   string
		rows   []models.RowString
		ids    []int
		insExp int
		updExp int
		delExp int
		errExp int
	}

	tCases := []testCase{
		testCase{
			name:   "test 1",
			insExp: 3,
			updExp: 0,
			delExp: 1,
			errExp: 1,
			rows: []models.RowString{
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
			},
			ids: []int{2},
		},
		testCase{
			name:   "test 2",
			insExp: 2,
			updExp: 1,
			delExp: 1,
			errExp: 1,
			rows: []models.RowString{
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
			},
			ids: []int{1, 2},
		},
		testCase{
			name:   "test 2",
			insExp: 4,
			updExp: 1,
			delExp: 0,
			errExp: 1,
			rows: []models.RowString{
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
			},
			ids: []int{},
		},
		testCase{
			name:   "test 2",
			insExp: 0,
			updExp: 4,
			delExp: 1,
			errExp: 0,
			rows: []models.RowString{
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
					OfferID:   "5",
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
			},
			ids: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			ins, upd, del, err := Validate(tc.rows, tc.ids)
			assert.Len(t, ins, tc.insExp)
			assert.Len(t, upd, tc.updExp)
			assert.Len(t, del, tc.delExp)
			assert.Equal(t, err, tc.errExp)
		})
	}
}
