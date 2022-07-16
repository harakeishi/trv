package trv

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const addSource = "add source"

type Trv struct {
	Config         Config
	Source         []string
	DB             Db
	Tables         []Table
	SourceSelecter *tview.DropDown
	TableViewer    *tview.List
	Searcher       *tview.InputField
	Pages          *tview.Pages
	InfoLayout     *tview.Grid
	InfoText       *tview.TextView
	InfoTable      *tview.Table
	App            *tview.Application
	Layout         *tview.Grid
}

func (t *Trv) Init() {
	//set data
	t.setConfig()
	t.setSource()
	// gui setting
	t.App = tview.NewApplication()
	t.App.EnableMouse(true)
	t.setSourceSelecter()
	t.setSearcher()
	t.setTableViewer()
	t.setInfoText()
	t.setInfoTable()
	t.setInfoLayout()
	t.setLayout()

	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlD:
			t.App.SetFocus(t.SourceSelecter)
		case tcell.KeyCtrlS:
			t.App.SetFocus(t.Searcher)
		case tcell.KeyCtrlR:
			t.App.SetFocus(t.TableViewer)
		}
		return event
	})
}

// Set Config
func (t *Trv) setConfig() {
	t.Config = loadConfig()
}

// Set Source
func (t *Trv) setSource() {
	t.Source = t.Config.getSourceList()
}

// set Source Selecter
func (t *Trv) setSourceSelecter() {
	t.SourceSelecter = tview.NewDropDown()
	t.SourceSelecter.SetLabel("data source: ").
		SetOptions(t.Source, func(text string, index int) {
			t.DB = t.Config.Source[index].setDbData()
			t.Tables = t.DB.tables
			t.TableViewer.Clear()
			t.filterList()
			t.App.SetFocus(t.Searcher)
		})
	t.SourceSelecter.SetBorder(true)
}

func (t *Trv) setTableViewer() {
	t.TableViewer = tview.NewList()
	t.TableViewer.SetTitle("Result")
	t.TableViewer.SetBorder(true)

	t.TableViewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			t.App.SetFocus(t.Searcher)
			return nil
		}
		return event
	})
}

func (t *Trv) setSearcher() {
	t.Searcher = tview.NewInputField()
	t.Searcher.SetLabel("serach:")
	t.Searcher.SetBorder(true)
	t.Searcher.SetChangedFunc(func(text string) {
		t.filterList()
	})
}

// set Source Info Text
func (t *Trv) setInfoText() {
	t.InfoText = tview.NewTextView()
	t.InfoText.SetText("")
}

// set Source Info Table
func (t *Trv) setInfoTable() {
	t.InfoTable = tview.NewTable().
		SetBorders(true).
		SetCell(0, 0, tview.NewTableCell("column").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(0, 1, tview.NewTableCell("type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(0, 2, tview.NewTableCell("comment").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)).
		SetCell(1, 0, tview.NewTableCell("")).
		SetCell(1, 1, tview.NewTableCell("")).
		SetCell(1, 2, tview.NewTableCell(""))
}

func (t *Trv) setPages() {
	t.Pages = tview.NewPages()
}

func (t *Trv) setInfoLayout() {
	t.InfoLayout = tview.NewGrid()
	t.InfoLayout.SetTitle("details").SetBorder(true)
	t.InfoLayout.SetSize(5, 5, 0, 0).
		AddItem(t.InfoText, 0, 0, 2, 5, 0, 0, true).
		AddItem(t.InfoTable, 2, 0, 3, 5, 2, 5, true)
	t.InfoLayout.SetOffset(1, 1)
}
func (t *Trv) setLayout() {
	t.Layout = tview.NewGrid()
	t.Layout.SetSize(10, 10, 0, 0).
		AddItem(t.SourceSelecter, 0, 0, 2, 3, 0, 0, true).
		AddItem(t.Searcher, 0, 3, 2, 7, 0, 0, true).
		AddItem(t.TableViewer, 2, 0, 8, 5, 0, 0, true).
		AddItem(t.InfoLayout, 2, 5, 8, 5, 0, 0, true)
}

//filterList(list *tview.List, items []Table, target string, textView *tview.TextView, table *tview.Table) *tview.List {
func (t *Trv) filterList() {
	target := t.Searcher.GetText()
	t.TableViewer.Clear()
	for _, r := range t.Tables {
		for i, c := range r.Columns {
			if strings.Contains(strings.ToLower(r.getFullName(i)), strings.ToLower(target)) || target == "" {
				t.TableViewer.AddItem(r.getFullName(i), c.Comment, 1, func() {})
				t.TableViewer.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
					for _, v := range t.Tables {
						for a, b := range v.Columns {
							if v.getFullName(a) == s1 {
								t.InfoText.SetText(fmt.Sprintf("table name: %s\ndetails: %s", v.Name, v.Description))
								t.InfoTable.RemoveRow(1)
								t.InfoTable.SetCell(1, 0, tview.NewTableCell(b.Name))
								t.InfoTable.SetCell(1, 1, tview.NewTableCell(b.Type))
								t.InfoTable.SetCell(1, 2, tview.NewTableCell(b.Comment))
							}
						}
					}
				})
			}
		}
	}
}

func (t Trv) Draw() {
	if err := t.App.SetRoot(t.Layout, true).SetFocus(t.SourceSelecter).Run(); err != nil {
		panic(err)
	}
}
