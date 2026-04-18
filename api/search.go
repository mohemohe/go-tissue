package api

import (
	"context"
	"net/url"
	"strconv"

	tissue "github.com/mohemohe/go-tissue"
)

type SearchOption struct {
	Query   string
	Page    int
	PerPage int
}

func (c *Client) SearchCheckins(ctx context.Context, option *SearchOption) ([]tissue.Checkin, error) {
	result := []tissue.Checkin{}
	if err := c.getJSON(ctx, "/v1/search/checkins", buildSearchQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) SearchCollections(ctx context.Context, option *SearchOption) ([]tissue.CollectionItem, error) {
	result := []tissue.CollectionItem{}
	if err := c.getJSON(ctx, "/v1/search/collections", buildSearchQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func buildSearchQuery(option *SearchOption) url.Values {
	query := url.Values{}
	if option == nil {
		return query
	}
	if option.Query != "" {
		query.Set("q", option.Query)
	}
	if option.Page > 0 {
		query.Set("page", strconv.Itoa(option.Page))
	}
	if option.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(option.PerPage))
	}
	return query
}
