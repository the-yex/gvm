package pipeline

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-yex/gvm/internal/github"
)

const (
	check      = "check online version"
	download   = "download version"
	unzip      = "unzip package"
	installing = "install version"
)

type upgradePipeline interface {
	String() string
	Do() tea.Cmd
	Next() upgradePipeline
	Version() string
}

type pipelineMsg struct {
	stage upgradePipeline
	info  string
}

type searchPipeline struct {
	asset *github.Asset
}

func (s searchPipeline) Version() string {
	return s.asset.Version()
}
func (s searchPipeline) String() string { return check }
func (s searchPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{}
		releases, b, err := github.NewReleaseUpdater().CheckForUpdates()
		if err != nil {
			msg.info = err.Error()
			return msg
		}
		if !b {
			msg.info = "Currently the latest version"
			return msg
		}
		asset, err := releases.FindAsset()
		if err != nil {
			msg.info = err.Error()
			return msg
		}
		s.asset = asset
		msg.stage = s.Next()
		return msg
	}
}

func (s searchPipeline) Next() upgradePipeline {
	return downloadPipeline{asset: s.asset}
}

type downloadPipeline struct {
	asset *github.Asset
}

func (d downloadPipeline) Version() string {
	return d.asset.Version()

}
func (d downloadPipeline) String() string { return download }
func (d downloadPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{stage: d.Next()}
		_, err := d.asset.Download()
		if err != nil {
			d.asset.Clean()
			msg.info = fmt.Sprintf("download failed 【 %s 】", err.Error())
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
	return u.asset.Version()
}

func (u unzipPipeline) String() string { return unzip }
func (u unzipPipeline) Do() tea.Cmd {
	return func() tea.Msg {
		msg := pipelineMsg{stage: u.Next()}
		err := u.asset.Unzip()
		if err != nil {
			u.asset.Clean()
			msg.info = err.Error()
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
	return i.asset.Version()
}

func (i installPipelineStage) String() string { return installing }
func (i installPipelineStage) Do() tea.Cmd {
	return func() tea.Msg {
		defer i.asset.Clean()
		msg := pipelineMsg{stage: i.Next()}
		i.asset.Install()
		return msg
	}
}
func (i installPipelineStage) Next() upgradePipeline {
	return nil
}
