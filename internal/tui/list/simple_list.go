package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-yex/gvm/internal/core"
	"github.com/the-yex/gvm/internal/version"
	"io"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/24 上午10:29
* @Package:
 */
const SimpleListHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().Underline(true).PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type SimpleListItemDelegate struct{}

func (d SimpleListItemDelegate) Height() int                             { return 1 }
func (d SimpleListItemDelegate) Spacing() int                            { return 0 }
func (d SimpleListItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SimpleListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*version.Version)
	if !ok {
		return
	}

	str := fmt.Sprintf("go%s", i.String())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

type SimpleListModel struct {
	list     list.Model
	index    int
	quitting bool
}
type ChoiceHandler func(choice string)

func NewSimpleListModel(list list.Model) SimpleListModel {
	list.Styles.Title = titleStyle
	list.Styles.HelpStyle = helpStyle
	list.Styles.PaginationStyle = paginationStyle
	return SimpleListModel{list: list}
}
func newSimpleListProgram(items []list.Item, title string, options ...tea.ProgramOption) *tea.Program {
	//listItems := make([]list.Item, 0, len(items))
	//for _, s := range items {
	//	listItems = append(listItems, SimpleListItem(s))
	//}
	l := list.New(items, SimpleListItemDelegate{}, 20, SimpleListHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.PaginationStyle = paginationStyle
	return tea.NewProgram(SimpleListModel{list: l}, options...)
}
func init() {
	core.NewSimpleListProgram = newSimpleListProgram
}

func (m SimpleListModel) Index() int {
	return m.index
}
func (m SimpleListModel) Init() tea.Cmd {
	return nil
}

func (m SimpleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			m.index = -1
			return m, tea.Quit

		case "enter":
			m.index = m.list.Index()
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SimpleListModel) View() string {
	if m.quitting {
		return "\n"
	}

	return "\n" + m.list.View()
}
