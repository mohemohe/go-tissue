package main

import (
	"context"

	tissue "github.com/mohemohe/go-tissue"
	"github.com/mohemohe/go-tissue/api"
)

type clientBundle struct {
	config   *Config
	scraping *tissue.Client
	api      *api.Client
}

func buildClient() *clientBundle {
	cfg := mustLoadConfig()
	b := &clientBundle{config: cfg}
	switch cfg.AuthMethod {
	case authMethodToken:
		c, err := api.NewClient(&api.ClientOption{
			BaseURL:     cfg.BaseURL,
			AccessToken: cfg.AccessToken,
		})
		if err != nil {
			die("failed to create api client: %v", err)
		}
		b.api = c
	case authMethodAccount:
		c, err := tissue.NewClient(&tissue.ClientOption{
			BaseURL:  cfg.BaseURL,
			Email:    cfg.Email,
			Password: cfg.Password,
		})
		if err != nil {
			die("failed to create client: %v", err)
		}
		b.scraping = c
	default:
		die("unknown auth_method: %q (run `tissue configure`)", cfg.AuthMethod)
	}
	return b
}

func (b *clientBundle) meName(ctx context.Context) string {
	switch b.config.AuthMethod {
	case authMethodToken:
		me, err := b.api.Me(ctx)
		if err != nil {
			die("%v", err)
		}
		return me.Name
	case authMethodAccount:
		me, err := b.scraping.Me(ctx)
		if err != nil {
			die("%v", err)
		}
		return me.Name
	}
	die("cannot resolve username for method %s", b.config.AuthMethod)
	return ""
}
