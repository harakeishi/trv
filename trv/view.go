package trv

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	defaultStyle = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("62"))
)

func Draw2() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type model struct {
	config Config
	dbs    list.Model
	tables list.Model
	// 現在のビューを追跡（"list" または "table"）
	currentView   string
	currentDb     DB
	width         int
	height        int
	currentSource Source
	currentColmun string
}

func (m *model) setConfig() error {
	if err := m.config.loadConfig(); err != nil {
		return fmt.Errorf("setConfig fail: %w", err)
	}
	return nil
}
func (m *model) setDbList() error {
	items := make([]list.Item, 0)
	for _, v := range m.config.Source {
		item := listItem{source: v}
		item.setTitle()
		items = append(items, item)
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Data sauce"
	m.dbs = l
	return nil
}

func (m *model) setTableList() error {
	err := m.currentDb.loadData(m.currentSource.Repo, m.currentSource.Path)
	fmt.Println(err)
	items := make([]list.Item, 0)
	for _, table := range m.currentDb.Tables {
		for _, column := range table.Columns {
			item := listItem{title: fmt.Sprintf("%s.%s", table.Name, column.Name), column: column}
			items = append(items, item)
		}
	}
	items = append(items, listItem{title: fmt.Sprint(err)})
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "show DB Info"
	l.SetShowHelp(false)
	m.tables = l
	return nil
}
func newModel() model {

	items2 := []list.Item{
		listItem{title: "項目 1"},
		listItem{title: "項目 2"},
		listItem{title: "項目 3"},
	}

	t := list.New(items2, list.NewDefaultDelegate(), 0, 0)
	t.Title = "tables"
	model := model{
		tables:      t,
		currentView: "list",
	}

	if err := model.setConfig(); err != nil {
		fmt.Println(err)
	}
	model.setDbList()
	return model
}

type listItem struct {
	title  string
	source Source
	column Column
}

func (l *listItem) setTitle() {
	l.title = fmt.Sprintf("%s/%s", l.source.Owner, l.source.Repo)
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlD:
			if m.currentView == "table" {
				m.currentView = "list"
			}
		case tea.KeyEnter:
			if m.currentView == "list" {
				selectedIndex := m.dbs.Index()
				item := m.dbs.Items()[selectedIndex].(listItem)
				m.currentSource = item.source
				m.setTableList()
				m.tables.SetSize((m.width), m.height/3*2)
				m.currentView = "table"
			} else if m.currentView == "table" {
				selectedIndex := m.tables.Index()
				item := m.tables.Items()[selectedIndex].(listItem)
				m.currentColmun = fmt.Sprintf("%+v", item.column)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.dbs.SetSize(msg.Width-h, msg.Height-v)
		m.tables.SetSize((msg.Width - h), (msg.Height - v))
	}

	if m.currentView == "list" {
		var cmd tea.Cmd
		m.dbs, cmd = m.dbs.Update(msg)
		return m, cmd
	} else if m.currentView == "table" {
		var cmd tea.Cmd
		m.tables, cmd = m.tables.Update(msg)
		return m, cmd
	}

	return m, nil
}

var tabls, column lipgloss.Style

func (m model) View() string {
	if m.currentView == "list" {
		return docStyle.Render(m.dbs.View())
	} else if m.currentView == "table" {
		tabls = defaultStyle.Width(m.width / 3 * 2).Height(m.height - 1)
		column = defaultStyle.Width(m.width / 3).Height(m.height - 1)
		return lipgloss.JoinHorizontal(lipgloss.Top, tabls.Render(fmt.Sprintf("%4s", m.tables.View())), column.Render(fmt.Sprintf("%4s", m.currentColmun))) + helpStyle.Render("\ntab: focus next • n: new  • q: exit\n")
	}
	return ""
}
