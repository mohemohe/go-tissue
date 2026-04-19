package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tissue "github.com/mohemohe/go-tissue"
	"github.com/mohemohe/go-tissue/api"
)

func cmdCheckin(args []string) {
	if len(args) == 0 {
		usageCheckin()
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "add":
		cmdCheckinAdd(rest)
	case "list":
		cmdCheckinList(rest)
	case "get":
		cmdCheckinGet(rest)
	case "update":
		cmdCheckinUpdate(rest)
	case "delete":
		cmdCheckinDelete(rest)
	case "-h", "--help", "help":
		usageCheckin()
	default:
		die("unknown checkin subcommand: %s", sub)
	}
}

func usageCheckin() {
	fmt.Fprintln(os.Stderr, "usage: tissue checkin <subcommand>")
	fmt.Fprintln(os.Stderr, "  add     チェックインを作成")
	fmt.Fprintln(os.Stderr, "  list    ユーザーのチェックイン一覧")
	fmt.Fprintln(os.Stderr, "  get     チェックイン詳細")
	fmt.Fprintln(os.Stderr, "  update  チェックイン更新")
	fmt.Fprintln(os.Stderr, "  delete  チェックイン削除")
}

func cmdCheckinAdd(args []string) {
	fs := flag.NewFlagSet("checkin add", flag.ExitOnError)
	tagList := fs.String("tags", "", "カンマ区切りのタグ")
	link := fs.String("link", "", "オカズリンク")
	note := fs.String("note", "", "ノート")
	private := fs.Bool("private", false, "非公開フラグ")
	sensitive := fs.Bool("sensitive", false, "過激フラグ")
	discard := fs.Bool("discard-elapsed-time", false, "経過時間を記録しない")
	at := fs.String("at", "", "チェックイン日時 (RFC3339)")
	_ = fs.Parse(args)

	var tags []string
	if *tagList != "" {
		for _, t := range strings.Split(*tagList, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	var checkedAt *time.Time
	if *at != "" {
		t, err := time.Parse(time.RFC3339, *at)
		if err != nil {
			die("invalid --at: %v", err)
		}
		checkedAt = &t
	}

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.CreateCheckin(ctx, &api.CreateCheckinOption{
			CheckedInAt:        checkedAt,
			Tags:               tags,
			Link:               *link,
			Note:               *note,
			IsPrivate:          *private,
			IsTooSensitive:     *sensitive,
			DiscardElapsedTime: *discard,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.CreateCheckin(ctx, &tissue.CreateCheckinOption{
			CheckedInAt:        checkedAt,
			Tags:               tags,
			Link:               *link,
			Note:               *note,
			IsPrivate:          *private,
			IsTooSensitive:     *sensitive,
			DiscardElapsedTime: *discard,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("unknown method: %s", cli.config.AuthMethod)
	}
}

func cmdCheckinList(args []string) {
	fs := flag.NewFlagSet("checkin list", flag.ExitOnError)
	page := fs.Int("page", 1, "ページ")
	perPage := fs.Int("per-page", 20, "1ページ当たり件数")
	_ = fs.Parse(args)

	cli := buildClient()
	ctx := context.Background()
	name := cli.meName(ctx)

	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.UserCheckins(ctx, name, &api.UserCheckinsOption{
			Page:    *page,
			PerPage: *perPage,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.UserCheckins(ctx, name, &tissue.UserCheckinsOption{
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

func cmdCheckinGet(args []string) {
	fs := flag.NewFlagSet("checkin get", flag.ExitOnError)
	setUsage(fs, "tissue checkin get <id>")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue checkin get <id>")
	}
	id, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid id: %v", err)
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.GetCheckin(ctx, id)
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.GetCheckin(ctx, id)
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("get is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCheckinUpdate(args []string) {
	fs := flag.NewFlagSet("checkin update", flag.ExitOnError)
	setUsage(fs, "tissue checkin update <id> [options]")
	note := fs.String("note", "", "ノート")
	link := fs.String("link", "", "オカズリンク")
	tagList := fs.String("tags", "", "カンマ区切りのタグ (空文字でクリア)")
	setTags := fs.Bool("set-tags", false, "--tags の値でタグを上書きする")
	private := fs.String("private", "", "非公開フラグ (true/false)")
	sensitive := fs.String("sensitive", "", "過激フラグ (true/false)")
	discard := fs.String("discard-elapsed-time", "", "経過時間を記録しない (true/false)")
	at := fs.String("at", "", "チェックイン日時 (RFC3339)")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue checkin update <id> [options]")
	}
	id, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid id: %v", err)
	}

	var notePtr, linkPtr *string
	var tagsPtr *[]string
	var privatePtr, sensitivePtr, discardPtr *bool
	var atPtr *time.Time

	if *note != "" {
		notePtr = note
	}
	if *link != "" {
		linkPtr = link
	}
	if *setTags {
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
	if b, ok := parseOptionalBool(*private); ok {
		privatePtr = &b
	}
	if b, ok := parseOptionalBool(*sensitive); ok {
		sensitivePtr = &b
	}
	if b, ok := parseOptionalBool(*discard); ok {
		discardPtr = &b
	}
	if *at != "" {
		t, err := time.Parse(time.RFC3339, *at)
		if err != nil {
			die("invalid --at: %v", err)
		}
		atPtr = &t
	}

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.UpdateCheckin(ctx, id, &api.UpdateCheckinOption{
			CheckedInAt:        atPtr,
			Tags:               tagsPtr,
			Link:               linkPtr,
			Note:               notePtr,
			IsPrivate:          privatePtr,
			IsTooSensitive:     sensitivePtr,
			DiscardElapsedTime: discardPtr,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.UpdateCheckin(ctx, id, &tissue.UpdateCheckinOption{
			CheckedInAt:        atPtr,
			Tags:               tagsPtr,
			Link:               linkPtr,
			Note:               notePtr,
			IsPrivate:          privatePtr,
			IsTooSensitive:     sensitivePtr,
			DiscardElapsedTime: discardPtr,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("update is not available for method %s", cli.config.AuthMethod)
	}
}

func cmdCheckinDelete(args []string) {
	fs := flag.NewFlagSet("checkin delete", flag.ExitOnError)
	setUsage(fs, "tissue checkin delete <id>")
	pos := parseMixed(fs, args)
	if len(pos) < 1 {
		die("usage: tissue checkin delete <id>")
	}
	id, err := strconv.ParseInt(pos[0], 10, 64)
	if err != nil {
		die("invalid id: %v", err)
	}
	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		if err := cli.api.DeleteCheckin(ctx, id); err != nil {
			die("%v", err)
		}
	case authMethodAccount:
		if err := cli.scraping.DeleteCheckin(ctx, id); err != nil {
			die("%v", err)
		}
	default:
		die("delete is not available for method %s", cli.config.AuthMethod)
	}
	fmt.Fprintln(os.Stderr, "deleted.")
}

func parseOptionalBool(s string) (bool, bool) {
	if s == "" {
		return false, false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		die("invalid bool %q: %v", s, err)
	}
	return b, true
}
