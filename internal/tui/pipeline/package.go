package pipeline

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/10/9 上午10:44
* @Package:
 */

var orderedStages = []string{check, download, unzip, installing}

type model struct {
	program         *tea.Program
	pipeline        upgradePipeline
	width           int
	height          int
	spinner         spinner.Model
	progress        progress.Model
	done            bool
	success         bool
	message         string
	currentVersion  string
	targetVersion   string
	platform        string
	assetName       string
	completed       []string
	failedStage     string
	downloadWritten int64
	downloadTotal   int64
}

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	activeStageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	doneStageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	failedStageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	pendingStageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	metaLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	dimStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	boxStyle          = lipgloss.NewStyle().Padding(1, 2)
)

func NewtProgram() *tea.Program {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	m := &model{
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
	}
	p := tea.NewProgram(m)
	m.program = p
	m.pipeline = searchPipeline{send: p.Send}
	return p
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.pipeline.Do(), m.spinner.Tick)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.progress.Width = max(20, msg.Width-16)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case downloadProgressMsg:
		m.downloadWritten = msg.written
		m.downloadTotal = msg.total
		if msg.total > 0 {
			cmd := m.progress.SetPercent(float64(msg.written) / float64(msg.total))
			return m, cmd
		}
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if pm, ok := newModel.(progress.Model); ok {
			m.progress = pm
		}
		return m, cmd
	case pipelineMsg:
		if msg.currentVersion != "" {
			m.currentVersion = msg.currentVersion
		}
		if msg.targetVersion != "" {
			m.targetVersion = msg.targetVersion
		}
		if msg.platform != "" {
			m.platform = msg.platform
		}
		if msg.assetName != "" {
			m.assetName = msg.assetName
		}
		if msg.info != "" {
			m.message = msg.info
		}
		if msg.stage == nil {
			m.done = true
			m.success = msg.success
			if m.success {
				if m.pipeline != nil && !containsStage(m.completed, m.pipeline.String()) {
					m.completed = append(m.completed, m.pipeline.String())
				}
			} else if m.pipeline != nil {
				m.failedStage = m.pipeline.String()
			}
			return m, tea.Quit
		}
		if !containsStage(m.completed, m.pipeline.String()) {
			m.completed = append(m.completed, m.pipeline.String())
		}
		if msg.stage.String() != download {
			m.downloadWritten = 0
			m.downloadTotal = 0
		}
		m.pipeline = msg.stage
		return m, m.pipeline.Do()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() string {
	lines := []string{
		titleStyle.Render("GVM Upgrade"),
		m.renderVersionLine(),
		m.renderMetaLine(),
		"",
		m.renderStages(),
	}
	if downloadInfo := m.renderDownloadInfo(); downloadInfo != "" {
		lines = append(lines, "", downloadInfo)
	}
	if m.message != "" {
		lines = append(lines, "", infoStyle.Render(m.message))
	}
	if !m.done {
		lines = append(lines, "", dimStyle.Render("按 q 退出"))
	} else if m.success {
		lines = append(lines, "", successStyle.Render("升级流程完成"))
		lines = append(lines, dimStyle.Render("可执行 gvm --version 验证新版本"))
	} else {
		lines = append(lines, "", errorStyle.Render("升级流程未完成"))
		lines = append(lines, dimStyle.Render("可重新执行 gvm upgrade 重试"))
	}
	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m *model) renderVersionLine() string {
	current := m.currentVersion
	if current == "" {
		current = "unknown"
	}
	target := m.targetVersion
	if target == "" {
		target = "-"
	}
	return fmt.Sprintf("当前版本: %s   目标版本: %s", current, target)
}

func (m *model) renderMetaLine() string {
	parts := make([]string, 0, 2)
	if m.platform != "" {
		parts = append(parts, fmt.Sprintf("%s %s", metaLabelStyle.Render("平台:"), m.platform))
	}
	if m.assetName != "" {
		parts = append(parts, fmt.Sprintf("%s %s", metaLabelStyle.Render("升级包:"), m.assetName))
	}
	return strings.Join(parts, "   ")
}

func (m *model) renderStages() string {
	parts := make([]string, 0, len(orderedStages))
	for _, stage := range orderedStages {
		switch {
		case containsStage(m.completed, stage):
			parts = append(parts, doneStageStyle.Render("✓ "+stage))
		case m.done && m.failedStage == stage:
			parts = append(parts, failedStageStyle.Render("x "+stage))
		case !m.done && m.pipeline != nil && m.pipeline.String() == stage:
			parts = append(parts, activeStageStyle.Render(m.spinner.View()+" "+stage))
		default:
			parts = append(parts, pendingStageStyle.Render("· "+stage))
		}
	}
	return strings.Join(parts, "\n")
}

func (m *model) renderDownloadInfo() string {
	if m.pipeline == nil || m.pipeline.String() != download || m.downloadTotal <= 0 {
		return ""
	}
	percent := float64(m.downloadWritten) / float64(m.downloadTotal)
	summary := fmt.Sprintf("%s   %s/%s",
		formatPercent(percent),
		formatBytes(m.downloadWritten),
		formatBytes(m.downloadTotal),
	)
	return m.progress.View() + "\n" + infoStyle.Render(summary)
}

func containsStage(stages []string, target string) bool {
	for _, stage := range stages {
		if stage == target {
			return true
		}
	}
	return false
}

func formatPercent(ratio float64) string {
	return fmt.Sprintf("%.0f%%", ratio*100)
}

func formatBytes(bytes int64) string {
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
	}
	return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
