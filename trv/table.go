package trv

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Defaul   bool   `json:"default"`
	Comment  string `json:"comment"`
}
type Table struct {
	Name        string   `json:"name"`
	Description string   `json:"comment"`
	Columns     []Column `json:"columns"`
}

// return table_name.column_name
func (t Table) getFullName(i int) string {
	return t.Name + "." + t.Columns[i].Name
}

// Get the schema information output in markdown from GitHub.
func (t *Table) fetchTableInfoInMarkdownFromGitHub(client *github.Client, ctx context.Context, owner, repo, path string) (schema, error) {
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("fetch table info fail:%w", err)
	}
	t.Name = strings.Replace(content.GetName(), ".md", "", -1)
	text, err := content.GetContent()
	if err != nil {
		return "", fmt.Errorf("fetch table info fail:%w", err)
	}

	return schema(text), nil
}
