package mock

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
)

// StoreMock ...
type StoreMock struct {
	db              *sql.DB
	mock            sqlmock.Sqlmock
	offerRepository *offerRepository
}

// New - инициализирует Store
func New(db *sql.DB, mock sqlmock.Sqlmock) *StoreMock {
	return &StoreMock{
		db:   db,
		mock: mock,
	}
}

// Offer ...
func (s *StoreMock) Offer() store.OfferRepository {
	if s.offerRepository == nil {
		s.offerRepository = &offerRepository{
			store: s,
		}
	}

	return s.offerRepository
}
