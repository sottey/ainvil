package limitless

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
	return "limitless"
}

func (c *Client) GetEntries(startDate, endDate string) ([]model.Entry, error) {
	// TODO: Replace with real API call using c.token
	// Simulate fake data for now
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var entries []model.Entry
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		entries = append(entries, model.Entry{
			ID:        fmt.Sprintf("limitless-%s", d.Format("20060102")),
			Source:    "limitless",
			Timestamp: d,
			Content:   fmt.Sprintf("Simulated Limitless log for %s", d.Format("Jan 2 2006")),
		})
	}

	return entries, nil
}

func (c *Client) GetAllEntries() ([]model.Entry, error) {
	// TODO: Pull all entries (paginate if needed)
	// For now, simulate 3 days of logs
	return c.GetEntries("2025-06-01", "2025-06-03")
}
