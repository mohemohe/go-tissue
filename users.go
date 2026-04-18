package go_tissue

import (
	"context"
	"net/url"
	"strconv"
)

type UserCheckinsOption struct {
	Page    int
	PerPage int
	HasLink *bool
}

func (c *Client) UserCheckins(ctx context.Context, user string, option *UserCheckinsOption) ([]UserCheckin, error) {
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
	}
	result := []UserCheckin{}
	if err := c.getJSON(ctx, "/api/users/"+user+"/checkins", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}
