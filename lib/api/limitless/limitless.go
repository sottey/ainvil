package limitless

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sottey/ainvil/lib/clientiface"
	"github.com/sottey/ainvil/lib/model"
)

const baseURL = "https://api.limitless.ai/v1/lifelogs"

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
	var all []model.Entry
	cursor := ""

	for {
		reqURL, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		q := reqURL.Query()
		q.Set("start", startDate)
		q.Set("end", endDate)
		q.Set("includeMarkdown", "true")
		q.Set("includeHeadings", "true")
		q.Set("limit", "25")
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		reqURL.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", reqURL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-API-Key", c.token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("limitless API error (%d): %s", resp.StatusCode, string(body))
		}

		var parsed struct {
			Data struct {
				Lifelogs []struct {
					ID        string `json:"id"`
					Markdown  string `json:"markdown"`
					StartTime string `json:"startTime"`
					Title     string `json:"title"`
				} `json:"lifelogs"`
			} `json:"data"`
			Meta struct {
				Lifelogs struct {
					NextCursor string `json:"nextCursor"`
				} `json:"lifelogs"`
			} `json:"meta"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
			return nil, err
		}

		for _, l := range parsed.Data.Lifelogs {
			t, _ := time.Parse(time.RFC3339, l.StartTime)
			all = append(all, model.Entry{
				ID:        l.ID,
				Source:    "limitless",
				Timestamp: t,
				Content:   l.Markdown,
				Metadata: map[string]string{
					"title": l.Title,
				},
			})
		}

		cursor = parsed.Meta.Lifelogs.NextCursor
		if cursor == "" {
			break
		}
	}

	return all, nil
}

func (c *Client) GetAllEntries() ([]model.Entry, error) {
	// Pass a very wide date range — adjust as needed
	return c.GetEntries("2000-01-01", time.Now().Format("2006-01-02"))
}
