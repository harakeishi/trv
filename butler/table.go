package butler

type column struct {
	name    string
	Type    string
	defaul  bool
	comment string
}
type table struct {
	name        string
	description string
	columns     []column
}

// return table_name.column_name
func (t table) getFullName(i int) string {
	return t.name + "." + t.columns[i].name
}
