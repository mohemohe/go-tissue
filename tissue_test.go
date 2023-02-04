package go_tissue

import (
	"os"
	"testing"
	"time"
)

func TestClient_CheckIn(t *testing.T) {
	if os.Getenv("TISSUE_SKIP_CHECKIN_TEST") != "" {
		t.Skip("skip checkin test")
	}

	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	checkInID, err := client.CheckIn(&CheckInOption{
		DateTime:     time.Now(),
		Tags:         []string{"test", "hoge"},
		Link:         "",
		Note:         "本番環境でテストしてすまん",
		Private:      true,
		TooSensitive: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if checkInID < 0 {
		t.Error("could not parse checkIn ID")
	}

	t.Log("checkin ID:", checkInID)
}

func TestClient_ListTags(t *testing.T) {
	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	result, err := client.ListTags(&ListTagsOption{
		Page: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("could not parse tags page")
	}
	if len(result) > 0 && result[0].Count == -1 {
		t.Error("could not parse checkin count")
	}

	for _, tag := range result {
		t.Log(tag)
	}
}

func TestClient_Search(t *testing.T) {
	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	result, err := client.Search(&SearchOption{
		Keyword: "VOICEROID",
		Page:    1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("could not parse search page")
	}
	for _, checkIn := range result {
		if checkIn.ID == -1 {
			t.Error("could not parse checkin ID")
		}
		if checkIn.DateTime == time.Unix(0, 0) {
			t.Error("could not parse checkin datetime")
		}
		if checkIn.User.ID == "" {
			t.Error("could not parse checkin user id")
		}
		if checkIn.User.DisplayName == "" {
			t.Error("could not parse checkin user name")
		}

		t.Log(checkIn)
	}
}

func TestClient_PublicTimeline(t *testing.T) {
	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	result, err := client.PublicTimeline(&TimelineOption{
		Page: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("could not parse search page")
	}
	for _, checkIn := range result {
		if checkIn.ID == -1 {
			t.Error("could not parse checkin ID")
		}
		if checkIn.DateTime == time.Unix(0, 0) {
			t.Error("could not parse checkin datetime")
		}
		if checkIn.User.ID == "" {
			t.Error("could not parse checkin user id")
		}
		if checkIn.User.DisplayName == "" {
			t.Error("could not parse checkin user name")
		}

		t.Log(checkIn)
	}
}

func TestClient_GetStatus(t *testing.T) {
	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	status, err := client.GetStatus()
	if err != nil {
		t.Fatal(err)
	}
	if status.User.ID == "" {
		t.Error("could not parse user ID")
	}
	if status.User.DisplayName == "" {
		t.Error("could not parse user display name")
	}
	if status.Session.Current == time.Duration(-1) {
		t.Error("could not parse current session")
	}
	if status.Session.ResetTo == time.Unix(0, 0) {
		t.Error("could not parse session reset datetime")
	}
	if status.Overview.Average == time.Duration(-1) {
		t.Error("could not parse average duration of overview")
	}
	if status.Overview.Median == time.Duration(-1) {
		t.Error("could not parse average duration of overview")
	}
	if status.Overview.Longest == time.Duration(-1) {
		t.Error("could not parse average duration of overview")
	}
	if status.Overview.Shortest == time.Duration(-1) {
		t.Error("could not parse average duration of overview")
	}
	if status.Overview.Sum == time.Duration(-1) {
		t.Error("could not parse average duration of overview")
	}
	if status.Overview.Count == -1 {
		t.Error("could not parse count of overview")
	}

	t.Log(status)
}
