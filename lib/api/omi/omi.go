package omi

import (
	"fmt"
	"time"

	"github.com/sottey/ainvil/lib/clientiface"
	"github.com/sottey/ainvil/lib/model"
)

type Client struct {
	token string
}

func NewClient(token string) clientiface.APIClient {
	return &Client{token: token}
}

func (c *Client) Name() string {
	return "omi"
}

func (c *Client) GetEntries(startDate, endDate string) ([]model.Entry, error) {
	// TODO: Replace this stub with real Omi API calls
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var entries []model.Entry
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		entries = append(entries, model.Entry{
			ID:        fmt.Sprintf("omi-%s", d.Format("20060102")),
			Source:    "omi",
			Timestamp: d,
			Content:   fmt.Sprintf("Simulated Omi log for %s", d.Format("Jan 2 2006")),
		})
	}

	return entries, nil
}

func (c *Client) GetAllEntries() ([]model.Entry, error) {
	// TODO: Make paginated API calls to retrieve all logs
	return c.GetEntries("2025-06-01", "2025-06-03")
}
