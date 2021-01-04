package data

import (
	"io"
	"net/http"
	"os"
	"sort"

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

// Validate - приводит данные к нужным типам и проводит валидацию. Вовзращает массив данных для
//				добавления, для обновления и количество строк с ошибками.
func Validate(dataToValidate []models.RowString, ids []int) ([]models.Row, []models.Row, []int, int) {
	dataToInsert := make([]models.Row, 0)
	idsToDelete := make([]int, 0)
	dataToUpdate := make([]models.Row, 0)
	dataWithErrors := 0

	for _, row := range dataToValidate {
		data := models.Row{}

		value, err := strconv.ParseInt(row.OfferID, 10, 64)
		if err != nil {
			dataWithErrors++
			continue
		}

		if value <= 0 {
			dataWithErrors++
			continue
		}

		data.OfferID = int(value)
		data.Name = row.Name

		valueFloat, err := strconv.ParseFloat(row.Price, 64)
		if err != nil {
			dataWithErrors++
			continue
		}

		if valueFloat <= 0 {
			dataWithErrors++
			continue
		}

		data.Price = valueFloat

		value, err = strconv.ParseInt(row.Quantity, 10, 64)
		if err != nil {
			dataWithErrors++
			continue
		}

		if value < 0 {
			dataWithErrors++
			continue
		}

		data.Quantity = int(value)

		valueBool, err := strconv.ParseBool(row.Available)
		if err != nil {
			dataWithErrors++
			continue
		}

		data.Available = valueBool

		// проверка на необходимость обновления
		// если в списке с id товаров только один элемент - проверяем соответствие,
		// иначе пользуемся поиском offer_id в этом списке

		if len(ids) == 1 {
			if ids[0] == data.OfferID {

				// если строка подлежит удалению, сразу добавляем ее id
				// в список для удаления
				if data.Available == true {
					dataToUpdate = append(dataToUpdate, data)
				} else {
					idsToDelete = append(idsToDelete, data.OfferID)
				}

				continue
			}
		} else {
			if inSlice := sort.SearchInts(ids, data.OfferID); inSlice != len(ids) {

				// если строка подлежит удалению, сразу добавляем ее id
				// в список для удаления
				if data.Available == true {
					dataToUpdate = append(dataToUpdate, data)
				} else {
					idsToDelete = append(idsToDelete, data.OfferID)
				}

				continue
			}
		}

		dataToInsert = append(dataToInsert, data)

	}

	return dataToInsert, dataToUpdate, idsToDelete, dataWithErrors
}
