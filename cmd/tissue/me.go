package main

import (
	"context"
	"flag"
)

func cmdMe(args []string) {
	fs := flag.NewFlagSet("me", flag.ExitOnError)
	_ = fs.Parse(args)

	cli := buildClient()
	ctx := context.Background()
	switch cli.config.AuthMethod {
	case authMethodToken:
		me, err := cli.api.Me(ctx)
		if err != nil {
			die("%v", err)
		}
		printJSON(me)
	case authMethodAccount:
		me, err := cli.scraping.Me(ctx)
		if err != nil {
			die("%v", err)
		}
		printJSON(me)
	default:
		die("me is not available for method %s", cli.config.AuthMethod)
	}
}
