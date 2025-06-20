package clientiface

import "github.com/sottey/ainvil/lib/model"

type APIClient interface {
	Name() string
	GetEntries(startDate, endDate string) ([]model.Entry, error)
	GetAllEntries() ([]model.Entry, error)
}
