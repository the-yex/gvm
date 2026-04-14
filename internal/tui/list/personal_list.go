package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/core"
	progress2 "github.com/the-yex/gvm/internal/tui/progress"
	"github.com/the-yex/gvm/internal/version"
	"strings"
	"sync"
	"time"
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
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8C8C8C")).
			Padding(0, 2, 1, 0)
	helpHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C7C7C")).
			Padding(0, 2)
	helpHintActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#56CCF2")).
				Padding(0, 2)
	progressInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5BC0DE")).
				Padding(0, 2)
)

type Model struct {
	progress      *progress2.Model
	program       *tea.Program
	list          list.Model
	keys          *keyMap
	index         int
	local         bool
	quitting      bool
	statusTracker *statusTracker
	footer        FooterInfo
	lastError     string
	retryAction   func() tea.Cmd
}

type FooterInfo struct {
	Mirror  string
	Timeout time.Duration
	Remote  bool
}

type uninstallResultMsg struct {
	item *version.Version
	err  error
}

func NewListProgram(items []list.Item, title title, footer FooterInfo) *tea.Program {
	keys := newKeyMap()
	isLocal := title == LOCAL
	if isLocal {
		keys.install.SetEnabled(false)
	}
	tracker := newStatusTracker()
	versionList := list.New(items, delegate{keys: keys, status: tracker}, 0, 0)
	versionList.Title = string(title)
	versionList.Styles.Title = titleStyle
	helpKeys := func() []key.Binding {
		return []key.Binding{
			keys.install, keys.uninstall, keys.use,
		}
	}
	versionList.AdditionalShortHelpKeys = helpKeys
	versionList.AdditionalFullHelpKeys = helpKeys
	model := &Model{
		list:          versionList,
		keys:          keys,
		local:         title == LOCAL,
		statusTracker: tracker,
		footer:        footer,
	}
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
	m.annotateStatus(item, "下载中...")
	m.progress = progress2.NewModel(m.program)
	defer func() { m.progress = nil }()
	err := core.MultiWriterInstall(item, m.progress.MultiWriter(nil), m.progress.SetSize)
	if err != nil {
		m.annotateStatus(item, "安装失败")
		msg := "安装失败: " + err.Error()
		m.setError(msg, func() tea.Cmd {
			go m.download(item)
			return nil
		})
		m.program.Send(m.list.NewStatusMessage(failMessageStyle(msg)))
		m.clearStatusAfter(item, 5*time.Second)
		return
	}
	item.DirName = fmt.Sprintf("go%s", item.String())
	m.annotateStatus(item, "安装完成")
	setCmd := m.list.SetItem(m.list.Index(), item)
	statusCmd := m.list.NewStatusMessage(successMessageStyle("success install " + item.String()))
	m.program.Send(tea.Batch(setCmd, statusCmd))
	m.clearStatusAfter(item, 2*time.Second)
	m.clearError()
}

func (m *Model) processUninstall(item *version.Version) {
	if item == nil {
		return
	}
	m.annotateStatus(item, "卸载中...")
	go func(it *version.Version) {
		err := core.UninstallVersion(it.LocalDir())
		m.program.Send(uninstallResultMsg{item: it, err: err})
	}(item)
}

func (m *Model) annotateStatus(item *version.Version, status string) {
	if m.statusTracker == nil || item == nil {
		return
	}
	if status == "" {
		m.statusTracker.Clear(item.String())
		return
	}
	m.statusTracker.Set(item.String(), status)
}

func (m *Model) setError(msg string, action func() tea.Cmd) {
	if msg == "" {
		m.clearError()
		return
	}
	m.lastError = msg
	m.retryAction = action
}

func (m *Model) clearError() {
	m.lastError = ""
	m.retryAction = nil
}

type statusTracker struct {
	mu       sync.RWMutex
	statuses map[string]string
}

func newStatusTracker() *statusTracker {
	return &statusTracker{
		statuses: make(map[string]string),
	}
}

func (s *statusTracker) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statuses[key]
}

func (s *statusTracker) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[key] = value
}

func (s *statusTracker) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.statuses, key)
}

func (m *Model) renderFooter() string {
	mirror := m.footer.Mirror
	if mirror == "" {
		mirror = viper.GetString(consts.CONFIG_MIRROR)
	}
	if mirror == "" {
		mirror = consts.DEFAULT_MIRROR
	}
	t := m.footer.Timeout
	if t <= 0 {
		t = 5 * time.Second
	}
	mode := "本地模式"
	if m.footer.Remote {
		mode = "远程模式"
	}
	return footerStyle.Render(fmt.Sprintf(" %s · 镜像：%s · 超时：%s", mode, mirror, durationLabel(t)))
}

func (m *Model) renderLastError() string {
	if m.lastError == "" {
		return ""
	}
	return failMessageStyle("⚠ " + m.lastError + " (按 r 重试)")
}

func (m *Model) renderContextHelp() string {
	item, _ := m.list.SelectedItem().(*version.Version)
	if item == nil {
		return ""
	}
	type hint struct {
		key       string
		label     string
		available bool
	}
	hints := []hint{
		{key: "i", label: "安装", available: item.LocalDir() == "" && m.progress == nil},
		{key: "x", label: "卸载", available: item.LocalDir() != "" && !item.CurrentUsed},
		{key: "u", label: "使用", available: item.DirName != "" && !item.CurrentUsed},
		{key: "r", label: "重试", available: m.retryAction != nil},
	}
	parts := make([]string, 0, len(hints))
	for _, h := range hints {
		style := helpHintStyle
		if h.available {
			style = helpHintActiveStyle
		}
		parts = append(parts, style.Render(fmt.Sprintf("[%s] %s", strings.ToUpper(h.key), h.label)))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderProgressSummary() string {
	if m.progress == nil {
		return ""
	}
	ratio := m.progress.Ratio()
	eta := formatDurationDesc(m.progress.Remaining())
	written := formatBytes(m.progress.WrittenBytes())
	total := formatBytes(m.progress.TotalBytes())
	return progressInfoStyle.Render(fmt.Sprintf("进度 %s · %s/%s · ETA %s",
		formatPercent(ratio), written, total, eta))
}

func formatPercent(ratio float64) string {
	return fmt.Sprintf("%.0f%%", ratio*100)
}

func formatSpeedDesc(speed float64) string {
	if speed >= 1024*1024 {
		return fmt.Sprintf("%.2f MB/s", speed/1024/1024)
	}
	return fmt.Sprintf("%.2f KB/s", speed/1024)
}

func formatDurationDesc(d time.Duration) string {
	if d < 0 {
		return "00:00"
	}
	min := int(d.Minutes())
	sec := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func formatBytes(bytes int64) string {
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
	}
	return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
}

func durationLabel(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}
	return fmt.Sprintf("%.0fs", d.Round(time.Second).Seconds())
}

func (m *Model) clearStatusAfter(item *version.Version, d time.Duration) {
	if d <= 0 || item == nil {
		return
	}
	go func() {
		time.Sleep(d)
		m.annotateStatus(item, "")
	}()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var item *version.Version

	if i, ok := m.list.SelectedItem().(*version.Version); ok {
		item = i
	}
	switch msg := msg.(type) {
	case uninstallResultMsg:
		item := msg.item
		if item == nil {
			return m, nil
		}
		if msg.err != nil {
			msgStr := "卸载失败: " + msg.err.Error()
			m.annotateStatus(item, "卸载失败")
			m.clearStatusAfter(item, 5*time.Second)
			m.setError(msgStr, func() tea.Cmd {
				m.processUninstall(item)
				return nil
			})
			return m, tea.Batch(m.list.NewStatusMessage(failMessageStyle(msgStr)))
		}
		m.annotateStatus(item, "卸载完成")
		m.clearStatusAfter(item, 2*time.Second)
		m.clearError()
		if m.local {
			m.list.RemoveItem(m.list.Index())
		} else {
			item.DirName = ""
			item.Installed = false
			m.list.SetItem(m.list.Index(), item)
		}
		return m, tea.Batch(m.list.NewStatusMessage(successMessageStyle("success uninstall " + item.String())))
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
			m.processUninstall(item)
			return m, nil
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
			m.clearError()
			return m, tea.Batch(setCmd, statusCmd)
		case key.Matches(msg, m.keys.retry):
			if m.retryAction == nil {
				return m, nil
			}
			return m, m.retryAction()
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
		if summary := m.renderProgressSummary(); summary != "" {
			s += "\n" + summary
		}
	}
	if help := m.renderContextHelp(); help != "" {
		s += "\n" + help
	}
	if errLine := m.renderLastError(); errLine != "" {
		s += "\n" + errLine
	}
	footer := m.renderFooter()
	if footer != "" {
		s += "\n" + footer
	}
	return s
}
