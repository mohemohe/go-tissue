package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func cmdConfigure(args []string) {
	fs := flag.NewFlagSet("configure", flag.ExitOnError)
	method := fs.String("method", "", "認証方式: token / account")
	baseURL := fs.String("base-url", "", "Tissue base URL (例: https://shikorism.net)")
	accessToken := fs.String("access-token", "", "個人用アクセストークン (method=token)")
	email := fs.String("email", "", "Email (method=account)")
	password := fs.String("password", "", "Password (method=account)")
	noPrompt := fs.Bool("no-prompt", false, "対話プロンプトを抑制し、指定されたフラグのみで保存")
	_ = fs.Parse(args)

	cfg, _ := loadConfig()
	if cfg == nil {
		cfg = &Config{}
	}

	reader := bufio.NewReader(os.Stdin)

	if *method != "" {
		cfg.AuthMethod = *method
	} else if !*noPrompt {
		cfg.AuthMethod = promptWithDefault(reader, "認証方式 (token/account)", defaultOr(cfg.AuthMethod, "token"))
	} else if cfg.AuthMethod == "" {
		cfg.AuthMethod = "token"
	}

	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	} else if !*noPrompt {
		cfg.BaseURL = promptWithDefault(reader, "Base URL", defaultOr(cfg.BaseURL, "https://shikorism.net"))
	}

	switch cfg.AuthMethod {
	case authMethodToken:
		if *accessToken != "" {
			cfg.AccessToken = *accessToken
		} else if !*noPrompt {
			cfg.AccessToken = promptWithDefault(reader, "個人用アクセストークン", cfg.AccessToken)
		}
		cfg.Email = ""
		cfg.Password = ""
	case authMethodAccount:
		if *email != "" {
			cfg.Email = *email
		} else if !*noPrompt {
			cfg.Email = promptWithDefault(reader, "Email", cfg.Email)
		}
		if *password != "" {
			cfg.Password = *password
		} else if !*noPrompt {
			cfg.Password = promptWithDefault(reader, "Password (平文保存されます)", cfg.Password)
		}
		cfg.AccessToken = ""
	default:
		die("unknown auth method: %q (want token/account)", cfg.AuthMethod)
	}

	if err := saveConfig(cfg); err != nil {
		die("failed to save config: %v", err)
	}
	p, _ := configPath()
	fmt.Fprintf(os.Stderr, "saved: %s\n", p)
}

func defaultOr(val, def string) string {
	if val != "" {
		return val
	}
	return def
}

func promptWithDefault(reader *bufio.Reader, label, def string) string {
	if def != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, def)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return def
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}
