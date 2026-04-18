package api

import (
	"context"
	"net/url"
	"time"

	tissue "github.com/mohemohe/go-tissue"
)

type UserStatsPeriodOption struct {
	Since time.Time
	Until time.Time
}

func (c *Client) UserDailyCheckinStats(ctx context.Context, name string, option *UserStatsPeriodOption) ([]tissue.DailyCheckinCount, error) {
	result := []tissue.DailyCheckinCount{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/stats/checkin/daily", buildPeriodQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserHourlyCheckinStats(ctx context.Context, name string, option *UserStatsPeriodOption) ([]HourlyCheckinSummary, error) {
	result := []HourlyCheckinSummary{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/stats/checkin/hourly", buildPeriodQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserTagStats(ctx context.Context, name string, option *UserStatsPeriodOption) ([]tissue.TagCount, error) {
	result := []tissue.TagCount{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/stats/tags", buildPeriodQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UserLinkStats(ctx context.Context, name string, option *UserStatsPeriodOption) ([]LinkCount, error) {
	result := []LinkCount{}
	if err := c.getJSON(ctx, "/v1/users/"+name+"/stats/links", buildPeriodQuery(option), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func buildPeriodQuery(option *UserStatsPeriodOption) url.Values {
	query := url.Values{}
	if option == nil {
		return query
	}
	if !option.Since.IsZero() {
		query.Set("since", option.Since.Format("2006-01-02"))
	}
	if !option.Until.IsZero() {
		query.Set("until", option.Until.Format("2006-01-02"))
	}
	return query
}
