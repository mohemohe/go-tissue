package main

import (
	"context"
	"flag"

	tissue "github.com/mohemohe/go-tissue"
	"github.com/mohemohe/go-tissue/api"
)

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	setUsage(fs, "tissue search <query> [options]")
	query := fs.String("q", "", "検索キーワード")
	page := fs.Int("page", 1, "ページ")
	perPage := fs.Int("per-page", 20, "1ページ当たり件数")
	pos := parseMixed(fs, args)
	if *query == "" && len(pos) > 0 {
		*query = pos[0]
	}
	if *query == "" {
		die("usage: tissue search <query> [--page N] [--per-page N]")
	}

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.SearchCheckins(ctx, &api.SearchOption{
			Query:   *query,
			Page:    *page,
			PerPage: *perPage,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.SearchCheckins(ctx, &tissue.SearchCheckinsOption{
			Query:   *query,
			Page:    *page,
			PerPage: *perPage,
		})
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("search is not available for method %s", cli.config.AuthMethod)
	}
}
