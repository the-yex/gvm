package base

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/the-yex/gvm/internal/version"
	"net/http"
	"net/url"
	"time"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/12 上午10:30
* @Package:
 */

type Registry struct {
	url string
	Url *url.URL
	Doc *goquery.Document
}

func NewBaseRegistry(mirrorUrl string) (*Registry, error) {
	r := &Registry{url: mirrorUrl}
	u, err := url.Parse(mirrorUrl)
	if err != nil {
		return nil, err
	}
	r.Url = u
	err = r.loadDocument()
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (r *Registry) loadDocument() error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, r.url, nil)
	if err != nil {
		return fmt.Errorf("invalid request for %s: %w", r.url, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", r.url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned %s", r.url, resp.Status)
	}

	r.Doc, err = goquery.NewDocumentFromReader(resp.Body)
	return err
}

func (r *Registry) StableVersions() (versions []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

func (r *Registry) UnstableVersions() (versions []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

func (r *Registry) ArchivedVersions() (versions []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

func (r *Registry) AllVersions() (versions []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}
