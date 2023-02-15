package go_tissue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	SearchOption struct {
		Keyword string
		Page    int
	}
	User struct {
		ID          string
		DisplayName string
	}
	CheckInResult struct {
		ID       int64
		DateTime time.Time
		Tags     []string
		Link     string
		Note     string
		User     User
	}
	TimelineOption struct {
		Page int
	}
	Session struct {
		Current time.Duration
		ResetTo time.Time
	}
	Overview struct {
		Average  time.Duration
		Median   time.Duration
		Longest  time.Duration
		Shortest time.Duration
		Sum      time.Duration
		Count    int
	}
	Status struct {
		User
		Session  Session
		Overview Overview
	}
	CreateCollectionOption struct {
		Title   string `json:"title"`
		Private bool   `json:"is_private"`
	}
	EditCollectionOption struct {
		ID      int64  `json:"-"`
		Title   string `json:"title"`
		Private bool   `json:"is_private"`
	}
	ListCollectionOption struct {
	}
	Collection struct {
		ID      int64  `json:"id"`
		Private bool   `json:"is_private"`
		Title   string `json:"title"`
		UserID  string `json:"user_name"`
	}
	CollectionItem struct {
		ID           int64    `json:"id"`
		CollectionID int64    `json:"collection_id"`
		Link         string   `json:"link"`
		Note         string   `json:"note"`
		NoteHTML     string   `json:"note_html"`
		Tags         []string `json:"tags"`
		CheckInURL   string   `json:"checkin_url"`
		UserID       string   `json:"user_name"`
	}
	CreateCollectionItemOption struct {
		CollectionID int64    `json:"-"`
		Link         string   `json:"link"`
		Note         string   `json:"note"`
		Tags         []string `json:"tags"`
		flash        bool     `json:"flash"`
	}
	EditCollectionItemOption struct {
		CollectionID int64    `json:"-"`
		ItemID       int64    `json:"-"`
		Note         string   `json:"note"`
		Tags         []string `json:"tags"`
		flash        bool     `json:"flash"`
	}
	ListCollectionItemOption struct {
		CollectionID int64 `json:"-"`
	}
	DeleteCollectionItemOption struct {
		CollectionID int64
		ItemID       int64
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

func (this *Client) initHttpRequest(ctx context.Context, client *http.Client, method string, spath string, body io.Reader) (*http.Request, error) {
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
	req.Header.Set("User-Agent", "go-tissue")
	req.Header.Set("Host", u.Host)

	return req, nil
}

func (this *Client) httpRequest(ctx context.Context, client *http.Client, method string, spath string, body io.Reader) (*http.Response, error) {
	req, err := this.initHttpRequest(ctx, client, method, spath, body)
	if err != nil {
		return nil, err
	}

	if (method != http.MethodGet) && (method != http.MethodHead) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return client.Do(req)
}

func (this *Client) jsonRequest(ctx context.Context, client *http.Client, method string, spath string, body io.Reader) (*http.Response, error) {
	req, err := this.initHttpRequest(ctx, client, method, spath, body)
	if err != nil {
		return nil, err
	}

	if (method != http.MethodGet) && (method != http.MethodHead) {
		req.Header.Set("Content-Type", "application/json")

		csrfToken, err := this.fetchToken("/", client)
		if err != nil {
			return nil, err
		}
		req.Header.Set("x-csrf-token", csrfToken)
	}
	return client.Do(req)
}

func (this *Client) fetchToken(spath string, client *http.Client) (string, error) {
	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	token := doc.Find("form input[name='_token']").First().AttrOr("value", "")
	if token == "" {
		return "", errors.New("_token not found")
	}
	return token, nil
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
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	tagNodes := doc.Find(".tags a")
	result = make([]ListTagsResult, tagNodes.Length())
	tagNodes.Each(func(i int, s *goquery.Selection) {
		result[i].Name = s.Find(".tag-name").Text()

		countText := s.Find(".checkins-count").Text()
		if strings.HasPrefix(countText, "(") {
			countText = countText[1:]
		}
		if strings.HasSuffix(countText, ")") {
			countText = countText[:len(countText)-1]
		}
		count, err := strconv.Atoi(countText)
		if err != nil {
			count = -1
		}
		result[i].Count = count
	})

	return result, nil
}

func (this *Client) parseChackIn(nodes *goquery.Selection) (result []CheckInResult, err error) {
	result = make([]CheckInResult, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		dateTimeNode := s.Find("a[href*='/checkin/']").First()
		checkInID, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(dateTimeNode.AttrOr("href", "-1"), this.option.BaseURL+"/checkin/")), 10, 64)
		if err != nil {
			checkInID = -1
		}
		result[i].ID = checkInID

		dateTime, err := time.Parse("2006/01/02 15:04", strings.TrimSpace(dateTimeNode.Text()))
		if err != nil {
			dateTime = time.Unix(0, 0)
		}
		result[i].DateTime = dateTime

		userNode := s.Find("a[href*='/user/']").First()
		result[i].User.DisplayName = strings.TrimSpace(userNode.Text())
		result[i].User.ID = strings.TrimSpace(strings.TrimPrefix(userNode.AttrOr("href", "-1"), this.option.BaseURL+"/user/"))

		tagNodes := s.Find(".tis-checkin-tags")
		tagNodes.Each(func(j int, t *goquery.Selection) {
			badgeNodes := t.Find(".badge")
			result[i].Tags = make([]string, badgeNodes.Length())
			badgeNodes.Each(func(k int, u *goquery.Selection) {
				result[i].Tags[k] = strings.TrimSpace(u.Text())
			})
		})

		result[i].Link = s.Find(".oi-link-intact + a").First().AttrOr("href", "")

		result[i].Note = strings.TrimSpace(s.Find(".tis-checkin-tags + div > p").First().Text())
	})

	return result, nil
}

func (this *Client) Search(option *SearchOption) (result []CheckInResult, err error) {
	spath := "/search/checkin" + "?q=" + url.QueryEscape(option.Keyword)

	client, err := this.httpClient()
	if err != nil {
		return nil, err
	}

	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	checkInNodes := doc.Find(".list-group-item")
	return this.parseChackIn(checkInNodes)
}

func (this *Client) PublicTimeline(option *TimelineOption) (result []CheckInResult, err error) {
	spath := "/timeline/public" + "?page=" + strconv.Itoa(option.Page)

	client, err := this.httpClient()
	if err != nil {
		return nil, err
	}

	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	checkInNodes := doc.Find(".container-fluid > .row > div[class*='col-']")
	return this.parseChackIn(checkInNodes)
}

func trimDurationFunc(r rune) bool {
	return r == '日' || r == '時' || r == '間' || r == '分' || r == '経' || r == '過'
}

func toRawNumberString(s string) string {
	return strings.ReplaceAll(strings.TrimFunc(s, trimDurationFunc), ",", "")
}

func elemToDuration(elem []string) time.Duration {
	if len(elem) == 3 {
		days, err1 := strconv.Atoi(toRawNumberString(elem[0]))
		hours, err2 := strconv.Atoi(toRawNumberString(elem[1]))
		minutes, err3 := strconv.Atoi(toRawNumberString(elem[2]))
		if err1 == nil && err2 == nil && err3 == nil {
			return (time.Duration(days) * time.Hour * 24) + (time.Duration(hours) * time.Hour) + (time.Duration(minutes) * time.Minute)
		}
	}
	return time.Duration(0)
}

func (this *Client) GetStatus() (result *Status, err error) {
	spath := "/"

	result = &Status{
		Session: Session{
			Current: time.Duration(-1),
			ResetTo: time.Unix(0, 0),
		},
		Overview: Overview{
			Average:  time.Duration(-1),
			Median:   time.Duration(-1),
			Longest:  time.Duration(-1),
			Shortest: time.Duration(-1),
			Sum:      time.Duration(-1),
			Count:    -1,
		},
	}

	client, err := this.httpClient()
	if err != nil {
		return result, err
	}

	res, err := this.httpRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return result, err
	}
	if res.StatusCode != http.StatusOK {
		return result, errors.New("something wrong: " + res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return result, err
	}

	infoNode := doc.Find(".container > .row > div[class*='col-'] > .card > .card-body").First()

	result.User.ID = strings.TrimSpace(infoNode.Find(".tis-profile-mini-display-name a").First().Text())

	result.User.DisplayName = strings.TrimSpace(strings.TrimPrefix(infoNode.Find(".tis-profile-mini-name a").First().Text(), "@"))

	currentSessionTexts := strings.Split(strings.TrimSpace(infoNode.Find("h6 + p").First().Text()), " ")
	result.Session.Current = elemToDuration(currentSessionTexts)

	sessionResetTexts := strings.Split(strings.TrimSpace(infoNode.Find("h6 + p + p").First().Text()), " ")
	if len(currentSessionTexts) == 3 {
		sessionResetTexts[0] = strings.Trim(sessionResetTexts[0], "(")
		sessionResetText := strings.Join(sessionResetTexts[0:2], " ")
		if resetTo, err := time.Parse("2006/01/02 15:04", sessionResetText); err == nil {
			result.Session.ResetTo = resetTo
		}
	}

	overviewNodes := infoNode.Find(".tis-profile-stats-table tbody tr")
	overviewNodes.Each(func(i int, s *goquery.Selection) {
		durationLabel := strings.TrimSpace(s.Find("th").First().Text())
		durationTexts := strings.Split(strings.TrimSpace(s.Find("td").First().Text()), " ")

		switch durationLabel {
		case "平均記録":
			result.Overview.Average = elemToDuration(durationTexts)
		case "中央値":
			result.Overview.Median = elemToDuration(durationTexts)
		case "最長記録":
			result.Overview.Longest = elemToDuration(durationTexts)
		case "最短記録":
			result.Overview.Shortest = elemToDuration(durationTexts)
		case "合計時間":
			result.Overview.Sum = elemToDuration(durationTexts)
		case "通算回数":
			if count, err := strconv.Atoi(strings.Trim(durationTexts[0], "回")); err == nil {
				result.Overview.Count = count
			}
		}
	})

	return result, nil
}

func (this *Client) CreateCollection(option *CreateCollectionOption) (result *Collection, err error) {
	spath := "/api/collections"

	client, err := this.httpClient()
	if err != nil {
		return result, err
	}

	payload, err := json.Marshal(option)
	if err != nil {
		return result, err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodPost, spath, bytes.NewReader(payload))
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return result, errors.New("something wrong: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	result = &Collection{}
	if err := json.Unmarshal(body, result); err != nil {
		return result, err
	}
	return result, nil
}

func (this *Client) EditCollection(option *EditCollectionOption) (result *Collection, err error) {
	spath := "/api/collections" + "/" + strconv.FormatInt(option.ID, 10)

	client, err := this.httpClient()
	if err != nil {
		return result, err
	}

	payload, err := json.Marshal(option)
	if err != nil {
		return result, err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodPut, spath, bytes.NewReader(payload))
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, errors.New("something wrong: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	result = &Collection{}
	if err := json.Unmarshal(body, result); err != nil {
		return result, err
	}
	return result, nil
}

func (this *Client) ListCollection(option *ListCollectionOption) (result []Collection, err error) {
	spath := "/api/collections"

	client, err := this.httpClient()
	if err != nil {
		return result, err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, errors.New("something wrong: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	result = []Collection{}
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (this *Client) DeleteCollection(collectionID int64) error {
	spath := "/api/collections" + "/" + strconv.FormatInt(collectionID, 10)

	client, err := this.httpClient()
	if err != nil {
		return err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodDelete, spath, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return errors.New("something wrong: " + res.Status)
	}
	return nil
}

func (this *Client) CreateCollectionItem(option *CreateCollectionItemOption) error {
	spath := "/api/collections" + "/" + strconv.FormatInt(option.CollectionID, 10) + "/items"

	client, err := this.httpClient()
	if err != nil {
		return err
	}

	option.flash = true
	payload, err := json.Marshal(option)
	if err != nil {
		return err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodPost, spath, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("something wrong: " + res.Status)
	}
	return nil
}

func (this *Client) EditCollectionItem(option *EditCollectionItemOption) (result *CollectionItem, err error) {
	spath := "/api/collections" + "/" + strconv.FormatInt(option.CollectionID, 10) + "/items" + "/" + strconv.FormatInt(option.ItemID, 10)

	client, err := this.httpClient()
	if err != nil {
		return nil, err
	}

	option.flash = true
	payload, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodPut, spath, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	result = &CollectionItem{}
	err = json.Unmarshal(body, result)

	return result, nil
}

func (this *Client) ListCollectionItem(option *ListCollectionItemOption) (result []CollectionItem, err error) {
	spath := "/api/collections" + "/" + strconv.FormatInt(option.CollectionID, 10) + "/items"

	client, err := this.httpClient()
	if err != nil {
		return nil, err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodGet, spath, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("something wrong: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	result = []CollectionItem{}
	err = json.Unmarshal(body, &result)

	return result, nil
}

// TODO: GetCollectionItem

func (this *Client) DeleteCollectionItem(option *DeleteCollectionItemOption) error {
	spath := "/api/collections" + "/" + strconv.FormatInt(option.CollectionID, 10) + "/items" + "/" + strconv.FormatInt(option.ItemID, 10)

	client, err := this.httpClient()
	if err != nil {
		return err
	}

	res, err := this.jsonRequest(context.TODO(), client, http.MethodDelete, spath, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return errors.New("something wrong: " + res.Status)
	}
	return nil
}
