package trv

import "strings"

type schema string

// Parses schema information output in markdown from tbls
func (s schema) parse() (Description string, columns []Column) {
	tmp := strings.Split(string(s), "#")
	d := strings.Split(tmp[3], "\n")
	Description = d[2]

	var columIndex int
	tmp = strings.Split(string(s), "#")
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
			return
		}
		columns = append(columns, Column{
			Name:    strings.TrimSpace(colum[nameIndex]),
			Type:    strings.TrimSpace(colum[typeIndex]),
			Comment: strings.TrimSpace(colum[commentIndex]),
		})
	}
	return
}

func index(a []string, query string) int {
	for i, v := range a {
		if strings.TrimSpace(v) == query {
			return i
		}
	}
	return -1
}
