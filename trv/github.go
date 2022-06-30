package butler

import (
	"strings"
)

func getTableInfo(source Source) []Table {
	var db db
	db.loadData(source.Repo, source.Path)
	client, ctx := source.NewClient()
	if len(db.tables) != 0 {
		return db.tables
	}
	_, contents, _, _ := client.Repositories.GetContents(ctx, source.Owner, source.Repo, source.Path, nil)
	for _, v := range contents {
		path := v.GetPath()
		if strings.Contains(path, ".md") {
			content, _, _, _ := client.Repositories.GetContents(ctx, source.Owner, source.Repo, path, nil)
			if strings.Replace(content.GetName(), ".md", "", -1) == "README" {
				continue
			}
			table := Table{Name: strings.Replace(content.GetName(), ".md", "", -1)}
			text, _ := content.GetContent()
			table.Description = GetDescriptionFromMarkdown(text)
			table.Columns = MarkdownParseTocolumn(text)
			db.tables = append(db.tables, table)
		}
	}
	db.saveData(source.Repo, source.Path)
	return db.tables
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
	for i, v := range rows {
		if i < 4 {
			continue
		}

		colum := strings.Split(v, "|")
		if len(colum) < 8 {
			continue
		}
		result = append(result, Column{
			Name:    strings.TrimSpace(colum[1]),
			Type:    strings.TrimSpace(colum[2]),
			Comment: strings.TrimSpace(colum[8]),
		})
	}
	return result
}
