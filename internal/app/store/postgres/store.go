package postgres

import (
	"database/sql"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
	_ "github.com/lib/pq"
)

// StorePSQL ...
type StorePSQL struct {
	offerRepository *offerRepository
	db              *sql.DB
}

// New - инициализирует Store
func New(db *sql.DB) *StorePSQL {
	return &StorePSQL{
		db: db,
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
