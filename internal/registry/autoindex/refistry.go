package autoindex

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/the-yex/gvm/internal/registry/base"
	"github.com/the-yex/gvm/internal/registry/internal"
	"github.com/the-yex/gvm/internal/version"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/11 下午6:23
* @Package:
 */

type Registry struct {
	base *base.Registry
}

func NewRegistry(mirrorUrl string) (*Registry, error) {
	baseRegistry, err := base.NewBaseRegistry(mirrorUrl)
	if err != nil {
		return nil, err
	}
	return &Registry{base: baseRegistry}, nil
}

func (r Registry) StableVersions() (versions []*version.Version, err error) {
	return r.AllVersions()
}

func (r Registry) UnstableVersions() (versions []*version.Version, err error) {
	return r.AllVersions()
}

func (r Registry) ArchivedVersions() (versions []*version.Version, err error) {
	return r.AllVersions()
}

func (r Registry) AllVersions() (versions []*version.Version, err error) {
	anchors := r.base.Doc.Find("pre").Find("a")
	items := make([]*internal.GoFileItem, 0, anchors.Length())

	anchors.Each(func(j int, anchor *goquery.Selection) {
		href := anchor.AttrOr("href", "")
		if !strings.HasPrefix(href, "go") || strings.HasSuffix(href, "/") {
			return
		}

		var size string
		if fields := strings.Fields(strings.TrimSpace(anchor.Nodes[0].NextSibling.Data)); len(fields) > 0 {
			size = fields[len(fields)-1]
		}

		items = append(items, &internal.GoFileItem{
			FileName: anchor.Text(),
			URL:      r.base.Url.JoinPath(href).String(),
			Size:     size,
		})
	})
	if len(items) == 0 {
		return nil, nil
	}
	return internal.Convert2Versions(items)
}
