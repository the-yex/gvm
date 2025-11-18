package pkg

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/12 下午1:54
* @Package:
 */

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/core"
	"github.com/the-yex/gvm/internal/registry"
	"github.com/the-yex/gvm/internal/utils"
	"github.com/the-yex/gvm/internal/version"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/11 下午3:45
* @Package:
 */

type VManager interface {
	// List 列出所有版本号
	List(kind consts.VersionKind) ([]*version.Version, error)
	// Install 安装指定版本号到本地
	Install(versionName string) error
	// Uninstall 从本地卸载指定版本号func DownloadFile(srcURL, filename string, flag int, perm fs.FileMode) (int64, error) {
	Uninstall(versionName string) error
}

type ManagerOption struct {
	WithLocal bool
}

func init() {
	core.MultiWriterInstall = remote{}.MultiWriterInstall
	core.UninstallVersion = local{}.UninstallDir
	core.InstallVersion = remote{}.Install
}
func WithLocal() func(option *ManagerOption) {
	return func(option *ManagerOption) {
		option.WithLocal = true
	}
}

// NewVManager
// r
func NewVManager(r bool, opts ...func(option *ManagerOption)) VManager {
	opt := &ManagerOption{}
	for _, o := range opts {
		o(opt)
	}
	if r {
		return &remote{withLocal: opt.WithLocal}
	}
	return &local{}
}

type local struct {
}

func (l local) List(kind consts.VersionKind) ([]*version.Version, error) {
	goRoots := viper.GetStringSlice(consts.CONFIG_GOROOT)
	currentVersion := l.currentUsedVersion()
	if len(goRoots) == 0 {
		return nil, nil
	}
	var versions []*version.Version
	for _, root := range goRoots {
		versionDirs, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, versionDir := range versionDirs {
			if !versionDir.IsDir() {
				continue
			}
			v, err := version.NewVersion(strings.TrimPrefix(versionDir.Name(), "go"))
			if err != nil || v == nil {
				continue
			}
			v.CurrentUsed = v.String() == currentVersion
			v.Installed = true
			v.Path = root
			v.DirName = versionDir.Name()
			versions = append(versions, v)
		}
	}
	return versions, nil
}
func (l local) Install(versionName string) error {
	return errors.New("not support")
}
func (l local) currentUsedVersion() string {
	p, _ := os.Readlink(consts.GO_ROOT)
	versionName := filepath.Base(p)
	return strings.TrimPrefix(versionName, "go")
}
func (l local) currentUsedVersionDir() string {
	p, _ := os.Readlink(consts.GO_ROOT)
	return p
}

func (l local) Uninstall(versionName string) error {
	if versionName == l.currentUsedVersion() {
		return fmt.Errorf("cannot uninstall version %s: it is currently in use\n", versionName)
	}
	for _, root := range viper.GetStringSlice(consts.CONFIG_GOROOT) {
		targetDir := filepath.Join(root, fmt.Sprintf("go%s", versionName))
		fmt.Println(targetDir)
		if finfo, err := os.Stat(targetDir); err != nil || !finfo.IsDir() {
			continue
		}
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("uninstall failed: %s\n", err.Error())
		}
		return nil
	}
	return fmt.Errorf("version %q is not installed\n", versionName)
}
func (l local) UninstallDir(versionDir string) error {
	if versionDir == l.currentUsedVersion() {
		return fmt.Errorf("cannot uninstall version %s: it is currently in use\n", versionDir)
	}
	if finfo, err := os.Stat(versionDir); err != nil || !finfo.IsDir() {
		return fmt.Errorf("version %q is not installed\n", versionDir)
	}
	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("uninstall failed: %s\n", err.Error())
	}
	return nil
}

type remote struct {
	withLocal bool
}

func (r remote) mergeInstalled(remoteVers []*version.Version, localVers []*version.Version) {
	m := make(map[string]*version.Version)
	for _, v := range localVers {
		m[v.Original()] = v // 用原始版本号做 key
	}
	for _, v := range remoteVers {
		v.Path = consts.VERSION_DIR
		if lv, ok := m[v.Original()]; ok {
			v.Installed = true
			v.CurrentUsed = lv.CurrentUsed
			v.Path = lv.Path
			v.DirName = lv.DirName
		}
	}
}

func (r remote) List(kind consts.VersionKind) (versions []*version.Version, err error) {
	p := core.NewSpinnerProgram(tea.WithAltScreen())
	wg := sync.WaitGroup{}
	wg.Go(func() {
		p.Run()
	})
	rg, err := registry.NewRegistry()
	if err != nil {
		return nil, err
	}
	switch kind {
	case consts.Stable:
		versions, err = rg.StableVersions()
	case consts.Unstable:
		versions, err = rg.UnstableVersions()
	case consts.Archived:
		versions, err = rg.ArchivedVersions()
	default:
		versions, err = rg.AllVersions()
	}
	if r.withLocal {
		installVersions, _ := local{}.List(kind)
		r.mergeInstalled(versions, installVersions)
	}
	p.Send(tea.Quit())
	wg.Wait()
	return versions, err
}

func (r remote) Install(versionName string) error {
	versions, err := (&remote{withLocal: false}).List(consts.All)
	if err != nil {
		return err
	}
	v, err := version.NewFinder(versions).Find(versionName)
	if err != nil {
		return err
	}
	if LocalInstalled(v.String()) != nil {
		return fmt.Errorf("%s has already been installed\n", v.String())
	}
	err = v.Install()
	if nil != err {
		return err
	}
	v.Path = consts.VERSION_DIR
	v.DirName = fmt.Sprintf("go%s", v.String())
	fmt.Println(v.LocalDir())
	return SwitchVersion(v.LocalDir())
}

func (r remote) MultiWriterInstall(item any, writer io.Writer, fn func(int642 int64)) error {
	v, ok := item.(*version.Version)
	if !ok {
		return errors.New("invalid version")
	}
	if LocalInstalled(v.String()) != nil {
		return fmt.Errorf("%s has already been installed\n", v.String())
	}
	artifact, err := v.FindArtifact()
	if err != nil {
		return err
	}
	err = artifact.MultiWriterInstall(v.String(), writer, fn)
	if nil != err {
		return err
	}
	return nil
}

func (r remote) Uninstall(version string) error {
	//TODO implement me
	return errors.New("not support")
}

/*
*
软连接go指定版本的本地目录
*/
func SwitchVersion(versionDir string) error {
	os.Remove(consts.GO_ROOT)
	_, err := os.Stat(versionDir)
	if err != nil {
		return err
	}
	if err = utils.Symlink(versionDir, consts.GO_ROOT); err != nil {
		return err
	}
	if output, err := exec.Command(filepath.Join(consts.GO_ROOT, "bin", "go"), "version").Output(); err == nil {
		fmt.Printf("Now using %s", strings.TrimPrefix(string(output), "go version "))
	}
	return nil
}

func LocalInstalled(versionName string) *version.Version {
	if versionName == "" {
		return nil
	}
	installVersions, _ := local{}.List(consts.All)
	for _, installVersion := range installVersions {
		if installVersion.String() == versionName {
			return installVersion
		}
	}
	return nil
}
