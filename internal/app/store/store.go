package store

// Store - ...
type Store interface {
	Offer() OfferRepository
}

// OfferRepository ...
type OfferRepository interface {
	GetOffersIDSBySalerID(id int) ([]int, error)
}
