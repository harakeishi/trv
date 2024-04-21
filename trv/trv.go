package trv

import (
	"fmt"
	"log"
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
	TableViewer    *tview.Table
	Searcher       *tview.InputField
	Pages          *tview.Pages
	InfoLayout     *tview.Grid
	InfoText       *tview.TextView
	App            *tview.Application
	Layout         *tview.Grid
	Modal          tview.Primitive
	ErrorModal     tview.Primitive
	Form           *tview.Form
	ErrorWindow    *tview.Modal
}

type Info struct {
	Table  Table
	Column Column
}

/*
Prepare to start trv.
Specifically, loading the configuration, loading the data, and preparing the TUI.
*/
func (t *Trv) Init() error {
	runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}
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
	t.setInfoLayout()
	t.setLayout()
	t.setForm()
	t.setModal()
	t.setErrorModal()
	t.setPages()

	t.Pages.HidePage("modal")
	t.Pages.HidePage("error")

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
			if err != nil {
				t.ErrorWindow.SetText(fmt.Sprintf("-ERROR-\n%s", err))
				t.ErrorWindow.Box.SetTitle("ERROR")
				t.Pages.ShowPage("error")
				return
			}
			t.Tables = t.DB.Tables
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
		t.Tables = t.DB.Tables
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
	t.TableViewer = tview.NewTable().SetBorders(false)
	t.TableViewer.SetTitle("Result(Ctrl+r)")
	t.TableViewer.SetBorder(true)
	t.TableViewer.SetSelectable(true, false)
	t.TableViewer.SetSelectedFunc(func(row int, column int) {
		cell := t.TableViewer.GetCell(row, column)
		info := cell.GetReference().(Info)
		t.InfoText.SetText(fmt.Sprintf("table name: %s\ndetails: %s\n\ncolumn: %s\ntype: %s\ncomment: %s", info.Table.Name, info.Table.Description, info.Column.Name, info.Column.Type, info.Column.Comment))
	})
}

// set search box
func (t *Trv) setSearcher() {
	t.Searcher = tview.NewInputField()
	t.Searcher.SetTitle("search(Ctrl+s)")
	t.Searcher.SetLabel("search:")
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
			t.Config.saveConfig()
			t.setSource()
			if err := t.addDropdownOption(); err != nil {
				log.Printf("addDropdownOption fail:%s", err)
			}
			t.Pages.HidePage("modal")
		}).
		AddButton("Quit", func() {
			t.Pages.HidePage("modal")
		})
	t.Form.SetBorder(true).SetTitle("add data source")
}

// set source modal
func (t *Trv) setModal() {
	t.Modal = tview.NewGrid().
		SetColumns(0, 4, 0).
		SetRows(0, 4, 0).
		AddItem(t.Form, 0, 0, 4, 4, 0, 0, true)
}

// set error modal
func (t *Trv) setErrorModal() {
	t.ErrorWindow = tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				t.App.Stop()
			}
		})

	t.ErrorModal = tview.NewGrid().
		SetColumns(0, 4, 0).
		SetRows(0, 4, 0).
		AddItem(t.ErrorWindow, 0, 0, 4, 4, 0, 0, true)
}

// set pages
func (t *Trv) setPages() {
	t.Pages = tview.NewPages().
		AddPage("background", t.Layout, true, true).
		AddPage("modal", t.Modal, true, true).
		AddPage("error", t.ErrorModal, true, true)
}

// Set the layout of the area displaying information about the column
func (t *Trv) setInfoLayout() {
	t.InfoLayout = tview.NewGrid()
	t.InfoLayout.SetTitle("details").SetBorder(true)
	t.InfoLayout.SetSize(5, 5, 0, 0).
		AddItem(t.InfoText, 0, 0, 5, 5, 0, 0, true)
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
	row := 0
	t.TableViewer.Clear()
	for _, r := range t.Tables {
		for i, c := range r.Columns {
			if strings.Contains(strings.ToLower(r.getFullName(i)), strings.ToLower(target)) || target == "" {
				t.TableViewer.SetCell(row, 0, tview.NewTableCell(r.getFullName(i)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft))
				t.TableViewer.SetCell(row, 1, tview.NewTableCell(c.Comment).
					SetTextColor(tcell.ColorBeige).
					SetAlign(tview.AlignLeft))
				cell := t.TableViewer.GetCell(row, 0)
				cell.SetReference(Info{Table: r, Column: c})
				row++
			}
		}
	}
	if t.TableViewer.GetRowCount() == 0 {
		t.TableViewer.SetCell(row, 0, tview.NewTableCell("No result").
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))
	}
}

// drawing
func (t Trv) Draw() {
	if err := t.App.SetRoot(t.Pages, true).Run(); err != nil {
		panic(err)
	}
}

// set Configuration form
func (t Trv) CreateConfig() {
	runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}
	t.App = tview.NewApplication()
	var source Source
	Form := tview.NewForm().
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
			t.Config.saveConfig()
			t.setSource()
			t.App.Stop()
		}).
		AddButton("Quit", func() {
			t.App.Stop()
		})
	Form.SetBorder(true).SetTitle("create config")
	if err := t.App.SetRoot(Form, true).Run(); err != nil {
		panic(err)
	}
}
