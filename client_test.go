package go_tissue

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	if os.Getenv("TISSUE_EMAIL") == "" || os.Getenv("TISSUE_PASSWORD") == "" {
		t.Skip("TISSUE_EMAIL / TISSUE_PASSWORD not set")
	}
	client, err := NewClient(&ClientOption{
		BaseURL:  os.Getenv("TISSUE_BASE_URL"),
		Email:    os.Getenv("TISSUE_EMAIL"),
		Password: os.Getenv("TISSUE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestClient_Me(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if me.ID == 0 {
		t.Error("empty user id")
	}
	if me.Name == "" {
		t.Error("empty user name")
	}
	t.Log(me)
}

func TestClient_LatestInformation(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	result, err := client.LatestInformation(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d informations", len(result))
}

func TestClient_DailyCheckinStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	result, err := client.DailyCheckinStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d daily stats", len(result))
}

func TestClient_RecentTags(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	result, err := client.RecentTags(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("recent tags: %v", result)
}

func TestClient_UserCheckins(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserCheckins(context.Background(), me.Name, &UserCheckinsOption{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d checkins for %s", len(result), me.Name)
}

func TestClient_UserTagStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.UserTagStats(context.Background(), me.Name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tag stats: %v", result)
}

func TestClient_UserDailyCheckinStats(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	me, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	until := time.Now()
	since := until.AddDate(-1, 0, 0)
	result, err := client.UserDailyCheckinStats(context.Background(), me.Name, &UserDailyCheckinStatsOption{
		Since: since,
		Until: until,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d daily stats (user)", len(result))
}

func TestClient_SearchCheckins(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	result, err := client.SearchCheckins(context.Background(), &SearchCheckinsOption{
		Query:   "test",
		Page:    1,
		PerPage: 24,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d search results", len(result))
}

func TestClient_CreateCheckin(t *testing.T) {
	if os.Getenv("TISSUE_SKIP_CHECKIN_TEST") == "1" {
		t.Skip("skip checkin test")
	}
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)

	private := true
	result, err := client.CreateCheckin(context.Background(), &CreateCheckinOption{
		Tags:           []string{"test", "hoge"},
		Note:           "本番環境でテストしてすまん",
		IsPrivate:      private,
		IsTooSensitive: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID == 0 {
		t.Error("empty checkin id")
	}
	t.Log(result)
}

func TestClient_CollectionLifecycle(t *testing.T) {
	defer time.Sleep(2 * time.Second)
	client := newTestClient(t)
	ctx := context.Background()

	created, err := client.CreateCollection(ctx, &CreateCollectionOption{
		Title:     "go-tissue test collection",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("empty collection id")
	}

	updated, err := client.UpdateCollection(ctx, &UpdateCollectionOption{
		ID:        created.ID,
		Title:     "go-tissue test collection (updated)",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "go-tissue test collection (updated)" {
		t.Errorf("title not updated: %q", updated.Title)
	}

	list, err := client.ListCollections(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, col := range list {
		if col.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("created collection not found in list")
	}

	item, err := client.CreateCollectionItem(ctx, &CreateCollectionItemOption{
		CollectionID: created.ID,
		Link:         "https://example.com",
		Note:         "test note",
		Tags:         []string{"test"},
	})
	if err != nil {
		t.Fatal(err)
	}

	note := "test note (updated)"
	tags := []string{"test", "updated"}
	updatedItem, err := client.UpdateCollectionItem(ctx, &UpdateCollectionItemOption{
		CollectionID: created.ID,
		ItemID:       item.ID,
		Note:         &note,
		Tags:         &tags,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedItem.Note != note {
		t.Errorf("note not updated: %q", updatedItem.Note)
	}

	items, err := client.ListCollectionItems(ctx, &ListCollectionItemsOption{
		CollectionID: created.ID,
		Page:         1,
		PerPage:      24,
	})
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
