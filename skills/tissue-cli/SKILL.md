---
name: tissue-cli
description: Use when the user works with the `tissue` CLI (shikorism.net / Tissue). Triggers on requests to check in, list/search checkins, manage collections or collection items, view tag stats, fetch user info, or configure authentication (`tissue configure`, `tissue checkin`, `tissue collection`, `tissue me`, `tissue search`, `tissue tags`). Also applies when discussing the `cmd/tissue` reference CLI in this repository or debugging its behavior.
---

# tissue CLI

`cmd/tissue` は本リポジトリのリファレンス実装 CLI。認証方式は **token** (個人用アクセストークン) / **account** (Email + Password) の2種類。一部サブコマンドは認証方式によって利用可否が異なる。

## インストール

```sh
go install github.com/mohemohe/go-tissue/cmd/tissue@latest
# もしくはリポジトリ内で:
go build -o tissue ./cmd/tissue
```

## 設定 (`tissue configure`)

設定ファイルは `$XDG_CONFIG_HOME/tissue/config.json` (既定 `~/.config/tissue/config.json`) に **0600** で保存される。対話/非対話の双方に対応。

```sh
tissue configure                                                    # 対話モード
tissue configure --method token   --access-token YOUR_TOKEN         # 非対話 (token)
tissue configure --method account --email user@example.com --password ...  # 非対話 (account)
```

個人用アクセストークンは [設定 → 個人用アクセストークン](https://shikorism.net/setting/profile) で発行する。

## サブコマンド早見表

| コマンド | 説明 | 対応認証 |
| --- | --- | --- |
| `tissue me` | 自分のユーザー情報 | token / account |
| `tissue checkin add` | チェックイン作成 | token / account |
| `tissue checkin list` | チェックイン一覧 | token / account |
| `tissue checkin get <id>` | チェックイン詳細 | **token のみ** |
| `tissue checkin update <id>` | チェックイン更新 | **token のみ** |
| `tissue checkin delete <id>` | チェックイン削除 | **token のみ** |
| `tissue collection list` | コレクション一覧 | token / account |
| `tissue collection create` | コレクション作成 | token / account |
| `tissue collection update <id>` | コレクション更新 | token / account |
| `tissue collection delete <id>` | コレクション削除 | token / account |
| `tissue collection item list <cid>` | アイテム一覧 | token / account |
| `tissue collection item add <cid>` | アイテム追加 | token / account |
| `tissue collection item update <cid> <iid>` | アイテム更新 | token / account |
| `tissue collection item delete <cid> <iid>` | アイテム削除 | token / account |
| `tissue search "<query>"` | チェックイン検索 | token / account |
| `tissue tags` | 最近使用タグ | **account のみ** |

## よく使うレシピ

### チェックインする

```sh
# タグ・リンク・ノート・プライバシー設定を付けて作成
tissue checkin add \
  --tags a,b \
  --link https://example.com \
  --note "memo" \
  --private
```

### 一覧・検索

```sh
tissue checkin list --user someone --page 1
tissue search "test"
```

### チェックインを編集/削除 (token 認証のみ)

```sh
tissue checkin get 123
tissue checkin update 123 --note "fixed"
tissue checkin delete 123
```

### コレクション操作

```sh
tissue collection list
tissue collection create --title "title" --private
tissue collection update 47 --title "new" --private
tissue collection delete 47
```

### コレクションアイテム操作

`update` では変更対象フィールドを `--set-*` フラグで明示的に指定する（未指定フィールドは保持）。

```sh
tissue collection item list 47
tissue collection item add 47 --link https://example.com --note memo --tags a,b
tissue collection item update 47 2346 --set-note --note updated --set-tags --tags c
tissue collection item delete 47 2346
```

## 開発・テスト

リポジトリ内で CLI を触る場合:

```sh
task build       # => ./tissue にビルド
task test:cli    # => go test -v -count=1 ./cmd/tissue
```

## トラブルシュート

- **`checkin get/update/delete` が動かない**: account 認証では非対応。`tissue configure --method token ...` でトークン認証に切り替える。
- **`tags` が動かない**: account 認証のみ対応。
- **設定が読めない**: `~/.config/tissue/config.json` が存在してパーミッション 0600 になっているか確認。`$XDG_CONFIG_HOME` が設定されている環境ではそちらが優先される。
- **401 / 認証エラー**: token の失効または Email / Password 変更を疑う。`tissue configure` を再実行。

## 関連リソース

- ライブラリ API: リポジトリルートの [`README.md`](../../README.md)
- REST API 仕様: [`doc/openapi.json`](../../doc/openapi.json)
- scraping 版実装: `client.go` ほかルート `.go` ファイル
- token 版実装: `api/` ディレクトリ
