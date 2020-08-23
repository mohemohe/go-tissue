package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"time"
)

type (
	Client struct {
		option *ClientOption
		cookie *cookiejar.Jar
	}
	ClientOption struct {
		BaseURL  string
		WebhookID string
	}
	CheckInOption struct {
		DateTime     time.Time `json:"-"`
		Tags         []string `json:"tags"`
		Link         string `json:"link"`
		Note         string `json:"note"`
		Private      bool `json:"is_private"`
		TooSensitive bool `json:"is_too_sensitive"`
	}
	checkInRequest struct {
		CheckInOption
		checkedInAt  string `json:"checked_in_at"`
	}
	checkInResponse struct {
		Status int `json:"status"`
		CheckIn CheckIn `json:"checkin"`
	}
	CheckIn struct {
		CheckInOption
		ID uint `json:"id"`
		Source string `json:"source"`
		CheckedInAt  string `json:"checked_in_at"`
	}
)

func NewClient(option *ClientOption) (*Client, error) {
	cookie, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if option.BaseURL == "" {
		option.BaseURL = "https://shikorism.net"
	}
	u, err := url.Parse(option.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api")
	option.BaseURL = u.String()

	return &Client{
		option: option,
		cookie: cookie,
	}, nil
}

func (this *Client) httpRequest(ctx context.Context, method string, spath string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(this.option.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-tissue")

	return req, nil
}

func (this *Client) CheckIn(ctx context.Context, option *CheckInOption) (result *CheckIn, err error) {
	p := path.Join("/webhooks/checkin/", this.option.WebhookID)

	checkInRequst := checkInRequest{
		CheckInOption: *option,
		checkedInAt: option.DateTime.Format("2006-01-02T15:04:05-0700"),
	}
	b, err := json.Marshal(checkInRequst)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)

	req, err := this.httpRequest(ctx, http.MethodPost, p, reader)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := checkInResponse{}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	r.CheckIn.DateTime, err = time.Parse("2006-01-02T15:04:05-0700", r.CheckIn.CheckedInAt)
	if err != nil {
		return nil, err
	}

	return &r.CheckIn, nil
}
