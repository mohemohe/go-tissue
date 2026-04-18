package go_tissue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
	"sync"
)

var errNilOption = errors.New("option is required")

type ClientOption struct {
	BaseURL  string
	Email    string
	Password string
}

type Client struct {
	option     *ClientOption
	httpClient *http.Client
	baseURL    *url.URL

	mu       sync.Mutex
	loggedIn bool
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		option:  option,
		baseURL: u,
		httpClient: &http.Client{
			Jar: jar,
		},
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

func (c *Client) xsrfToken() (string, error) {
	cookies := c.httpClient.Jar.Cookies(c.baseURL)
	for _, ck := range cookies {
		if ck.Name == "XSRF-TOKEN" {
			decoded, err := url.QueryUnescape(ck.Value)
			if err != nil {
				return "", err
			}
			return decoded, nil
		}
	}
	return "", errors.New("XSRF-TOKEN cookie not found")
}

func (c *Client) ensureLoggedIn(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loggedIn {
		return nil
	}
	if err := c.login(ctx); err != nil {
		return err
	}
	c.loggedIn = true
	return nil
}

func (c *Client) login(ctx context.Context) error {
	loginURL := c.resolveURL("/login", nil)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loginURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "go-tissue")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()

	token, err := c.xsrfToken()
	if err != nil {
		return fmt.Errorf("fetch XSRF token: %w", err)
	}

	form := url.Values{
		"email":    []string{c.option.Email},
		"password": []string{c.option.Password},
	}
	postReq, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	postReq.Header.Set("User-Agent", "go-tissue")
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("X-XSRF-TOKEN", token)
	postReq.Header.Set("Accept", "text/html,application/xhtml+xml")

	origCheckRedirect := c.httpClient.CheckRedirect
	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() { c.httpClient.CheckRedirect = origCheckRedirect }()

	postRes, err := c.httpClient.Do(postReq)
	if err != nil {
		return err
	}
	defer postRes.Body.Close()
	_, _ = io.Copy(io.Discard, postRes.Body)

	if postRes.StatusCode >= 400 {
		return fmt.Errorf("login failed: %s", postRes.Status)
	}
	if postRes.StatusCode == http.StatusOK {
		return errors.New("login failed: credentials rejected")
	}
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, spath string, query url.Values, body io.Reader, contentType string) (*http.Response, error) {
	if err := c.ensureLoggedIn(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, c.resolveURL(spath, query), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-tissue")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if body != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if method != http.MethodGet && method != http.MethodHead {
		if token, err := c.xsrfToken(); err == nil {
			req.Header.Set("X-XSRF-TOKEN", token)
		}
	}
	return c.httpClient.Do(req)
}

func (c *Client) getJSON(ctx context.Context, spath string, query url.Values, out interface{}) error {
	res, err := c.doRequest(ctx, http.MethodGet, spath, query, nil, "")
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
	res, err := c.doRequest(ctx, method, spath, nil, body, "application/json")
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
	res, err := c.doRequest(ctx, http.MethodDelete, spath, nil, nil, "")
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
