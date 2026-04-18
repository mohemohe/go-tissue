package api

import (
	"context"
	"net/url"
	"strconv"
	"time"

	tissue "github.com/mohemohe/go-tissue"
)

type UserCheckinsOption struct {
	Page    int
	PerPage int
	HasLink *bool
	Since   time.Time
	Until   time.Time
	Order   string
}

type PageOption struct {
	Page    int
	PerPage int
}

func (c *Client) GetUser(ctx context.Context, name string) (*tissue.User, error) {
	result := &tissue.User{}
	if err := c.getJSON(ctx, "/v1/users/"+name, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserCheckins(ctx context.Context, name string, option *UserCheckinsOption) ([]tissue.Checkin, error) {
	query := url.Values{}
	if option != nil {
		if option.Page > 0 {
			query.Set("page", strconv.Itoa(option.Page))
		}
		if option.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(option.PerPage))
		}
		if option.HasLink != nil {
			query.Set("has_link", strconv.FormatBool(*option.HasLink))
		}
		if !option.Since.IsZero() {
			query.Set("since", option.Since.Format("2006-01-02"))
		}
		if !option.Until.IsZero() {
			query.Set("until", option.Until.Format("2006-01-02"))
		}
		if option.Order != "" {
			query.Set("order", option.Order)
		}
	}
	result := []tissue.Checkin{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/checkins", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserLikes(ctx context.Context, name string, option *PageOption) ([]tissue.Checkin, error) {
	query := applyPageOption(nil, option)
	result := []tissue.Checkin{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/likes", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserCollections(ctx context.Context, name string, option *PageOption) ([]tissue.Collection, error) {
	query := applyPageOption(nil, option)
	result := []tissue.Collection{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/collections", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func applyPageOption(query url.Values, option *PageOption) url.Values {
	if query == nil {
		query = url.Values{}
	}
	if option == nil {
		return query
	}
	if option.Page > 0 {
		query.Set("page", strconv.Itoa(option.Page))
	}
	if option.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(option.PerPage))
	}
	return query
}
