package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	authMethodToken   = "token"
	authMethodAccount = "account"
)

type Config struct {
	BaseURL     string `json:"base_url,omitempty"`
	AuthMethod  string `json:"auth_method"`
	AccessToken string `json:"access_token,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tissue", "config.json"), nil
}

func loadConfig() (*Config, error) {
	p, err := configPath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func saveConfig(cfg *Config) error {
	p, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0600)
}

func mustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		die("config not loaded (run `tissue configure` first): %v", err)
	}
	return cfg
}
