package butler

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type column struct {
	name    string
	Type    string
	defaul  bool
	comment string
}
type table struct {
	name    string
	columns []column
}

func Viewer() {
	tables := []table{{name: "user", columns: []column{{name: "name", comment: "氏名"}, {name: "age", comment: "年齢"}}},
		{name: "domain", columns: []column{{name: "name", comment: "ドメイン名"}, {name: "state", comment: "ドメインの状態（0:無効、1:有効）"}}}}
	app := tview.NewApplication()
	dropdown := tview.NewDropDown().
		SetLabel("data sorce: ").
		SetOptions([]string{"sampleDB1", "sampleDB2", "sampleDB3"}, nil)
	dropdown.SetBorder(true)

	listView := tview.NewList()
	listView.SetTitle("Result")
	listView.SetBorder(true)
	for _, r := range tables {
		for _, c := range r.columns {
			listView.AddItem(r.name+"."+c.name, c.comment, 1, nil)
		}
	}

	inputField := tview.NewInputField()
	inputField.SetLabel("serach: ")
	inputField.SetBorder(true)

	grid := tview.NewGrid()
	grid.SetSize(10, 10, 0, 0).
		AddItem(listView, 1, 0, 9, 10, 0, 0, true).
		AddItem(dropdown, 0, 0, 1, 3, 0, 0, true).
		AddItem(inputField, 0, 3, 1, 7, 0, 0, true)

	dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			app.SetFocus(inputField)
			return nil
		}
		return nil
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			app.SetFocus(listView)
			return nil

		}
		return event
	})
	inputField.SetChangedFunc(func(text string) {
		listView = filterList(listView, tables, inputField.GetText())
	})

	listView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			app.SetFocus(inputField)
			return nil
		}
		return event
	})

	if err := app.SetRoot(grid, true).SetFocus(dropdown).Run(); err != nil {
		panic(err)
	}
}

func filterList(list *tview.List, items []table, target string) *tview.List {
	list.Clear()
	for _, r := range items {
		for _, c := range r.columns {
			if strings.Contains(c.name, target) {
				list.AddItem(r.name+"."+c.name, c.comment, 1, nil)
			}
		}
	}
	return list
}
