package postgres

import (
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
)

// StorePSQL ...
type StorePSQL struct {
	offerRepository  *offerRepository
	connectionString string
}

// New - инициализирует Store
func New(conString string) *StorePSQL {
	return &StorePSQL{
		connectionString: conString,
	}
}

// Offer ...
func (s *StorePSQL) Offer() store.OfferRepository {
	if s.offerRepository == nil {
		s.offerRepository = &offerRepository{
			store: s,
		}
	}

	return s.offerRepository
}
