package models

// Row - товар из файла с данными в формате .xlsx
type Row struct {
	OfferID   int
	SalerID   int
	Name      string
	Price     float64
	Quantity  int
	Available bool `json:"available,omitempty"`
}

// RowString ...
type RowString struct {
	OfferID   string
	Name      string
	Price     string
	Quantity  string
	Available string
}

// WorkerResult - краткая статистика, результат работы воркера
type WorkerResult struct {
	NInserted   int `json:"количество созданных строк"`
	NUpdated    int `json:"количество обновленных строк"`
	NDeleted    int `json:"количество удаленных строк"`
	NWithErrors int `json:"количество строк с ошибками"`
}
