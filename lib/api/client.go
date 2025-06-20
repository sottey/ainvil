package api

import (
	"github.com/sottey/ainvil/lib/model"
)

// APIClient defines the interface for supported wearable APIs.
type APIClient interface {
	// Name returns the unique identifier of the API (e.g., "omi", "limitless")
	Name() string

	// GetEntries fetches entries between the given date range (inclusive).
	// Dates are in the format "YYYY-MM-DD".
	GetEntries(startDate, endDate string) ([]model.Entry, error)

	// GetAllEntries fetches all entries available from the API.
	GetAllEntries() ([]model.Entry, error)
}
