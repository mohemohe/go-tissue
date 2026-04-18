package go_tissue

import (
	"context"
	"net/http"
	"strconv"
	"time"
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

func (c *Client) CreateCheckin(ctx context.Context, option *CreateCheckinOption) (*Checkin, error) {
	if option == nil {
		option = &CreateCheckinOption{}
	}
	result := &Checkin{}
	if err := c.sendJSON(ctx, http.MethodPost, "/api/checkins", option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetCheckin(ctx context.Context, id int64) (*Checkin, error) {
	result := &Checkin{}
	if err := c.getJSON(ctx, "/api/checkins/"+strconv.FormatInt(id, 10), nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCheckin(ctx context.Context, id int64, option *UpdateCheckinOption) (*Checkin, error) {
	if option == nil {
		option = &UpdateCheckinOption{}
	}
	result := &Checkin{}
	if err := c.sendJSON(ctx, http.MethodPatch, "/api/checkins/"+strconv.FormatInt(id, 10), option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCheckin(ctx context.Context, id int64) error {
	return c.deleteRequest(ctx, "/api/checkins/"+strconv.FormatInt(id, 10))
}
