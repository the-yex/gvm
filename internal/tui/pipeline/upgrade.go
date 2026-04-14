package pipeline

import (
	"fmt"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/github"
)

const (
	check      = "检查线上版本"
	download   = "下载升级包"
	unzip      = "解压升级包"
	installing = "安装新版本"
)

type upgradePipeline interface {
	String() string
	Do() tea.Cmd
	Next() upgradePipeline
	Version() string
}

type downloadProgressMsg struct {
	written int64
	total   int64
}

type pipelineMsg struct {
	stage          upgradePipeline
	info           string
	currentVersion string
	targetVersion  string
	platform       string
	assetName      string
	success        bool
}

type searchPipeline struct {
	asset *github.Asset
	send  func(tea.Msg)
}

func (s searchPipeline) Version() string {
	if s.asset == nil {
		return ""
	}
	return s.asset.Version()
}

func (s searchPipeline) String() string { return check }

func (s searchPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{
			currentVersion: displayVersion(consts.Version),
			platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		}
		release, hasUpdate, err := github.NewReleaseUpdater().CheckForUpdates()
		if err != nil {
			msg.info = "检查更新失败: " + err.Error()
			return msg
		}
		msg.targetVersion = displayVersion(release.TagName)
		if !hasUpdate {
			msg.success = true
			msg.info = fmt.Sprintf("当前已经是最新版本 %s", msg.targetVersion)
			return msg
		}

		asset, err := release.FindAsset()
		if err != nil {
			msg.info = "未找到当前平台可用的升级包: " + err.Error()
			return msg
		}
		s.asset = asset
		msg.stage = s.Next()
		msg.assetName = asset.Name
		msg.info = fmt.Sprintf("发现新版本 %s，开始升级", msg.targetVersion)
		return msg
	}
}

func (s searchPipeline) Next() upgradePipeline {
	return downloadPipeline{asset: s.asset, send: s.send}
}

type downloadPipeline struct {
	asset *github.Asset
	send  func(tea.Msg)
}

func (d downloadPipeline) Version() string {
	if d.asset == nil {
		return ""
	}
	return d.asset.Version()
}

func (d downloadPipeline) String() string { return download }

func (d downloadPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{
			stage:         d.Next(),
			targetVersion: d.Version(),
			assetName:     d.asset.Name,
		}
		if _, err := d.asset.DownloadWithProgress(func(written, total int64) {
			if d.send != nil {
				d.send(downloadProgressMsg{written: written, total: total})
			}
		}); err != nil {
			d.asset.Clean()
			msg.info = fmt.Sprintf("下载失败: %s", err.Error())
			msg.stage = nil
		}
		return msg
	}
}

func (d downloadPipeline) Next() upgradePipeline {
	return unzipPipeline{asset: d.asset}
}

type unzipPipeline struct {
	asset *github.Asset
}

func (u unzipPipeline) Version() string {
	if u.asset == nil {
		return ""
	}
	return u.asset.Version()
}

func (u unzipPipeline) String() string { return unzip }

func (u unzipPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{
			stage:         u.Next(),
			targetVersion: u.Version(),
			assetName:     u.asset.Name,
		}
		if err := u.asset.Unzip(); err != nil {
			u.asset.Clean()
			msg.info = "解压失败: " + err.Error()
			msg.stage = nil
		}
		return msg
	}
}

func (u unzipPipeline) Next() upgradePipeline {
	return installPipelineStage{asset: u.asset}
}

type installPipelineStage struct {
	asset *github.Asset
}

func (i installPipelineStage) Version() string {
	if i.asset == nil {
		return ""
	}
	return i.asset.Version()
}

func (i installPipelineStage) String() string { return installing }

func (i installPipelineStage) Do() tea.Cmd {
	return func() tea.Msg {
		defer i.asset.Clean()
		msg := pipelineMsg{
			stage:          i.Next(),
			success:        true,
			currentVersion: displayVersion(consts.Version),
			targetVersion:  i.Version(),
			platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			assetName:      i.asset.Name,
		}
		if err := i.asset.Install(); err != nil {
			msg.success = false
			msg.info = fmt.Sprintf("安装失败: %s", err.Error())
			msg.stage = nil
			return msg
		}
		msg.info = fmt.Sprintf("升级完成: %s -> %s", msg.currentVersion, msg.targetVersion)
		return msg
	}
}

func (i installPipelineStage) Next() upgradePipeline {
	return nil
}

func displayVersion(v string) string {
	v = strings.TrimSpace(strings.TrimPrefix(v, "v"))
	if v == "" {
		return "unknown"
	}
	return v
}
