package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"time"
)

type CheckInOption struct {
	DateTime     time.Time `json:"-"`
	Tags         []string  `json:"tags"`
	Link         string    `json:"link"`
	Note         string    `json:"note"`
	Private      bool      `json:"is_private"`
	TooSensitive bool      `json:"is_too_sensitive"`
}

type webhookCheckInRequest struct {
	CheckInOption
	CheckedInAt string `json:"checked_in_at,omitempty"`
}

type webhookCheckInResponse struct {
	Status  int     `json:"status"`
	CheckIn CheckIn `json:"checkin"`
}

type CheckIn struct {
	CheckInOption
	ID          uint      `json:"id"`
	Source      string    `json:"source"`
	CheckedInAt string    `json:"checked_in_at"`
	DateTime    time.Time `json:"-"`
}

func (c *Client) CheckIn(ctx context.Context, option *CheckInOption) (*CheckIn, error) {
	if c.option.WebhookID == "" {
		return nil, errors.New("webhook id is required")
	}
	if option == nil {
		option = &CheckInOption{}
	}

	spath := path.Join("/webhooks/checkin", c.option.WebhookID)

	body := webhookCheckInRequest{CheckInOption: *option}
	if !option.DateTime.IsZero() {
		body.CheckedInAt = option.DateTime.Format("2006-01-02T15:04:05-0700")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.resolveURL(spath, nil), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-tissue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, readErrorResponse(res)
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := webhookCheckInResponse{}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.CheckIn.CheckedInAt != "" {
		if t, err := time.Parse("2006-01-02T15:04:05-0700", r.CheckIn.CheckedInAt); err == nil {
			r.CheckIn.DateTime = t
		}
	}
	return &r.CheckIn, nil
}
