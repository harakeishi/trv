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
	tables []Table
}

// If there is DB data locally, load it and return it.
func (d *DB) loadData(repo, path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("loadData fail:%s", err)
	}
	raw, _ := ioutil.ReadFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, strings.Replace(path, "/", "-", -1)))
	json.Unmarshal(raw, &d.tables)
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
	file, err := json.MarshalIndent(d.tables, "", " ")
	if err != nil {
		log.Printf("saveData fail:%s", err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, strings.Replace(path, "/", "-", -1)), file, 0644); err != nil {
		log.Printf("saveData fail:%s", err)
	}
}

func (d *DB) fetchDBInfo(client *github.Client, ctx context.Context, source Source) error {
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
			if err := table.fetchTableInfo(client, ctx, source.Owner, source.Repo, path); err != nil {
				return fmt.Errorf("fech DB info fail:%w", err)
			}
			d.tables = append(d.tables, table)
		}
	}
	return nil
}
