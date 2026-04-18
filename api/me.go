package api

import (
	"context"

	tissue "github.com/mohemohe/go-tissue"
)

func (c *Client) Me(ctx context.Context) (*tissue.Me, error) {
	result := &tissue.Me{}
	if err := c.getJSON(ctx, "/v1/me", nil, result); err != nil {
		return nil, err
	}
	return result, nil
}
