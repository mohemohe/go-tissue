go-tissue
====

## checkin (Scraping)

```go
package example

import (
    tissue "github.com/mohemohe/go-tissue"
    "log"
    "time"
)

func main() {
    client, _ := tissue.NewClient(&tissue.ClientOption{
        Email:    "user@example.com",
        Password: "dolphin",
    })
    id, _ := client.CheckIn(&tissue.CheckInOption{
        DateTime:     time.Now(),
        Tags:         []string{"test", "shibafu528"},
        Link:         "https://example.com",
        Note:         "golangでチェックインしたい人生だった",
        Private:      true,
        TooSensitive: false,
    })
    log.Println("checkin id:", id)
}
```

## checkin (WebHook)

```go
package example

import (
    "context"
    tissue "github.com/mohemohe/go-tissue/api"
    "time"
    "log"
)

func main() {
    client, _ := tissue.NewClient(&tissue.ClientOption{
        WebhookID: "dolphin",
    })
    result, _ := client.CheckIn(context.TODO(), &tissue.CheckInOption{
        DateTime:     time.Now(),
        Tags:         []string{"test", "shibafu528"},
        Link:         "https://example.com",
        Note:         "golangでチェックインしたい人生だった",
        Private:      true,
        TooSensitive: false,
    })
    log.Println("checkin id:", result.ID)
}
```

## 免責

しばふに怒られても責任はとれません