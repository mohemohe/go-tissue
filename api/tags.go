package api

import "context"

func (c *Client) RecentTags(ctx context.Context) ([]string, error) {
	result := []string{}
	if err := c.getJSON(ctx, "/recent-tags", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
