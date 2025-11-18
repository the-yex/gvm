package official

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/the-yex/gvm/internal/registry/base"
	"github.com/the-yex/gvm/internal/version"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/11 下午6:23
* @Package: official
 */

// Registry 实现了官方Go下载页面的解析逻辑
type Registry struct {
	base *base.Registry
}

// NewRegistry 创建一个新的官方Registry实例
func NewRegistry(mirrorUrl string) (*Registry, error) {
	baseRegistry, err := base.NewBaseRegistry(mirrorUrl)
	if err != nil {
		return nil, err
	}
	return &Registry{base: baseRegistry}, nil
}

// artifacts 从HTML表格中提取构件信息
func (r Registry) artifacts(table *goquery.Selection) (artifacts []version.ArtifactInfo) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		if strings.HasPrefix(href, "/") { // relative paths
			href = fmt.Sprintf("%s://%s%s", r.base.Url.Scheme, r.base.Url.Host, href)
		}
		artifacts = append(artifacts, version.ArtifactInfo{
			FileName:  td.Eq(0).Find("a").Text(),
			URL:       href,
			Kind:      version.Kind(strings.ToLower(td.Eq(1).Text())),
			OS:        version.OS(strings.ToLower(td.Eq(2).Text())),
			Arch:      version.ARCH(strings.ToLower(td.Eq(3).Text())),
			Size:      td.Eq(4).Text(),
			Checksum:  td.Eq(5).Text(),
			Algorithm: alg,
		})
	})
	return artifacts
}

// StableVersions 返回所有稳定版本
func (r Registry) StableVersions() (versions []*version.Version, err error) {
	versions = make([]*version.Version, 0, 10) // 预分配空间，减少扩容
	var divs *goquery.Selection
	if r.hasUnstableVersions() {
		divs = r.base.Doc.Find("#stable").NextUntil("#unstable")
	} else {
		divs = r.base.Doc.Find("#stable").NextUntil("#archive")
	}

	err = r.parseVersions(divs, &versions)
	return versions, err
}

// hasUnstableVersions 检查页面中是否存在不稳定版本
func (r Registry) hasUnstableVersions() bool {
	return r.base.Doc.Find("#unstable").Length() > 0
}

// UnstableVersions 返回所有不稳定版本
func (r Registry) UnstableVersions() (versions []*version.Version, err error) {
	versions = make([]*version.Version, 0, 5) // 预分配空间，减少扩容
	divs := r.base.Doc.Find("#unstable").NextUntil("#archive")
	err = r.parseVersions(divs, &versions)
	return versions, err
}

// ArchivedVersions 返回所有归档版本
func (r Registry) ArchivedVersions() (versions []*version.Version, err error) {
	versions = make([]*version.Version, 0, 20) // 预分配空间，减少扩容
	divs := r.base.Doc.Find("#archive").Find("div.toggle")
	err = r.parseVersions(divs, &versions)
	return versions, err
}

func (r Registry) parseVersions(divs *goquery.Selection, versions *[]*version.Version) error {
	var err error
	divs.EachWithBreak(func(i int, div *goquery.Selection) bool {
		versionName, ok := div.Attr("id")
		if !ok {
			return true
		}

		var v *version.Version
		if v, err = version.NewGoVersion(
			versionName,
			version.WithArtifacts(r.artifacts(div.Find("table").First())),
		); err != nil {
			return false
		}

		*versions = append(*versions, v)
		return true
	})
	return err
}

// AllVersions 返回所有版本（稳定版、不稳定版和归档版）
func (r Registry) AllVersions() (versions []*version.Version, err error) {
	versions = make([]*version.Version, 0, 50)
	// 获取稳定版本
	stableVersions, err := r.StableVersions()
	if err != nil {
		return nil, fmt.Errorf("获取稳定版本失败: %w", err)
	}
	versions = append(versions, stableVersions...)

	// 获取不稳定版本
	unstableVersions, err := r.UnstableVersions()
	if err != nil {
		return nil, fmt.Errorf("获取不稳定版本失败: %w", err)
	}
	versions = append(versions, unstableVersions...)

	// 获取归档版本
	archivedVersions, err := r.ArchivedVersions()
	if err != nil {
		return nil, fmt.Errorf("获取归档版本失败: %w", err)
	}
	versions = append(versions, archivedVersions...)
	return versions, nil

}
