package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Context              GithubContext
	Theme                string
	TwitterUsername      string
	FeedURL              string
	ExtraRepo            string
	ExtraRepoDescription string
}

type GithubContext struct {
	Token      string `json:"token"`
	Repository string `json:"repository"`
}

func Derive() (*Config, error) {
	var c Config
	if err := json.Unmarshal([]byte(os.Getenv("GITHUB_CONTEXT")), &c.Context); err != nil {
		c.Context.Repository = os.Getenv("GITHUB_REPOSITORY")
		c.Context.Token = os.Getenv("GITHUB_TOKEN")
	}
	c.Context.Token = readInput("token", c.Context.Token)
	c.Theme = readInput("theme", "dark")
	c.TwitterUsername = readInput("twitter_username", "")
	c.FeedURL = readInput("feed_url", "")
	c.ExtraRepo = readInput("extra_repo", "")
	c.ExtraRepoDescription = readInput("extra_repo_description", "")
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
