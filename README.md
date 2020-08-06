go-tissue
====

## import

```go
import tissue github.com/mohemohe/go-tissue
```

## checkin

```go
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
```

## 免責

しばふに怒られても責任はとれません