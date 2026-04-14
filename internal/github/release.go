// Copyright (c) 2019 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/mholt/archiver/v3"
	"github.com/the-yex/gvm/internal/consts"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Release represents a software version release.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// FindAsset 当前版本是否支持该系统
func (r *Release) FindAsset() (*Asset, error) {
	for _, asset := range r.Assets {
		if strings.Contains(asset.Name, runtime.GOOS) && strings.Contains(asset.Name, runtime.GOARCH) {
			return &asset, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("current version not support %s-%s", runtime.GOOS, runtime.GOARCH))
}

// Asset contains downloadable resource files.
type Asset struct {
	Name               string `json:"name"`
	tempDir            string `json:"temp_dir"`
	ContentType        string `json:"content_type"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ProgressFunc func(written, total int64)

func (a *Asset) Version() string {
	re := regexp.MustCompile(`\d+\.\d+\.\d+`)
	return re.FindString(a.Name)
}

// IsCompressedFile checks if the file is in compressed format.
func (a *Asset) IsCompressedFile() bool {
	return a.ContentType == "application/zip" || a.ContentType == "application/x-gzip"
}

// Download saves the remote resource to local file with progress support.
func (a *Asset) Download() (size int64, err error) {
	return a.DownloadWithProgress(nil)
}

// DownloadWithProgress saves the remote resource to a local file and reports download progress.
func (a *Asset) DownloadWithProgress(onProgress ProgressFunc) (size int64, err error) {
	url := a.BrowserDownloadURL
	if source := os.Getenv("GVM_SOURCE"); source != "" {
		url = strings.Replace(url, "gitlab", source, 1)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "gvm")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	tmpDir, err := os.MkdirTemp(consts.GVM_HOME, strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return 0, err
	}
	a.tempDir = tmpDir
	f, err := os.OpenFile(filepath.Join(tmpDir, a.Name), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	if onProgress != nil {
		onProgress(0, resp.ContentLength)
	}
	pw := &progressWriter{
		total:      resp.ContentLength,
		onProgress: onProgress,
	}
	size, err = io.Copy(io.MultiWriter(f, pw), resp.Body)
	if err != nil {
		return 0, err
	}
	if onProgress != nil {
		onProgress(size, resp.ContentLength)
	}
	return size, nil
}

func (a *Asset) Unzip() error {
	return archiver.Unarchive(filepath.Join(a.tempDir, a.Name), a.tempDir)
}
func (a *Asset) Clean() {
	if a.tempDir != "" {
		os.RemoveAll(a.tempDir)
	}
}
func (a *Asset) Install() error {
	src := filepath.Join(a.tempDir, "gvm")
	dst := filepath.Join(consts.GVM_HOME, "gvm")
	tmp := dst + ".tmp"

	if err := a.mv(src, tmp); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0755); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
func (a *Asset) mv(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

type progressWriter struct {
	written    int64
	total      int64
	onProgress ProgressFunc
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.written += int64(n)
	if pw.onProgress != nil {
		pw.onProgress(pw.written, pw.total)
	}
	return n, nil
}

// ReleaseUpdater handles version update checks and operations.
type ReleaseUpdater struct {
}

// NewReleaseUpdater creates a release update handler instance.
func NewReleaseUpdater() *ReleaseUpdater {
	return new(ReleaseUpdater)
}

// CheckForUpdates verifies if newer version exists. https://api.github.com/repos/the-yex/gvm/releases/latest
func (up ReleaseUpdater) CheckForUpdates() (rel *Release, yes bool, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", consts.AUTHOR, consts.NAME)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, errors.New("unexpected status code: " + resp.Status)
	}

	var latest Release
	if err = json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return nil, false, err
	}

	latestVersion, err := semver.NewVersion(latest.TagName)
	if err != nil {
		return nil, false, err
	}
	currentVersion, err := semver.NewVersion(strings.TrimPrefix(consts.Version, "v"))
	if err != nil {
		return &latest, true, nil
	}
	if latestVersion.GreaterThan(currentVersion) {
		return &latest, true, nil
	}
	return &latest, false, nil
}
