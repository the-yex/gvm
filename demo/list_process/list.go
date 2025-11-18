package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	progress2 "github.com/the-yex/gvm/internal/tui/progress"
	"github.com/the-yex/gvm/internal/utils"
	"io"
	"os"
	"strings"
	"time"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/26 上午11:12
* @Package:
 */

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/24 上午10:29
* @Package:
 */
const SimpleListHeight = 14
const (
	padding  = 2
	maxWidth = 80
)

type tickMsg time.Time

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().Underline(true).PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	titleStyle        = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#25A065")).
				Padding(0, 0)
)

type SimpleListItem string

func (i SimpleListItem) FilterValue() string { return "" }

type SimpleListItemDelegate struct{}

func (d SimpleListItemDelegate) Height() int                             { return 1 }
func (d SimpleListItemDelegate) Spacing() int                            { return 0 }
func (d SimpleListItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SimpleListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SimpleListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("go%s", i)

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
	progress *progress2.Model
	program  *tea.Program
	index    int
	quitting bool
}

func NewListProgram(list list.Model) *tea.Program {
	list.Styles.Title = titleStyle
	list.Styles.HelpStyle = helpStyle
	list.Styles.PaginationStyle = paginationStyle
	model := &SimpleListModel{list: list}
	p := tea.NewProgram(model, tea.WithAltScreen())
	model.program = p
	return p
}

func (m *SimpleListModel) Index() int {
	return m.index
}
func (m *SimpleListModel) Init() tea.Cmd {
	return nil
}
func (m *SimpleListModel) download(pm *progress2.Model) {
	sourceUrl := "https://golang.google.cn/dl/go1.24.6.darwin-arm64.tar.gz"
	f, err := os.OpenFile("./go1.24.6.darwin-arm64.tar.gz", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	utils.Download(sourceUrl, pm.MultiWriter(f), pm.SetSize)
	m.progress = nil
}
func (m *SimpleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			m.index = -1
			return m, tea.Quit

		case "enter":
			if m.progress != nil {
				return m, nil
			}
			m.progress = progress2.NewModel(m.program)
			go m.download(m.progress)
			return m, tickCmd()
		}
	}
	var cmds []tea.Cmd
	if m.progress != nil {
		pm, pCmd := m.progress.Update(msg)
		cmds = append(cmds, pCmd)
		m.progress = pm.(*progress2.Model)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *SimpleListModel) View() string {
	if m.quitting {
		return "\n"
	}
	s := "\n" + m.list.View()
	if m.progress != nil {
		pad := strings.Repeat(" ", padding)
		s += "\n" +
			pad + m.progress.View()
	}
	return s
}
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
