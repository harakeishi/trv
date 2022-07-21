package trv

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
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

// return table_name.column_name
func (t Table) getFullName(i int) string {
	return t.Name + "." + t.Columns[i].Name
}

func (t *Table) fetchTableInfo(client *github.Client, ctx context.Context, owner, repo, path string) error {
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return fmt.Errorf("fetch table info fail:%w", err)
	}
	t.Name = strings.Replace(content.GetName(), ".md", "", -1)
	text, err := content.GetContent()
	if err != nil {
		return fmt.Errorf("fetch table info fail:%w", err)
	}
	t.Description = GetDescriptionFromMarkdown(text)
	t.Columns = MarkdownParseTocolumn(text)
	return nil
}

func GetDescriptionFromMarkdown(text string) string {
	tmp := strings.Split(text, "#")
	d := strings.Split(tmp[3], "\n")
	return d[2]
}

func MarkdownParseTocolumn(text string) []Column {
	var result []Column
	var columIndex int
	tmp := strings.Split(text, "#")
	for i, v := range tmp {
		if strings.Contains(v, "Columns") {
			columIndex = i
		}
	}

	rows := strings.Split(tmp[columIndex], "\n")
	header := strings.Split(rows[2], "|")
	nameIndex := index(header, "Name")
	typeIndex := index(header, "Type")
	commentIndex := index(header, "Comment")

	for i, v := range rows {
		if i < 4 {
			continue
		}
		colum := strings.Split(v, "|")
		if len(colum) < 8 {
			return result
		}
		result = append(result, Column{
			Name:    strings.TrimSpace(colum[nameIndex]),
			Type:    strings.TrimSpace(colum[typeIndex]),
			Comment: strings.TrimSpace(colum[commentIndex]),
		})
	}
	return result
}

func index(a []string, query string) int {
	for i, v := range a {
		if strings.TrimSpace(v) == query {
			return i
		}
	}
	return -1
}
