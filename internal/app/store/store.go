package store

import (
	"context"
	"database/sql"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/models"
)

// Store - ...
type Store interface {
	Offer() OfferRepository
}

// OfferRepository ...
type OfferRepository interface {
	GetOffersIDSBySalerID(id int) ([]int, error)
	InsertOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error)
	UpdateOffers(ctx context.Context, tx *sql.Tx, rows []models.Row, salerID int) (int, error)
	DeleteOffers(ctx context.Context, tx *sql.Tx, ids []int) (int, error)
	WorkerPipeline(rowsToInsert []models.Row, rowsToUpdate []models.Row, idsToDelete []int, salerID int) (*models.WorkerResult, error)
	GetOffers(ctx context.Context, offerID, salerID, offer string) ([]models.Row, error)
}
