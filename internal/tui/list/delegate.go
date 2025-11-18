package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-yex/gvm/internal/version"
	"io"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/25 上午10:28
* @Package:
 */
type delegate struct {
	keys *keyMap
}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*version.Version)
	if !ok {
		return
	}

	data := fmt.Sprintf(" %s", i.String())
	if i.LocalDir() != "" {
		data = fmt.Sprintf("%s  %s", data, i.LocalDir())
	}

	normal := lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("#888888")) // 未安装灰色
	local := lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("#00FF00"))  // 已安装绿色
	selected := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)
	current := lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("#FFD700")).Bold(true) // 当前使用金色
	fn := normal.Render

	switch {
	case i.CurrentUsed:
		fn = current.Render
		if index == m.Index() {
			fn = func(s ...string) string {
				return current.Underline(true).Render("> " + strings.Join(s, " "))
			}
		}
	case index == m.Index():
		fn = func(s ...string) string {
			return selected.Underline(true).Render("> " + strings.Join(s, " "))
		}
	case i.LocalDir() != "":
		fn = local.Render

	}

	fmt.Fprint(w, fn(data))
}

type keyMap struct {
	install   key.Binding
	uninstall key.Binding
	use       key.Binding
}

func (d keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.install,
		d.uninstall,
		d.use,
	}
}

func (d keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.install,
			d.uninstall,
			d.use,
		},
	}
}

func newKeyMap() *keyMap {
	return &keyMap{
		install: key.NewBinding(
			key.WithKeys("enter", "i"),
			key.WithHelp("i", "install"),
		),
		uninstall: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "uninstall"),
		),
		use: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "use")),
	}
}
