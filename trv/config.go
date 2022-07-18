package trv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Source []Source `json:"source"`
}

func (c *Config) loadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/.trv/config.json", home))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &c); err != nil {
		return err
	}
	return nil
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

func (s Source) NewClient() (*github.Client, context.Context, error) {
	var client *github.Client
	var ts oauth2.TokenSource
	var err error
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.Token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	if s.IsEnterprise {
		client, err = github.NewEnterpriseClient(s.BaseURL, s.UploadURL, tc)
		if err != nil {
			return nil, nil, fmt.Errorf("new client fail:%w", err)
		}
	} else {
		client = github.NewClient(tc)
	}
	return client, ctx, nil
}
