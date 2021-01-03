package models

// Row - товар из файла с данными в формате .xlsx
type Row struct {
	OfferID   int
	Name      string
	Price     float64
	Quantity  int
	Available bool
}

// RowString ...
type RowString struct {
	OfferID   string
	Name      string
	Price     string
	Quantity  string
	Available string
}
