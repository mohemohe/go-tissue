package go_tissue

import "context"

func (c *Client) Me(ctx context.Context) (*Me, error) {
	result := &Me{}
	if err := c.getJSON(ctx, "/api/me", nil, result); err != nil {
		return nil, err
	}
	return result, nil
}
