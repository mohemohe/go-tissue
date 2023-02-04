package go_tissue

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()

	os.Exit(m.Run())
}

func TestClient_CheckIn(t *testing.T) {
	if os.Getenv("TISSUE_SKIP_CHECKIN_TEST") == "1" {
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

var createdCollection Collection

func TestClient_CreateCollection(t *testing.T) {
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

	createOption := &CreateCollectionOption{
		Title:   "test collection",
		Private: true,
	}
	collection, err := client.CreateCollection(createOption)
	if err != nil {
		t.Fatal(err)
	}
	if collection.ID <= 0 {
		t.Error("invalid collection ID")
	}
	if collection.Title != createOption.Title {
		t.Error("invalid collection title")
	}
	if collection.Private != createOption.Private {
		t.Error("invalid collection visibility")
	}
	if collection.UserID == "" {
		t.Error("invalid collection user ID")
	}

	createdCollection = collection

	t.Log(collection)
}

func TestClient_EditCollection(t *testing.T) {
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

	editOption := &EditCollectionOption{
		ID:      createdCollection.ID,
		Title:   "test collection 2",
		Private: true,
	}
	collection, err := client.EditCollection(editOption)
	if err != nil {
		t.Fatal(err)
	}
	if collection.ID != createdCollection.ID {
		t.Error("invalid collection ID")
	}
	if collection.Title != editOption.Title {
		t.Error("invalid collection title")
	}
	if collection.Private != editOption.Private {
		t.Error("invalid collection visibility")
	}
	if collection.UserID != createdCollection.UserID {
		t.Error("invalid collection user ID")
	}

	t.Log(collection)
}

func TestClient_ListCollection(t *testing.T) {
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

	collections, err := client.ListCollection()
	if err != nil {
		t.Fatal(err)
	}
	if len(collections) == 0 {
		return
	}
	for _, collection := range collections {
		if collection.ID <= 0 {
			t.Error("invalid collection ID")
		}
		if collection.Title == "" {
			t.Error("invalid collection title")
		}
		if collection.UserID == "" {
			t.Error("invalid collection user ID")
		}

		t.Log(collection)
	}
}

func TestClient_DeleteCollection(t *testing.T) {
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

	if err := client.DeleteCollection(createdCollection.ID); err != nil {
		t.Fatal(err)
	}
}
