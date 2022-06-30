package butler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Column struct {
	Name    string
	Type    string
	Defaul  bool
	Comment string
}
type Table struct {
	Name        string
	Description string
	Columns     []Column
	UpdateDate  time.Time
}
type db struct {
	tables []Table
}

// return table_name.column_name
func (t Table) getFullName(i int) string {
	return t.Name + "." + t.Columns[i].Name
}

func (d *db) loadData(repo, path string) {
	home, _ := os.UserHomeDir()
	raw, err := ioutil.ReadFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, path))
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(raw, &d.tables)
}

func (d *db) saveData(repo, path string) {
	home, _ := os.UserHomeDir()
	if f, err := os.Stat(fmt.Sprintf("%s/.trv", home)); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(fmt.Sprintf("%s/.trv", home), 0777); err != nil {
			fmt.Println(err)
		}
	}
	file, _ := json.MarshalIndent(d.tables, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("%s/.trv/%s-%s.json", home, repo, path), file, 0644)
}
