package go_tissue

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

type (
	Client struct {
		option   *ClientOption
		cookie   *cookiejar.Jar
		loggedIn bool
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
	ListTagsOption struct {
		Page int
	}
	ListTagsResult struct {
		Name  string
		Count int
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
		option:   option,
		cookie:   cookie,
		loggedIn: false,
	}, nil
}

func (this *Client) httpClient() (*http.Client, error) {
	loginPath := "/login"
	client := &http.Client{
		Jar: this.cookie,
	}
	if this.loggedIn {
		return client, nil
	}
	token, err := this.fetchToken(loginPath, client)

	postForm := url.Values{
		"_token":   []string{token},
		"email":    []string{this.option.Email},
		"password": []string{this.option.Password},
	}

	res, err := this.httpRequest(context.TODO(), client, http.MethodPost, loginPath, strings.NewReader(postForm.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("unauthorized")
	}
	this.loggedIn = true
	return client, nil
}

func (this *Client) httpRequest(ctx context.Context, client *http.Client, method string, spath string, body io.Reader) (*http.Response, error) {
	u, err := url.Parse(this.option.BaseURL)
	if err != nil {
		return nil, err
	}
	endpoint := this.option.BaseURL + spath
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if (method != http.MethodGet) && (method != http.MethodHead) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", "go-tissue")
	req.Header.Set("Host", u.Host)

	return client.Do(req)
}

func (this *Client) fetchToken(spath string, client *http.Client) (string, error) {
	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
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
	spath := "/checkin"

	client, err := this.httpClient()
	if err != nil {
		return -1, err
	}

	token, err := this.fetchToken(spath, client)
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

	res, err := this.httpRequest(context.TODO(), client, http.MethodPost, spath, strings.NewReader(postForm.Encode()))
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		return -1, errors.New("something wrong: " + res.Status)
	}

	location := res.Header.Get("location")
	if location == "" {
		return -1, nil
	}
	path := strings.Split(location, "/")
	return strconv.ParseInt(path[len(path)-1], 10, 64)
}

func (this *Client) ListTags(option *ListTagsOption) (result []ListTagsResult, err error) {
	spath := "/tag" + "?page=" + strconv.Itoa(option.Page)

	client, err := this.httpClient()
	if err != nil {
		return nil, err
	}

	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	tagNameNodes, err := htmlquery.QueryAll(doc, "//*[@class=\"tag-name\"]")
	if err != nil {
		return nil, err
	}

	result = make([]ListTagsResult, len(tagNameNodes))
	for i, tagNameNode := range tagNameNodes {
		countNode := htmlquery.FindOne(tagNameNode.Parent, "//*[@class=\"checkins-count\"]")
		countData := countNode.FirstChild.Data
		if strings.HasPrefix(countData, "(") {
			countData = countData[1:]
		}
		if strings.HasSuffix(countData, ")") {
			countData = countData[:len(countData)-1]
		}
		count, err := strconv.Atoi(countData)
		if err != nil {
			count = -1
		}
		result[i] = ListTagsResult{
			Name:  tagNameNode.FirstChild.Data,
			Count: count,
		}
	}

	return result, nil
}
