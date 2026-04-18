package go_tissue

import (
	"context"
	"net/url"
	"time"
)

const statsDateLayout = "2006-01-02"

func (c *Client) DailyCheckinStats(ctx context.Context) ([]DailyCheckinCount, error) {
	result := []DailyCheckinCount{}
	if err := c.getJSON(ctx, "/api/stats/checkin/daily", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type UserDailyCheckinStatsOption struct {
	Since time.Time
	Until time.Time
}

func (c *Client) UserDailyCheckinStats(ctx context.Context, user string, option *UserDailyCheckinStatsOption) ([]DailyCheckinCount, error) {
	query := url.Values{}
	if option != nil {
		if !option.Since.IsZero() {
			query.Set("since", option.Since.Format(statsDateLayout))
		}
		if !option.Until.IsZero() {
			query.Set("until", option.Until.Format(statsDateLayout))
		}
	}
	result := []DailyCheckinCount{}
	if err := c.getJSON(ctx, "/api/users/"+user+"/stats/checkin/daily", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserTagStats(ctx context.Context, user string) ([]TagCount, error) {
	result := []TagCount{}
	if err := c.getJSON(ctx, "/api/users/"+user+"/stats/tags", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
