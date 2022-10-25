package trv

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type Trv struct {
	Config         Config
	Source         []string
	DB             DB
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
	Modal          tview.Primitive
	Form           *tview.Form
}

/*
Prepare to start trv.
Specifically, loading the configuration, loading the data, and preparing the TUI.
*/
func (t *Trv) Init() error {
	os.Setenv("LC_CTYPE", "en_US.UTF-8")
	//set data
	if err := t.setConfig(); err != nil {
		return err
	}
	t.setSource()

	// gui setting
	t.App = tview.NewApplication()
	t.App.EnableMouse(true)
	if err := t.setSourceSelecter(); err != nil {
		return err
	}
	t.setSearcher()
	t.setTableViewer()
	t.setInfoText()
	t.setInfoTable()
	t.setInfoLayout()
	t.setLayout()
	t.setForm()
	t.setModal()
	t.setPages()

	t.Pages.HidePage("modal")

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
	return nil
}

// Set Config
func (t *Trv) setConfig() error {
	if err := t.Config.loadConfig(); err != nil {
		return fmt.Errorf("setConfig fail: %w", err)
	}
	return nil
}

// Set Source
func (t *Trv) setSource() {
	t.Source = t.Config.getSourceList()
}

// set Source Selecter
func (t *Trv) setSourceSelecter() error {
	var err error
	t.SourceSelecter = tview.NewDropDown()
	t.SourceSelecter.SetTitle("data source(Ctrl+d)")
	t.SourceSelecter.SetLabel("data source: ").
		SetOptions(t.Source, func(text string, index int) {
			if index > len(t.Source)-1 {
				return
			}
			t.DB, err = t.Config.Source[index].setDbData()
			t.Tables = t.DB.tables
			t.filterList()
			t.App.SetFocus(t.Searcher)
		})
	if err != nil {
		return fmt.Errorf("set source selecter fail:%w", err)
	}
	t.SourceSelecter.AddOption("add source", func() {
		t.Pages.ShowPage("modal")
	})
	t.SourceSelecter.SetBorder(true)
	return nil
}

// add dropdown option
func (t *Trv) addDropdownOption() error {
	var err error
	currentOptionCount := t.SourceSelecter.GetOptionCount()
	lastOptionIndex := len(t.Source) - 1
	t.SourceSelecter.RemoveOption(currentOptionCount - 1)
	t.SourceSelecter.AddOption(t.Source[lastOptionIndex], func() {
		t.DB, err = t.Config.Source[lastOptionIndex].setDbData()
		t.Tables = t.DB.tables
		t.filterList()
		t.App.SetFocus(t.Searcher)
	})
	if err != nil {
		return fmt.Errorf("add dropdown option fail:%w", err)
	}
	t.SourceSelecter.AddOption("add source", func() {
		t.Pages.ShowPage("modal")
	})
	return nil
}

// set table view
func (t *Trv) setTableViewer() {
	t.TableViewer = tview.NewList()
	t.TableViewer.SetTitle("Result(Ctrl+r)")
	t.TableViewer.SetBorder(true)
}

// set search box
func (t *Trv) setSearcher() {
	t.Searcher = tview.NewInputField()
	t.Searcher.SetTitle("serach(Ctrl+s)")
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

// set new source form
func (t *Trv) setForm() {
	var source Source
	t.Form = tview.NewForm().
		AddInputField("Owner(required)", "", 20, nil, func(text string) {
			source.Owner = text
		}).
		AddInputField("Repo(required)", "", 20, nil, func(text string) {
			source.Repo = text
		}).
		AddInputField("Path(required)", "", 20, nil, func(text string) {
			source.Path = text
		}).
		AddPasswordField("Token(required)", "", 50, '*', func(text string) {
			source.Token = text
		}).
		AddCheckbox("IsEnterprise", false, func(checked bool) {
			source.IsEnterprise = checked
		}).
		AddInputField("BaseURL", "", 50, nil, func(text string) {
			source.BaseURL = text
		}).
		AddInputField("UploadURL", "", 50, nil, func(text string) {
			source.UploadURL = text
		}).
		AddButton("Save", func() {
			t.Config.addSource(source)
			t.setSource()
			t.addDropdownOption()
			t.Pages.HidePage("modal")
		}).
		AddButton("Quit", func() {
			t.Pages.HidePage("modal")
		})
	t.Form.SetBorder(true).SetTitle("add data source")
}

// set add source modal
func (t *Trv) setModal() {
	t.Modal = tview.NewGrid().
		SetColumns(0, 4, 0).
		SetRows(0, 4, 0).
		AddItem(t.Form, 0, 0, 4, 4, 0, 0, true)
}

// set pages
func (t *Trv) setPages() {
	t.Pages = tview.NewPages().
		AddPage("background", t.Layout, true, true).
		AddPage("modal", t.Modal, true, true)
}

// Set the layout of the area displaying information about the column
func (t *Trv) setInfoLayout() {
	t.InfoLayout = tview.NewGrid()
	t.InfoLayout.SetTitle("details").SetBorder(true)
	t.InfoLayout.SetSize(5, 5, 0, 0).
		AddItem(t.InfoText, 0, 0, 2, 5, 0, 0, true).
		AddItem(t.InfoTable, 2, 0, 3, 5, 2, 5, true)
	t.InfoLayout.SetOffset(1, 1)
}

// set layout
func (t *Trv) setLayout() {
	t.Layout = tview.NewGrid()
	t.Layout.SetSize(10, 10, 0, 0).
		AddItem(t.SourceSelecter, 0, 0, 2, 3, 0, 0, true).
		AddItem(t.Searcher, 0, 3, 2, 7, 0, 0, true).
		AddItem(t.TableViewer, 2, 0, 8, 5, 0, 0, true).
		AddItem(t.InfoLayout, 2, 5, 8, 5, 0, 0, true)
}

// Filter and display data
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
								comment := tview.NewTableCell(b.Comment)
								comment.SetMaxWidth(10)
								comment.SetExpansion(10)
								t.InfoTable.SetCell(1, 2, comment)
							}
						}
					}
				})
			}
		}
	}
}

func (t Trv) Draw() {
	runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}
	if err := t.App.SetRoot(t.Pages, true).Run(); err != nil {
		panic(err)
	}
}
