package trv

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

func fetchDbInfo(client *github.Client, ctx context.Context, source Source) []Table {
	var tables []Table
	_, contents, _, _ := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
	for _, v := range contents {
		path := v.GetPath()
		if strings.Contains(path, ".md") {
			if strings.Replace(path, ".md", "", -1) == "README" {
				continue
			}
			table := fetchTableInfo(client, ctx, source.Owner, source.Repo, path)
			tables = append(tables, table)
		}
	}
	return tables
}

func fetchTableInfo(client *github.Client, ctx context.Context, owner, repo, path string) Table {
	var table Table
	content, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	table = Table{Name: strings.Replace(content.GetName(), ".md", "", -1)}
	text, _ := content.GetContent()
	table.Description = GetDescriptionFromMarkdown(text)
	table.Columns = MarkdownParseTocolumn(text)
	return table
}

func GetDescriptionFromMarkdown(text string) string {
	tmp := strings.Split(text, "#")
	d := strings.Split(tmp[3], "\n")
	return d[2]
}
func MarkdownParseTocolumn(text string) []Column {
	var result []Column
	tmp := strings.Split(text, "#")
	rows := strings.Split(tmp[5], "\n")

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
