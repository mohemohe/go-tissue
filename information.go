package go_tissue

import "context"

func (c *Client) LatestInformation(ctx context.Context) ([]Information, error) {
	result := []Information{}
	if err := c.getJSON(ctx, "/api/information/latest", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
