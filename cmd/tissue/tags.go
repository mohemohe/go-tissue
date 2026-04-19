package main

import (
	"context"
	"flag"
)

func cmdTags(args []string) {
	fs := flag.NewFlagSet("tags", flag.ExitOnError)
	_ = fs.Parse(args)

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		result, err := cli.api.RecentTags(ctx)
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	case authMethodAccount:
		result, err := cli.scraping.RecentTags(ctx)
		if err != nil {
			die("%v", err)
		}
		printJSON(result)
	default:
		die("tags is not available for method %s", cli.config.AuthMethod)
	}
}
