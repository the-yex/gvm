package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-yex/gvm/internal/core"
	progress2 "github.com/the-yex/gvm/internal/tui/progress"
	"github.com/the-yex/gvm/internal/version"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/25 上午10:46
* @Package:
 */
type title string

const (
	Remote title = "remote versions"
	LOCAL  title = "local versions"
)

var (
	appStyle = lipgloss.NewStyle().Padding(0, 0).Margin(1, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 0)

	successMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	warnMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#FFD700"}).
				Render

	failMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FF4B4B", Dark: "#FF4B4B"}).
				Render
)

type Model struct {
	progress *progress2.Model
	program  *tea.Program
	list     list.Model
	keys     *keyMap
	index    int
	local    bool
	quitting bool
}

func NewListProgram(items []list.Item, title title) *tea.Program {
	keys := newKeyMap()
	isLocal := title == LOCAL
	if isLocal {
		keys.install.SetEnabled(false)
	}
	versionList := list.New(items, delegate{keys: keys}, 0, 0)
	versionList.Title = string(title)
	versionList.Styles.Title = titleStyle
	helpKeys := func() []key.Binding {
		return []key.Binding{
			keys.install, keys.uninstall, keys.use,
		}
	}
	versionList.AdditionalShortHelpKeys = helpKeys
	versionList.AdditionalFullHelpKeys = helpKeys
	model := &Model{list: versionList, keys: keys, local: title == LOCAL}
	program := tea.NewProgram(model, tea.WithAltScreen())
	model.program = program
	return program
}

func (m *Model) Index() int {
	return m.index
}
func (m *Model) Init() tea.Cmd {
	return nil
}
func (m *Model) download(item *version.Version) {
	m.progress = progress2.NewModel(m.program)
	defer func() { m.progress = nil }()
	err := core.MultiWriterInstall(item, m.progress.MultiWriter(nil), m.progress.SetSize)
	if err != nil {
		return
	}
	item.DirName = fmt.Sprintf("go%s", item.String())
	setCmd := m.list.SetItem(m.list.Index(), item)
	statusCmd := m.list.NewStatusMessage(successMessageStyle("success install " + item.String()))
	m.program.Send(tea.Batch(setCmd, statusCmd))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var item *version.Version

	if i, ok := m.list.SelectedItem().(*version.Version); ok {
		item = i
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.install):
			if item.LocalDir() != "" {
				return m, nil
			}
			if m.progress != nil {
				return m, nil
			}
			go m.download(item)
			return m, nil
		case key.Matches(msg, m.keys.uninstall):
			if item.CurrentUsed {
				return m, m.list.NewStatusMessage(warnMessageStyle("can not uninstall current used version " + item.String()))
			}
			if item.LocalDir() == "" {
				return m, nil
			}
			err := core.UninstallVersion(item.LocalDir())
			if err != nil {
				return m, m.list.NewStatusMessage(failMessageStyle(err.Error()))
			}
			if m.local {
				m.list.RemoveItem(m.list.Index())
			} else {
				item.DirName = ""
				m.list.SetItem(m.list.Index(), item)
			}
			return m, m.list.NewStatusMessage(successMessageStyle("success uninstall " + item.String()))
		case key.Matches(msg, m.keys.use):
			if item.CurrentUsed || item.DirName == "" {
				return m, nil
			}
			err := core.SwitchVersion(item.LocalDir())
			if err != nil {
				return m, m.list.NewStatusMessage(failMessageStyle(err.Error()))
			}
			for i, v := range m.list.Items() {
				if vi := v.(*version.Version); vi.CurrentUsed {
					vi.CurrentUsed = false
					m.list.SetItem(i, vi)
				}
			}
			item.CurrentUsed = true
			setCmd := m.list.SetItem(m.list.Index(), item)
			statusCmd := m.list.NewStatusMessage(successMessageStyle("current use " + item.String()))
			return m, tea.Batch(setCmd, statusCmd)
		}
	}
	if m.progress != nil {
		pm, pCmd := m.progress.Update(msg)
		cmds = append(cmds, pCmd)
		m.progress = pm.(*progress2.Model)
		return m, tea.Batch(cmds...)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	s := appStyle.Render(m.list.View())
	if m.progress != nil {
		pad := strings.Repeat(" ", 2)
		s += "\n" +
			pad + m.progress.View()
	}
	return s
}
