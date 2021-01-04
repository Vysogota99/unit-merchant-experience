package mock

// OfferRepository ...
type offerRepository struct {
	store *StoreMock
}

// GetOffersIDBySalerID - возвращает массив id товаров определенного продавца
func (o *offerRepository) GetOffersIDSBySalerID(id int) ([]int, error) {
	query := `
		SELECT id 
		FROM offers
		WHERE saler_id = $1
		ORDER BY saler_id asc 
	`

	result := make([]int, 0)
	rows, err := o.store.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var lotID int
		if err := rows.Scan(&lotID); err != nil {
			return nil, err
		}

		result = append(result, lotID)
	}

	return result, nil
}
