package spinner

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/23 下午3:37
* @Package:
 */

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/the-yex/gvm/internal/core"
)

type errMsg error

type spinnerModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func newSpinnerProgram(options ...tea.ProgramOption) *tea.Program {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return tea.NewProgram(spinnerModel{spinner: s}, options...)
}
func init() {
	core.NewSpinnerProgram = newSpinnerProgram
}
func NewSpinner() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spinnerModel{spinner: s}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading ...\n\n", m.spinner.View())
	if m.quitting {
		return "\n"
	}
	return str
}
