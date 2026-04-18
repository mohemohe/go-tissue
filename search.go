package go_tissue

import (
	"context"
	"net/url"
	"strconv"
)

type SearchCheckinsOption struct {
	Query   string
	Page    int
	PerPage int
}

func (c *Client) SearchCheckins(ctx context.Context, option *SearchCheckinsOption) ([]Checkin, error) {
	query := url.Values{}
	if option != nil {
		if option.Query != "" {
			query.Set("q", option.Query)
		}
		if option.Page > 0 {
			query.Set("page", strconv.Itoa(option.Page))
		}
		if option.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(option.PerPage))
		}
	}
	result := []Checkin{}
	if err := c.getJSON(ctx, "/api/search/checkins", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}
