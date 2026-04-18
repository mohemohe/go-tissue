package api

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestClient_CheckIn(t *testing.T) {
	if os.Getenv("TISSUE_WEBHOOK_ID") == "" {
		t.Skip("TISSUE_WEBHOOK_ID not set")
	}
	if os.Getenv("TISSUE_SKIP_CHECKIN_TEST") == "1" {
		t.Skip("skip checkin test")
	}
	defer time.Sleep(2 * time.Second)

	client, err := NewClient(&ClientOption{
		BaseURL:   os.Getenv("TISSUE_BASE_URL"),
		WebhookID: os.Getenv("TISSUE_WEBHOOK_ID"),
	})
	if err != nil {
		t.Fatal(err)
	}

	checkIn, err := client.CheckIn(context.TODO(), &CheckInOption{
		DateTime:     time.Now(),
		Tags:         []string{"test", "hoge"},
		Link:         "https://github.com/mohemohe/go-tissue",
		Note:         "go-tissue webhook test checkin",
		Private:      true,
		TooSensitive: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if checkIn.ID == 0 {
		t.Error("empty checkin id")
	}
	t.Log(checkIn)
}
