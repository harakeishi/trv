package butler

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Token  string   `json:"token"`
	Source []Source `json:"source"`
}
type Source struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Path  string `json:"path"`
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
		sourceList = append(sourceList, v.Repo)
	}
	return sourceList
}
