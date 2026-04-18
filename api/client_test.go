package api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func newTokenClient(t *testing.T) *Client {
	t.Helper()
	if os.Getenv("TISSUE_ACCESS_TOKEN") == "" {
		t.Skip("TISSUE_ACCESS_TOKEN not set")
	}
	client, err := NewClient(&ClientOption{
		BaseURL:     os.Getenv("TISSUE_BASE_URL"),
		AccessToken: os.Getenv("TISSUE_ACCESS_TOKEN"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestClient_Me(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if me.Name == "" {
		t.Error("empty user name")
	}
	t.Log(me)
}

func TestClient_GetUser(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	user, err := client.GetUser(context.Background(), me.Name)
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != me.Name {
		t.Errorf("unexpected user name: %s", user.Name)
	}
}

func TestClient_UserCheckins(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserCheckins(context.Background(), me.Name, &UserCheckinsOption{
		Page:    1,
		PerPage: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d checkins", len(result))
}

func TestClient_UserLikes(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserLikes(context.Background(), me.Name, &PageOption{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d likes", len(result))
}

func TestClient_UserCollections(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserCollections(context.Background(), me.Name, &PageOption{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d collections", len(result))
}

func TestClient_UserDailyCheckinStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	until := time.Now()
	since := until.AddDate(-1, 0, 0)
	result, err := client.UserDailyCheckinStats(context.Background(), me.Name, &UserStatsPeriodOption{
		Since: since,
		Until: until,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d daily stats", len(result))
}

func TestClient_UserHourlyCheckinStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserHourlyCheckinStats(context.Background(), me.Name, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d hourly stats", len(result))
}

func TestClient_UserTagStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserTagStats(context.Background(), me.Name, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tag stats: %v", result)
}

func TestClient_UserLinkStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserLinkStats(context.Background(), me.Name, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("link stats: %v", result)
}

func TestClient_SearchCheckins(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)

	result, err := client.SearchCheckins(context.Background(), &SearchOption{
		Query:   "test",
		Page:    1,
		PerPage: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d search results", len(result))
}

func TestClient_CheckinLifecycle(t *testing.T) {
	if os.Getenv("TISSUE_SKIP_CHECKIN_TEST") == "1" {
		t.Skip("skip checkin test")
	}
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)
	ctx := context.Background()

	created, err := client.CreateCheckin(ctx, &CreateCheckinOption{
		Tags:           []string{"test", "hoge"},
		Note:           "go-tissue api v1 test",
		IsPrivate:      true,
		IsTooSensitive: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("empty checkin id")
	}

	got, err := client.GetCheckin(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("id mismatch: got %d want %d", got.ID, created.ID)
	}

	newNote := "go-tissue api v1 test (updated)"
	updated, err := client.UpdateCheckin(ctx, created.ID, &UpdateCheckinOption{
		Note: &newNote,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Note != newNote {
		t.Errorf("note not updated: %q", updated.Note)
	}

	if err := client.DeleteCheckin(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
}

func TestClient_CollectionLifecycle(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTokenClient(t)
	ctx := context.Background()

	created, err := client.CreateCollection(ctx, &CreateCollectionOption{
		Title:     "go-tissue api v1 test",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("empty collection id")
	}

	got, err := client.GetCollection(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("id mismatch")
	}

	updated, err := client.UpdateCollection(ctx, created.ID, &UpdateCollectionOption{
		Title:     "go-tissue api v1 test (updated)",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "go-tissue api v1 test (updated)" {
		t.Errorf("title not updated: %q", updated.Title)
	}

	item, err := client.CreateCollectionItem(ctx, created.ID, &CreateCollectionItemOption{
		Link: "https://example.com",
		Note: "item test",
		Tags: []string{"test"},
	})
	if err != nil {
		t.Fatal(err)
	}

	newNote := "item test (updated)"
	updatedItem, err := client.UpdateCollectionItem(ctx, created.ID, item.ID, &UpdateCollectionItemOption{
		Note: &newNote,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedItem.Note != newNote {
		t.Errorf("note not updated: %q", updatedItem.Note)
	}

	items, err := client.ListCollectionItems(ctx, created.ID, &PageOption{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Error("no items returned")
	}

	if err := client.DeleteCollectionItem(ctx, created.ID, item.ID); err != nil {
		t.Fatal(err)
	}
	if err := client.DeleteCollection(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
}
