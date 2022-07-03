package trv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Source []Source `json:"source"`
}

func loadConfig() Config {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func (c Config) getSourceList() []string {
	var sourceList []string
	for _, v := range c.Source {
		sourceList = append(sourceList, fmt.Sprintf("%s/%s", v.Repo, v.Path))
	}
	return sourceList
}

func (s Source) NewClient() (*github.Client, context.Context) {
	var client *github.Client
	var ts oauth2.TokenSource
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.Token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	if s.IsEnterprise {
		client, _ = github.NewEnterpriseClient(s.BaseURL, s.UploadURL, tc)
	} else {
		client = github.NewClient(tc)
	}
	return client, ctx
}
