package go_tissue

import (
	"os"
	"testing"
	"time"
)

func TestClient_CheckIn(t *testing.T) {
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
}
