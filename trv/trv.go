package trv

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateSerachField() *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("serach:")
	inputField.SetBorder(true)
	return inputField
}

func Viewer() {
	var tables []Table
	config := loadConfig()

	source := config.getSourceList()

	app := tview.NewApplication()
	app.EnableMouse(true)
	inputField := CreateSerachField()
	textView := tview.NewTextView()
	textView.SetText("")
	listView := tview.NewList()
	listView.SetTitle("Result")
	listView.SetBorder(true)
	table := tview.NewTable().
		SetBorders(true)
	table.SetCell(0, 0, tview.NewTableCell("column").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("comment").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(1, 0, tview.NewTableCell(""))
	table.SetCell(1, 1, tview.NewTableCell(""))
	table.SetCell(1, 2, tview.NewTableCell(""))

	dropdown := tview.NewDropDown().
		SetLabel("data source: ").
		SetOptions(source, func(text string, index int) {
			db := config.Source[index].setDbData()
			tables = db.tables
			listView.Clear()
			filterList(listView, tables, inputField.GetText(), textView, table)
			app.SetFocus(inputField)
		})
	dropdown.SetBorder(true)

	detailsBox := tview.NewGrid()
	detailsBox.SetTitle("details").SetBorder(true)
	detailsBox.SetSize(5, 5, 0, 0).
		AddItem(textView, 0, 0, 2, 5, 0, 0, true).
		AddItem(table, 2, 0, 3, 5, 2, 5, true)
	detailsBox.SetOffset(1, 1)
	grid := tview.NewGrid()
	grid.SetSize(10, 10, 0, 0).
		AddItem(dropdown, 0, 0, 1, 3, 0, 0, true).
		AddItem(inputField, 0, 3, 1, 7, 0, 0, true).
		AddItem(listView, 1, 0, 6, 10, 0, 0, true).
		AddItem(detailsBox, 7, 0, 3, 10, 0, 0, true)

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
		listView = filterList(listView, tables, inputField.GetText(), textView, table)
	})

	listView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			app.SetFocus(inputField)
			return nil
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlD:
			app.SetFocus(dropdown)
		case tcell.KeyCtrlS:
			app.SetFocus(inputField)
		case tcell.KeyCtrlR:
			app.SetFocus(listView)
		}
		return event
	})

	if err := app.SetRoot(grid, true).SetFocus(dropdown).Run(); err != nil {
		panic(err)
	}
}

func filterList(list *tview.List, items []Table, target string, textView *tview.TextView, table *tview.Table) *tview.List {
	list.Clear()
	for _, r := range items {
		for i, c := range r.Columns {
			if strings.Contains(strings.ToLower(r.getFullName(i)), strings.ToLower(target)) || target == "" {
				list.AddItem(r.getFullName(i), c.Comment, 1, func() {})
				list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
					for _, v := range items {
						for a, b := range v.Columns {
							if v.getFullName(a) == s1 {
								textView.SetText(fmt.Sprintf("table name: %s\ndetails: %s", v.Name, v.Description))
								table.RemoveRow(1)
								table.SetCell(1, 0, tview.NewTableCell(b.Name))
								table.SetCell(1, 1, tview.NewTableCell(b.Type))
								table.SetCell(1, 2, tview.NewTableCell(b.Comment))
							}
						}
					}
				})
			}
		}
	}
	return list
}
