package trv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Source []Source `json:"source"`
}

// Check if there is a config file
func (c Config) Exists() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	if f, err := os.Stat(fmt.Sprintf("%s/.trv/config.json", home)); os.IsNotExist(err) || f.IsDir() {
		return false, nil
	} else {
		return true, nil
	}
}

// Load the config file
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

// Retrieve the source list from the config.
func (c Config) getSourceList() []string {
	var sourceList []string
	for _, v := range c.Source {
		sourceList = append(sourceList, fmt.Sprintf("%s/%s", v.Repo, v.Path))
	}
	return sourceList
}

// Add the source to the config.
func (c *Config) addSource(s Source) {
	c.Source = append(c.Source, s)
}

// save the config.
func (c Config) saveConfig() {
	home, _ := os.UserHomeDir()
	file, _ := json.MarshalIndent(c, "", " ")
	dir := filepath.Join(home, ".trv")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0755)
	}
	_ = ioutil.WriteFile(filepath.Join(dir, "config.json"), file, 0644)
}

// Generate GitHub client
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
