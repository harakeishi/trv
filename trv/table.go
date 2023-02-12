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

type Repositories interface {
	GetContents(ctx context.Context, owner string, repo string, path string, opt *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

// Get the schema information output in markdown from GitHub.
func (t *Table) fetchTableInfoInMarkdownFromGitHub(repositories Repositories, ctx context.Context, owner, repo, path string) (schema, error) {
	content, _, _, err := repositories.GetContents(ctx, owner, repo, path, nil)
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
