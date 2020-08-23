package api

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestClient_CheckIn(t *testing.T) {
	client, err := NewClient(&ClientOption{
		BaseURL: os.Getenv("TISSUE_BASE_URL"),
		WebhookID: os.Getenv("TISSUE_WEBHOOK_ID"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}

	checkIn, err := client.CheckIn(context.TODO(), &CheckInOption{
		DateTime:     time.Now(),
		Tags:         []string{"test", "hoge"},
		Link:         "https://github.com/mohemohe/go-tissue",
		Note:         "go-tissue test checkin",
		Private:      true,
		TooSensitive: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if checkIn == nil {
		t.Error("something wrong")
	}
	if checkIn.ID == 0 {
		t.Error("something wrong")
	}
}
