package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Context GithubContext
	Theme   string
}

type GithubContext struct {
	Token      string `json:"token"`
	Repository string `json:"repository"`
}

func Derive() (*Config, error) {
	var c Config
	if err := json.Unmarshal([]byte(os.Getenv("GITHUB_CONTEXT")), &c.Context); err != nil {
		return nil, fmt.Errorf("github context is missing or invalid: %s", err)
	}
	c.Theme = readInput("theme", "dark")
	return &c, nil
}

func readInput(name string, def string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = fmt.Sprintf("INPUT_%s", strings.ToUpper(name))
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return def
}
