package trv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

type DB struct {
	Tables []Table `json:"tables"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
}

// If there is DB data locally, load it and return it.
func (d *DB) loadData(repo, path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("loadData fail:%s", err)
	}
	raw, _ := ioutil.ReadFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, strings.Replace(path, "/", "-", -1)))
	json.Unmarshal(raw, &d.Tables)
}

// Store DB data locally.
func (d *DB) saveData(repo, path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("saveData fail:%s", err)
	}
	if f, err := os.Stat(fmt.Sprintf("%s/.trv", home)); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(fmt.Sprintf("%s/.trv", home), 0777); err != nil {
			log.Printf("saveData fail:%s", err)
		}
	}
	file, err := json.MarshalIndent(d.Tables, "", " ")
	if err != nil {
		log.Printf("saveData fail:%s", err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, strings.Replace(path, "/", "-", -1)), file, 0644); err != nil {
		log.Printf("saveData fail:%s", err)
	}
}

func (d *DB) fetchDBInfo(client *github.Client, ctx context.Context, source Source) error {
	content, _, _, _ := client.Repositories.GetContents(ctx, source.Owner, source.Repo, fmt.Sprintf("%s/schema.json", source.Path), nil)
	if content != nil {
		text, err := content.GetContent()
		if err != nil {
			return fmt.Errorf("fetch table info fail:%w", err)
		}
		var table DB
		err = json.Unmarshal([]byte(text), &table)
		if err != nil {
			return fmt.Errorf("fetch table info fail:%w", err)
		}
		d.Tables = table.Tables
	} else {
		if len(d.Tables) != 0 {
			return nil
		}
		// Processing in the absence of schema.json
		_, contents, _, err := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
		if err != nil {
			return fmt.Errorf("fech DB info fail:%w", err)
		}
		for _, v := range contents {
			path := v.GetPath()
			if strings.Contains(path, ".md") {
				if strings.Contains(path, "README.md") {
					continue
				}
				var table Table
				if err := table.fetchTableInfoFromMarkdown(client, ctx, source.Owner, source.Repo, path); err != nil {
					return fmt.Errorf("fech DB info fail:%w", err)
				}
				d.Tables = append(d.Tables, table)
			}
		}
	}
	return nil
}
