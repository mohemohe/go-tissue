package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type ClientOption struct {
	BaseURL     string
	WebhookID   string
	AccessToken string
}

type Client struct {
	option     *ClientOption
	httpClient *http.Client
	baseURL    *url.URL
}

func NewClient(option *ClientOption) (*Client, error) {
	if option == nil {
		return nil, errors.New("option is required")
	}
	if option.BaseURL == "" {
		option.BaseURL = "https://shikorism.net"
	}
	u, err := url.Parse(option.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api")

	return &Client{
		option:     option,
		httpClient: &http.Client{},
		baseURL:    u,
	}, nil
}

func (c *Client) resolveURL(spath string, query url.Values) string {
	u := *c.baseURL
	u.Path = path.Join(u.Path, spath)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

func (c *Client) doRequest(ctx context.Context, method, spath string, query url.Values, body io.Reader, contentType string, needAuth bool) (*http.Response, error) {
	if needAuth && c.option.AccessToken == "" {
		return nil, errors.New("access token is required")
	}
	req, err := http.NewRequestWithContext(ctx, method, c.resolveURL(spath, query), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-tissue")
	req.Header.Set("Accept", "application/json")
	if body != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if needAuth {
		req.Header.Set("Authorization", "Bearer "+c.option.AccessToken)
	}
	return c.httpClient.Do(req)
}

func (c *Client) getJSON(ctx context.Context, spath string, query url.Values, out interface{}) error {
	res, err := c.doRequest(ctx, http.MethodGet, spath, query, nil, "", true)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return readErrorResponse(res)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) sendJSON(ctx context.Context, method, spath string, in, out interface{}) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	res, err := c.doRequest(ctx, method, spath, nil, body, "application/json", true)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return readErrorResponse(res)
	}
	if out == nil || res.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) deleteRequest(ctx context.Context, spath string) error {
	res, err := c.doRequest(ctx, http.MethodDelete, spath, nil, nil, "", true)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return readErrorResponse(res)
	}
	return nil
}

func readErrorResponse(res *http.Response) error {
	b, _ := io.ReadAll(res.Body)
	if len(b) > 0 {
		return fmt.Errorf("unexpected status %s: %s", res.Status, strings.TrimSpace(string(b)))
	}
	return fmt.Errorf("unexpected status %s", res.Status)
}
