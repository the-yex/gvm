package core

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/utils"
	"io"
	"os"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/25 下午4:11
* @Package:
 */

var (
	UninstallVersion   func(version string) error
	InstallVersion     func(version string) error
	MultiWriterInstall func(version any, writer io.Writer, fn func(int642 int64)) error
)

func SwitchVersion(versionDir string) error {
	os.Remove(consts.GO_ROOT)
	_, err := os.Stat(versionDir)
	if err != nil {
		return err
	}
	if err = utils.Symlink(versionDir, consts.GO_ROOT); err != nil {
		return err
	}
	return nil
}

var (
	NewSpinnerProgram    func(options ...tea.ProgramOption) *tea.Program
	NewSimpleListProgram func(items []list.Item, title string, options ...tea.ProgramOption) *tea.Program
)
