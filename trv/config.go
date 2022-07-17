package trv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Source []Source `json:"source"`
}

func (c *Config) loadConfig() {
	home, _ := os.UserHomeDir()
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/.trv/config.json", home))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &c); err != nil {
		log.Fatal(err)
	}
}

func (c Config) getSourceList() []string {
	var sourceList []string
	for _, v := range c.Source {
		sourceList = append(sourceList, fmt.Sprintf("%s/%s", v.Repo, v.Path))
	}
	return sourceList
}
func (c *Config) addSource(s Source) {
	c.Source = append(c.Source, s)
	c.saveConfig()
}

func (c Config) saveConfig() {
	home, _ := os.UserHomeDir()
	file, _ := json.MarshalIndent(c, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("%s/.trv/config.json", home), file, 0644)
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
