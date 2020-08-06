package go_tissue

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	Client struct {
		option *ClientOption
		cookie *cookiejar.Jar
	}
	ClientOption struct {
		BaseURL  string
		Email    string
		Password string
	}
	CheckInOption struct {
		DateTime     time.Time
		Tags         []string
		Link         string
		Note         string
		Private      bool
		TooSensitive bool
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

	return &Client{
		option: option,
		cookie: cookie,
	}, nil
}

func (this *Client) httpClient() (*http.Client, error) {
	loginURL := this.option.BaseURL + "/login"
	client := &http.Client{
		Jar: this.cookie,
	}
	token, err := this.fetchToken(loginURL, client)

	postForm := url.Values{
		"_token":   []string{token},
		"email":    []string{this.option.Email},
		"password": []string{this.option.Password},
	}
	res, err := client.PostForm(loginURL, postForm)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("unauthorized")
	}

	return client, nil
}

func (this *Client) fetchToken(url string, client *http.Client) (string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return "", err
	}
	input := htmlquery.FindOne(doc, "//input[@name='_token']/@value")
	if input == nil {
		return "", errors.New("_token not found")
	}
	return input.FirstChild.Data, nil
}

func (this *Client) CheckIn(option *CheckInOption) (checkInID int64, err error) {
	checkInURL := this.option.BaseURL + "/checkin"

	client, err := this.httpClient()
	if err != nil {
		return -1, err
	}

	token, err := this.fetchToken(checkInURL, client)
	if err != nil {
		return -1, err
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	date := option.DateTime.Format("2006/01/02")
	time := option.DateTime.Format("15:04")
	tags := strings.Join(option.Tags, " ")

	postForm := url.Values{
		"_token": []string{token},
		"date":   []string{date},
		"time":   []string{time},
		"tags":   []string{tags},
		"link":   []string{option.Link},
		"note":   []string{option.Note},
	}
	if option.Private {
		postForm["is_private"] = []string{"on"}
	}
	if option.TooSensitive {
		postForm["is_too_sensitive"] = []string{"on"}
	}

	res, err := client.PostForm(checkInURL, postForm)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		return -1, errors.New("something wrong")
	}

	location := res.Header.Get("location")
	if location == "" {
		return -1, nil
	}
	path := strings.Split(location, "/")
	return strconv.ParseInt(path[len(path)-1], 10, 64)
}
