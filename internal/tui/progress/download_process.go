package progress

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
	"time"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/26 下午2:28
* @Package:
 */

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding    = 1
	maxWidth   = 100
	speedQueue = 5
)

type progressMsg struct {
	ratio      float64
	speed      float64
	remain     time.Duration
	written    int64
	totalBytes int64
}
type Model struct {
	program    *tea.Program
	progress   progress.Model
	cancel     bool
	speed      float64
	remain     time.Duration
	written    int64
	totalBytes int64
	writer     *ProgressWriter
}

func NewModel(program *tea.Program) *Model {
	m := &Model{
		progress: progress.New(progress.WithDefaultGradient()),
		program:  program,
	}
	if m.program == nil {
		m.program = tea.NewProgram(m, tea.WithAltScreen())
	}
	return m
}

func (m *Model) Start() {
	m.program.Run()
}

func (m *Model) Quit() {
	m.program.Quit()
}
func (m *Model) SetSize(size int64) {
	m.writer.total = size
}
func (m *Model) Size() int64 {
	return m.written
}
func (m *Model) MultiWriter(fw io.Writer) io.Writer {
	w := &ProgressWriter{
		start:        time.Now(),
		speedHistory: make([]float64, 0, speedQueue),
		onProgress:   func(msg tea.Msg) { m.program.Send(msg) },
	}
	m.writer = w

	if fw == nil {
		return w
	}
	return io.MultiWriter(fw, w)
}
func (m *Model) IsCancel() bool {
	return m.cancel
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.cancel = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case progressMsg:
		m.speed = msg.speed
		m.remain = msg.remain
		m.written = msg.written
		m.totalBytes = msg.totalBytes
		return m, m.progress.SetPercent(msg.ratio)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *Model) View() string {
	pad := strings.Repeat(" ", padding)
	sizeInfo := fmt.Sprintf("%s/%s", formatSize(m.written), formatSize(m.totalBytes))
	return "\n" + pad + m.progress.View() + "\n\n" +
		pad + helpStyle(fmt.Sprintf("press q to  cancel | speed：%s | ETA：%s | %s",
		formatSpeed(m.speed), formatETA(m.remain), sizeInfo))
}

func (m *Model) CurrentSpeed() float64 {
	return m.speed
}

func (m *Model) Remaining() time.Duration {
	return m.remain
}

func (m *Model) WrittenBytes() int64 {
	return m.written
}

func (m *Model) TotalBytes() int64 {
	return m.totalBytes
}

func (m *Model) Ratio() float64 {
	if m.totalBytes == 0 {
		return 0
	}
	return float64(m.written) / float64(m.totalBytes)
}

func formatSpeed(speed float64) string {
	if speed >= 1024*1024 {
		return fmt.Sprintf("%.2f MB/s", speed/1024/1024)
	}
	return fmt.Sprintf("%.2f KB/s", speed/1024)
}

func formatETA(d time.Duration) string {
	min := int(d.Minutes())
	sec := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

// 格式化文件大小
func formatSize(bytes int64) string {
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
	}
	return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
}

type ProgressWriter struct {
	total        int64
	written      int64
	start        time.Time
	speedHistory []float64
	onProgress   func(msg tea.Msg)
}

func (pw *ProgressWriter) SetSize(size int64) {
	pw.total = size
}
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.written += int64(n)

	now := time.Now()
	elapsed := now.Sub(pw.start).Seconds()
	if elapsed <= 0 {
		elapsed = 0.001
	}

	instSpeed := float64(n) / elapsed
	pw.speedHistory = append(pw.speedHistory, instSpeed)
	if len(pw.speedHistory) > speedQueue {
		pw.speedHistory = pw.speedHistory[1:]
	}

	var sum float64
	for _, s := range pw.speedHistory {
		sum += s
	}
	avgSpeed := sum / float64(len(pw.speedHistory))
	remain := time.Duration(float64(pw.total-pw.written)/avgSpeed) * time.Second

	if pw.onProgress != nil && pw.total > 0 {
		pw.onProgress(progressMsg{
			ratio:      float64(pw.written) / float64(pw.total),
			speed:      avgSpeed,
			remain:     remain,
			written:    pw.written,
			totalBytes: pw.total,
		})
	}
	pw.start = now
	return n, nil
}
