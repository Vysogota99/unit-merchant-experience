package data

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"
	"github.com/tealeg/xlsx/v3"
)

// DownloadFile - скачивает файл по ссылке и сохраняет в дирректорию /static
func DownloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ReadXLSX - считывает данные из файла.xlsx
func ReadXLSX(fileName string) ([]models.RowString, error) {
	wb, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	result := make([]models.RowString, 0)
	for _, sh := range wb.Sheets {
		nRows := sh.MaxRow
		rowsInSheet := make([]models.RowString, nRows-1)
		for i := 1; i < nRows; i++ {
			val, err := sh.Cell(i, 0)
			if err != nil {
				return nil, err
			}
			rowsInSheet[i-1].OfferID = val.Value

			val, err = sh.Cell(i, 1)
			if err != nil {
				return nil, err
			}
			rowsInSheet[i-1].Name = val.Value

			val, err = sh.Cell(i, 2)
			if err != nil {
				return nil, err
			}
			rowsInSheet[i-1].Price = val.Value

			val, err = sh.Cell(i, 3)
			if err != nil {
				return nil, err
			}

			rowsInSheet[i-1].Quantity = val.Value

			val, err = sh.Cell(i, 4)
			if err != nil {
				return nil, err
			}

			rowsInSheet[i-1].Available = val.Value
		}

		result = append(result, rowsInSheet...)
	}

	return result, err
}

// Validate - приводит данные к нужным типам и проводит валидацию
func Validate(dataToValidate []models.RowString) ([]models.Row, error) {
	result := make([]models.Row, len(dataToValidate))
	for i, row := range dataToValidate {
		data := models.Row{}
		value, err := strconv.ParseInt(row.OfferID, 10, 64)
		if err != nil {
			return nil, err
		}

		if value <= 0 {
			return nil, fmt.Errorf("Индекс товара должен быть положительным")
		}

		data.OfferID = int(value)
		data.Name = row.Name

		valueFloat, err := strconv.ParseFloat(row.Price, 64)
		if err != nil {
			return nil, err
		}

		if valueFloat <= 0 {
			return nil, fmt.Errorf("цена должна быть положительной")
		}

		data.Price = valueFloat

		value, err = strconv.ParseInt(row.Quantity, 10, 64)
		if err != nil {
			return nil, err
		}

		if value < 0 {
			return nil, fmt.Errorf("Количество товаров не может быть отрицательным")
		}

		data.Quantity = int(value)

		valueBool, err := strconv.ParseBool(row.Available)
		if err != nil {
			return nil, err
		}
		data.Available = valueBool

		result[i] = data
	}
	return result, nil
}
