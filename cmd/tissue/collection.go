package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tissue "github.com/mohemohe/go-tissue"
	"github.com/mohemohe/go-tissue/api"
)

func cmdCollection(args []string) {
	if len(args) == 0 {
		usageCollection()
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "list":
		cmdCollectionList(rest)
	case "create":
		cmdCollectionCreate(rest)
	case "update":
		cmdCollectionUpdate(rest)
	case "delete":
		cmdCollectionDelete(rest)
	case "item":
		cmdCollectionItem(rest)
	case "-h", "--help", "help":
		usageCollection()
	default:
		die("unknown collection subcommand: %s", sub)
	}
}

func usageCollection() {
	fmt.Fprintln(os.Stderr, "usage: tissue collection <subcommand>")
	fmt.Fprintln(os.Stderr, "  list    コレクション一覧")
	fmt.Fprintln(os.Stderr, "  create  コレクション作成")
	fmt.Fprintln(os.Stderr, "  update  コレクション更新")
	fmt.Fprintln(os.Stderr, "  delete  コレクション削除")
	fmt.Fprintln(os.Stderr, "  item    コレクションアイテム操作 (list/add/update/delete)")
}

func cmdCollectionList(args []string) {
	fs := flag.NewFlagSet("collection list", flag.ExitOnError)
	page := fs.Int("page", 1, "ページ")
	perPage := fs.Int("per-page", 20, "1ページ当たり件数")
	_ = fs.Parse(args)

	cli := buildClient()
	ctx := context.Background()

	switch cli.config.AuthMethod {
	case authMethodToken:
		name := cli.meName(ctx)
		result, err := cli.api.UserCollections(ctx, name, &api.PageOption{Page: *page, PerPage: *perPage})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.ListCollections(ctx, &tissue.ListCollectionsOption{
			Page:    *page,
			PerPage: *perPage,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("list is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionCreate(args []string) {
	fs := flag.NewFlagSet("collection create", flag.ExitOnError)
	title := fs.String("title", "", "タイトル")
	private := fs.Bool("private", false, "非公開フラグ")
	_ = fs.Parse(args)
	if *title == "" {
		die("--title is required")
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.CreateCollection(ctx, &api.CreateCollectionOption{
			Title:     *title,
			IsPrivate: *private,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.CreateCollection(ctx, &tissue.CreateCollectionOption{
			Title:     *title,
			IsPrivate: *private,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("create is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionUpdate(args []string) {
	fs := flag.NewFlagSet("collection update", flag.ExitOnError)
	setUsage(fs, "tissue collection update <id> [options]")
	title := fs.String("title", "", "タイトル")
	private := fs.Bool("private", false, "非公開フラグ")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue collection update <id> [options]")
	}
	id, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid id: %v", err)
	}
	if *title == "" {
		die("--title is required")
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.UpdateCollection(ctx, id, &api.UpdateCollectionOption{
			Title:     *title,
			IsPrivate: *private,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.UpdateCollection(ctx, &tissue.UpdateCollectionOption{
			ID:        id,
			Title:     *title,
			IsPrivate: *private,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("update is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionDelete(args []string) {
	fs := flag.NewFlagSet("collection delete", flag.ExitOnError)
	setUsage(fs, "tissue collection delete <id>")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue collection delete <id>")
	}
	id, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid id: %v", err)
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		if err := cli.api.DeleteCollection(ctx, id); err != nil {
			die("%v", err)
		}
	case authMethodAccount:
		if err := cli.scraping.DeleteCollection(ctx, id); err != nil {
			die("%v", err)
		}
	default:
		die("delete is not available for method %s", cli.config.AuthMethod)
	}
	fmt.Fprintln(os.Stderr, "deleted.")
}

func cmdCollectionItem(args []string) {
	if len(args) == 0 {
		usageCollectionItem()
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "list":
		cmdCollectionItemList(rest)
	case "add":
		cmdCollectionItemAdd(rest)
	case "update":
		cmdCollectionItemUpdate(rest)
	case "delete":
		cmdCollectionItemDelete(rest)
	case "-h", "--help", "help":
		usageCollectionItem()
	default:
		die("unknown item subcommand: %s", sub)
	}
}

func usageCollectionItem() {
	fmt.Fprintln(os.Stderr, "usage: tissue collection item <subcommand>")
	fmt.Fprintln(os.Stderr, "  list    コレクション内アイテム一覧")
	fmt.Fprintln(os.Stderr, "  add     アイテムを追加")
	fmt.Fprintln(os.Stderr, "  update  アイテムを更新")
	fmt.Fprintln(os.Stderr, "  delete  アイテムを削除")
}

func cmdCollectionItemList(args []string) {
	fs := flag.NewFlagSet("collection item list", flag.ExitOnError)
	setUsage(fs, "tissue collection item list <collection-id> [options]")
	page := fs.Int("page", 1, "ページ")
	perPage := fs.Int("per-page", 20, "1ページ当たり件数")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue collection item list <collection-id>")
	}
	cid, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid collection id: %v", err)
	}

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.ListCollectionItems(ctx, cid, &api.PageOption{Page: *page, PerPage: *perPage})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.ListCollectionItems(ctx, &tissue.ListCollectionItemsOption{
			CollectionID: cid,
			Page:         *page,
			PerPage:      *perPage,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("item list is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionItemAdd(args []string) {
	fs := flag.NewFlagSet("collection item add", flag.ExitOnError)
	setUsage(fs, "tissue collection item add <collection-id> --link <url> [options]")
	link := fs.String("link", "", "オカズリンク (必須)")
	note := fs.String("note", "", "ノート")
	tagList := fs.String("tags", "", "カンマ区切りのタグ")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue collection item add <collection-id> --link <url> [--note <note>] [--tags a,b]")
	}
	cid, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid collection id: %v", err)
	}
	if *link == "" {
		die("--link is required")
	}
	var tags []string
	if *tagList != "" {
		for _, t := range strings.Split(*tagList, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.CreateCollectionItem(ctx, cid, &api.CreateCollectionItemOption{
			Link: *link,
			Note: *note,
			Tags: tags,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.CreateCollectionItem(ctx, &tissue.CreateCollectionItemOption{
			CollectionID: cid,
			Link:         *link,
			Note:         *note,
			Tags:         tags,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("item add is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionItemUpdate(args []string) {
	fs := flag.NewFlagSet("collection item update", flag.ExitOnError)
	setUsage(fs, "tissue collection item update <collection-id> <item-id> [options]")
	note := fs.String("note", "", "ノート")
	noteSet := fs.Bool("set-note", false, "--note の値でノートを上書き (空値許可)")
	tagList := fs.String("tags", "", "カンマ区切りのタグ")
	tagsSet := fs.Bool("set-tags", false, "--tags の値でタグを上書き (空値許可)")
	pos := parseMixed(fs, args)
	if len(pos) < 2 {
		die("usage: tissue collection item update <collection-id> <item-id> [options]")
	}
	cid, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid collection id: %v", err)
	}
	iid, err := strconv.ParseInt(pos[1], 10, 64)
	if err != nil {
		die("invalid item id: %v", err)
	}

	cli := buildClient()
	ctx := context.Background()

	var notePtr *string
	var tagsPtr *[]string
	if *noteSet {
		notePtr = note
	}
	if *tagsSet {
		var tags []string
		if *tagList != "" {
			for _, t := range strings.Split(*tagList, ",") {
				if trimmed := strings.TrimSpace(t); trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}
		tagsPtr = &tags
	}

	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.UpdateCollectionItem(ctx, cid, iid, &api.UpdateCollectionItemOption{
			Note: notePtr,
			Tags: tagsPtr,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.UpdateCollectionItem(ctx, &tissue.UpdateCollectionItemOption{
			CollectionID: cid,
			ItemID:       iid,
			Note:         notePtr,
			Tags:         tagsPtr,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("item update is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCollectionItemDelete(args []string) {
	fs := flag.NewFlagSet("collection item delete", flag.ExitOnError)
	setUsage(fs, "tissue collection item delete <collection-id> <item-id>")
	pos := parseMixed(fs, args)
	if len(pos) < 2 {
		die("usage: tissue collection item delete <collection-id> <item-id>")
	}
	cid, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid collection id: %v", err)
	}
	iid, err := strconv.ParseInt(pos[1], 10, 64)
	if err != nil {
		die("invalid item id: %v", err)
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		if err := cli.api.DeleteCollectionItem(ctx, cid, iid); err != nil {
			die("%v", err)
		}
	case authMethodAccount:
		if err := cli.scraping.DeleteCollectionItem(ctx, cid, iid); err != nil {
			die("%v", err)
		}
	default:
		die("item delete is not available for method %s", cli.config.AuthMethod)
	}
	fmt.Fprintln(os.Stderr, "deleted.")
}
