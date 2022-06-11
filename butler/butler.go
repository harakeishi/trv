package butler

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var tables []table

func Viewer() {
	config := loadConfig()

	source := config.getSourceList()

	app := tview.NewApplication()

	inputField := tview.NewInputField()
	inputField.SetLabel("serach: ")
	inputField.SetBorder(true)
	textView := tview.NewTextView()
	textView.SetTitle("details")
	textView.SetBorder(true)
	textView.SetText("")
	listView := tview.NewList()
	listView.SetTitle("Result")
	listView.SetBorder(true)
	dropdown := tview.NewDropDown().
		SetLabel("data sorce: ").
		SetOptions(source, func(text string, index int) {
			tables = getTableInfo(config.Token, config.Source[index].Owner, config.Source[index].Repo, config.Source[index].Path)
			listView.Clear()
			filterList(listView, tables, inputField.GetText(), textView)
		})
	dropdown.SetBorder(true)

	grid := tview.NewGrid()
	grid.SetSize(10, 10, 0, 0).
		AddItem(dropdown, 0, 0, 1, 3, 0, 0, true).
		AddItem(inputField, 0, 3, 1, 7, 0, 0, true).
		AddItem(listView, 1, 0, 9, 6, 0, 0, true).
		AddItem(textView, 1, 6, 9, 4, 0, 0, true)

	dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			app.SetFocus(inputField)
			return nil
		}
		return event
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
		listView = filterList(listView, tables, inputField.GetText(), textView)
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

func filterList(list *tview.List, items []table, target string, textView *tview.TextView) *tview.List {
	list.Clear()
	for _, r := range items {
		for i, c := range r.columns {
			if strings.Contains(r.getFullName(i), target) || target == "" {
				list.AddItem(r.getFullName(i), c.comment, 1, func() {})
				list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
					for _, v := range items {
						for a, b := range v.columns {
							if v.getFullName(a) == s1 {
								textView.SetText(fmt.Sprintf("table name: %s\ndetails: %s\n\ncolumn: %s\ntype: %s\ncomment: %s\n", v.name, v.description, b.name, b.Type, b.comment))
							}
						}
					}
				})
			}
		}
	}
	return list
}
