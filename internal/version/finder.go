package version

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/14 下午5:23
* @Package:
 */

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-yex/gvm/internal/core"
	"runtime"
	"sort"
)

type Finder struct {
	kind   Kind
	goos   string
	goarch string
	items  []*Version
}

// NewFinder creates a new Finder instance with sorted versions and applied options.
func NewFinder(items []*Version) *Finder {
	sort.Sort(Collection(items)) // Sort in ascending order.

	fdr := Finder{
		kind:   ArchiveKind,
		goos:   runtime.GOOS,
		goarch: runtime.GOARCH,
		items:  items,
	}

	return &fdr
}

func (fdr *Finder) Find(vname string) (*Version, error) {
	if vname == Latest {
		return fdr.findLatest()
	}

	for i := len(fdr.items) - 1; i >= 0; i-- {
		if fdr.items[i].String() == vname && fdr.items[i].match(fdr.goos, fdr.goarch) {
			return fdr.items[i], nil
		}
	}

	cs, err := NewConstraint(vname)
	if err != nil {
		return nil, fmt.Errorf("version not found %q [%s,%s]", vname, fdr.goos, fdr.goarch)
	}

	versionFound := false
	var (
		vsn []list.Item
		vs  []*Version
	)
	for i := len(fdr.items) - 1; i >= 0; i-- { // Prefer higher versions first.
		if cs.Check(fdr.items[i]) {
			versionFound = true
			if fdr.items[i].match(fdr.goos, fdr.goarch) {
				vs = append(vs, fdr.items[i])
				vsn = append(vsn, fdr.items[i])
			}
		}
	}
	if len(vsn) > 0 {
		type IndexProvider interface {
			Index() int
		}
		finalVersion, _ := core.NewSimpleListProgram(vsn, "select a fixed version to install", tea.WithAltScreen()).Run()
		if simpleModel, ok := finalVersion.(IndexProvider); ok {
			if simpleModel.Index() < 0 {
				return nil, fmt.Errorf("\n")
			}
			return vs[simpleModel.Index()], nil
		}
	}
	if versionFound {
		return nil, fmt.Errorf("package not found [%s,%s,%s]", string(fdr.kind), fdr.goos, fdr.goarch)
	}
	return nil, fmt.Errorf("version not found %q [%s,%s]", vname, fdr.goos, fdr.goarch)
}

// MustFind returns matched version or panics on error.
func (fdr *Finder) MustFind(vname string) *Version {
	v, err := fdr.Find(vname)
	if err != nil {
		panic(err)
	}
	return v
}

// Latest represents the current stable release.
const Latest = "latest"

func (fdr *Finder) findLatest() (*Version, error) {
	if len(fdr.items) == 0 {
		return nil, fmt.Errorf("version not found %q [%s,%s]", Latest, fdr.goos, fdr.goarch)
	}

	for i := len(fdr.items) - 1; i >= 0; i-- {
		if fdr.items[i].match(fdr.goos, fdr.goarch) {
			return fdr.items[i], nil
		}
	}
	return nil, fmt.Errorf("package not found [%s,%s,%s]", string(fdr.kind), fdr.goos, fdr.goarch)

}
