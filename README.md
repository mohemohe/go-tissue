go-tissue
====

[shikorism.net](https://shikorism.net) (Tissue) の非公式 Go クライアントライブラリ。

認証方式は以下の2種類:

- **スクレイピング版** (`github.com/mohemohe/go-tissue`): Email / Password でセッション Cookie を取得する方式
- **API トークン版** (`github.com/mohemohe/go-tissue/api`): 個人用アクセストークンで公開 REST API (`/api/v1/...`) にアクセスする方式。Webhook によるチェックインも同パッケージで提供

## checkin (スクレイピング版)

```go
package example

import (
    "context"
    "log"

    tissue "github.com/mohemohe/go-tissue"
)

func main() {
    client, _ := tissue.NewClient(&tissue.ClientOption{
        Email:    "user@example.com",
        Password: "dolphin",
    })

    result, err := client.CreateCheckin(context.Background(), &tissue.CreateCheckinOption{
        Tags:           []string{"test", "shibafu528"},
        Link:           "https://example.com",
        Note:           "golangでチェックインしたい人生だった",
        IsPrivate:      true,
        IsTooSensitive: false,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Println("checkin id:", result.ID)
}
```

## checkin (API トークン版)

個人用アクセストークンは [設定 → 個人用アクセストークン](https://shikorism.net/setting/profile) で発行できます。

```go
package example

import (
    "context"
    "log"

    tissue "github.com/mohemohe/go-tissue/api"
)

func main() {
    client, _ := tissue.NewClient(&tissue.ClientOption{
        AccessToken: "YOUR_PERSONAL_ACCESS_TOKEN",
    })

    result, err := client.CreateCheckin(context.Background(), &tissue.CreateCheckinOption{
        Tags:           []string{"test", "shibafu528"},
        Link:           "https://example.com",
        Note:           "golangでチェックインしたい人生だった",
        IsPrivate:      true,
        IsTooSensitive: false,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Println("checkin id:", result.ID)
}
```

## checkin (Webhook)

```go
package example

import (
    "context"
    "log"
    "time"

    tissue "github.com/mohemohe/go-tissue/api"
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

## 提供 API

### スクレイピング版 (`go-tissue`)

初回呼び出し時に `GET /login` → `POST /login` で自動的にセッションを確立する。

- `Me(ctx)` — 自分のユーザー情報とチェックイン概況
- `LatestInformation(ctx)` — サイトのお知らせ一覧
- `DailyCheckinStats(ctx)` — サイト全体の日次統計
- `UserDailyCheckinStats(ctx, user, option)` — ユーザー単位の日次統計 (since/until)
- `UserTagStats(ctx, user)` — ユーザーのタグ使用統計
- `UserCheckins(ctx, user, option)` — ユーザーのチェックイン一覧 (page / per_page / has_link)
- `SearchCheckins(ctx, option)` — チェックイン検索
- `RecentTags(ctx)` — 最近使用したタグ
- `CreateCheckin(ctx, option)` — チェックインの作成
- `ListCollections(ctx, option)`, `CreateCollection(ctx, option)`, `UpdateCollection(ctx, option)`, `DeleteCollection(ctx, id)` — コレクション操作 (page / per_page)
- `ListCollectionItems(ctx, option)`, `CreateCollectionItem(ctx, option)`, `UpdateCollectionItem(ctx, option)`, `DeleteCollectionItem(ctx, collectionID, itemID)` — コレクションアイテム操作

### API トークン版 (`go-tissue/api`)

[`doc/openapi.json`](./doc/openapi.json) の `/api/v1/...` エンドポイントをカバーする。全メソッドは `Authorization: Bearer <AccessToken>` ヘッダーで認証する。

- `Me(ctx)` — 自分のユーザー情報
- `GetUser(ctx, name)` — ユーザー情報
- `UserCheckins(ctx, name, option)` — ユーザーのチェックイン一覧 (page / per_page / has_link / since / until / order)
- `UserLikes(ctx, name, option)` — ユーザーがいいねしたチェックイン
- `UserCollections(ctx, name, option)` — ユーザーのコレクション
- `UserDailyCheckinStats(ctx, name, option)`, `UserHourlyCheckinStats(ctx, name, option)`, `UserTagStats(ctx, name, option)`, `UserLinkStats(ctx, name, option)` — ユーザー統計
- `CreateCheckin(ctx, option)`, `GetCheckin(ctx, id)`, `UpdateCheckin(ctx, id, option)`, `DeleteCheckin(ctx, id)` — チェックイン CRUD
- `CreateCollection(ctx, option)`, `GetCollection(ctx, id)`, `UpdateCollection(ctx, id, option)`, `DeleteCollection(ctx, id)` — コレクション CRUD
- `ListCollectionItems(ctx, collectionID, option)`, `CreateCollectionItem(ctx, collectionID, option)`, `UpdateCollectionItem(ctx, collectionID, itemID, option)`, `DeleteCollectionItem(ctx, collectionID, itemID)` — アイテム CRUD
- `SearchCheckins(ctx, option)`, `SearchCollections(ctx, option)` — 検索
- `CheckIn(ctx, option)` — Webhook 経由のチェックイン (WebhookID が必要)

## CLI (`cmd/tissue`)

リファレンス実装の CLI。認証方式は `token` (個人用アクセストークン) / `account` (Email + Password) の2種類。

### インストール

```sh
go install github.com/mohemohe/go-tissue/cmd/tissue@latest
```

### 初期設定

```sh
tissue configure
# または非対話:
tissue configure --method token --access-token YOUR_TOKEN
tissue configure --method account --email user@example.com --password ...
```

設定ファイルは `$XDG_CONFIG_HOME/tissue/config.json` (既定 `~/.config/tissue/config.json`) にパーミッション 0600 で保存される。

### 主要コマンド

```sh
tissue me                                              # 自分のユーザー情報
tissue checkin add --tags a,b --note "memo" --private  # チェックイン
tissue checkin list --user someone --page 1            # チェックイン一覧
tissue checkin get 123                                 # (token) 詳細
tissue checkin update 123 --note "fixed"               # (token) 更新
tissue checkin delete 123                              # (token) 削除

tissue collection list
tissue collection create --title "title" --private
tissue collection update 47 --title "new" --private
tissue collection delete 47

tissue collection item list 47
tissue collection item add 47 --link https://example.com --note memo --tags a,b
tissue collection item update 47 2346 --set-note --note updated --set-tags --tags c
tissue collection item delete 47 2346

tissue search "test"
tissue tags                                            # (account のみ)
```

一部のコマンドは認証方式によって制限がある (例: `checkin get/update/delete` は token 認証のみ)。

## 免責

しばふに怒られても責任はとれません
