package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	tissue "github.com/mohemohe/go-tissue"
)

type CreateCheckinOption struct {
	CheckedInAt        *time.Time `json:"checked_in_at,omitempty"`
	Tags               []string   `json:"tags,omitempty"`
	Link               string     `json:"link,omitempty"`
	Note               string     `json:"note,omitempty"`
	IsPrivate          bool       `json:"is_private"`
	IsTooSensitive     bool       `json:"is_too_sensitive"`
	DiscardElapsedTime bool       `json:"discard_elapsed_time"`
}

type UpdateCheckinOption struct {
	CheckedInAt        *time.Time `json:"checked_in_at,omitempty"`
	Tags               *[]string  `json:"tags,omitempty"`
	Link               *string    `json:"link,omitempty"`
	Note               *string    `json:"note,omitempty"`
	IsPrivate          *bool      `json:"is_private,omitempty"`
	IsTooSensitive     *bool      `json:"is_too_sensitive,omitempty"`
	DiscardElapsedTime *bool      `json:"discard_elapsed_time,omitempty"`
}

func (c *Client) CreateCheckin(ctx context.Context, option *CreateCheckinOption) (*tissue.Checkin, error) {
	if option == nil {
		option = &CreateCheckinOption{}
	}
	result := &tissue.Checkin{}
	if err := c.sendJSON(ctx, http.MethodPost, "/v1/checkins", option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetCheckin(ctx context.Context, id int64) (*tissue.Checkin, error) {
	result := &tissue.Checkin{}
	if err := c.getJSON(ctx, "/v1/checkins/"+strconv.FormatInt(id, 10), nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCheckin(ctx context.Context, id int64, option *UpdateCheckinOption) (*tissue.Checkin, error) {
	if option == nil {
		option = &UpdateCheckinOption{}
	}
	result := &tissue.Checkin{}
	if err := c.sendJSON(ctx, http.MethodPatch, "/v1/checkins/"+strconv.FormatInt(id, 10), option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCheckin(ctx context.Context, id int64) error {
	return c.deleteRequest(ctx, "/v1/checkins/"+strconv.FormatInt(id, 10))
}
