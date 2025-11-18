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

package internal

import (
	"github.com/the-yex/gvm/internal/utils"
	"github.com/the-yex/gvm/internal/version"
	"regexp"
	"strings"
)

var goVersionPattern = regexp.MustCompile(`^go(\d+(?:\.\d+){0,2}(?:[a-z]+\d+)?)`)

type GoFileItem struct {
	FileName string
	URL      string
	Size     string
}

func (item GoFileItem) getGoVersion() string {
	matches := goVersionPattern.FindStringSubmatch(item.FileName)
	if len(matches) > 1 {
		return matches[1] // e.g. 1.23, 1.23.1, 1.24rc1
	}
	return ""
}

func (item GoFileItem) isSHA256File() bool {
	return strings.HasSuffix(item.FileName, ".sha256")
}

func (item GoFileItem) isPackageFile() bool {
	return strings.HasSuffix(item.FileName, ".tar.gz") ||
		strings.HasSuffix(item.FileName, ".pkg") ||
		strings.HasSuffix(item.FileName, ".zip") ||
		strings.HasSuffix(item.FileName, ".msi")
}

func (item GoFileItem) getKind() version.Kind {
	if strings.HasSuffix(item.FileName, ".src.tar.gz") {
		return version.SourceKind
	}
	if strings.HasSuffix(item.FileName, ".tar.gz") || strings.HasSuffix(item.FileName, ".zip") {
		return version.ArchiveKind
	}
	if strings.HasSuffix(item.FileName, ".pkg") || strings.HasSuffix(item.FileName, ".msi") {
		return version.InstallerKind
	}
	return "Unknown"
}

var osMapping = []struct {
	key string
	val version.OS
}{
	{"linux", version.Linux},
	{"darwin", version.MacOS},
	{"windows", version.Windows},
	{"freebsd", version.FreeBSD},
	{"netbsd", version.NetBSD},
	{"openbsd", version.OpenBSD},
	{"solaris", version.Solaris},
	{"plan9", version.Plan9},
	{"aix", version.AIX},
	{"dragonfly", version.Dragonfly},
	{"illumos", version.Illumos},
}

func (item GoFileItem) getOS() version.OS {
	for _, kv := range osMapping {
		if strings.Contains(item.FileName, kv.key) {
			return kv.val
		}
	}
	return ""
}

var archMapping = []struct {
	Key string
	Val version.ARCH
}{
	{"-loong64.", version.LOONGARCH64},
	{"-loong32", version.LOONGARCH32},
	{"-riscv64.", version.RISCV64},
	{"-s390x.", version.S390X},
	{"-mips64le.", version.MIPS64LE},
	{"-mips64.", version.MIPS64},
	{"-mipsle.", version.MIPSLE},
	{"-mips.", version.MIPS},
	{"-ppc64le.", version.PPC64LE},
	{"-ppc64.", version.PPC64},
	{"-arm64.", version.ARM64},
	{"-armv6l.", version.ARMV6},
	{"-arm.", version.ARMV6}, // 注意顺序：放在 arm64 之后
	{"-amd64.", version.X8664},
	{"-386.", version.X86},
}

func (item GoFileItem) getArch() version.ARCH {
	for _, m := range archMapping {
		if strings.Contains(item.FileName, m.Key) {
			return m.Val
		}
	}
	return version.UnknownARCH
}

func Convert2Versions(items []*GoFileItem) (versions []*version.Version, err error) {
	artifactInfos := make(map[string][]version.ArtifactInfo, 20)

	for _, item := range items {
		ver := item.getGoVersion()
		if ver == "" {
			continue
		}
		if _, ok := artifactInfos[ver]; !ok {
			artifactInfos[ver] = make([]version.ArtifactInfo, 0, 20)
		}

		if item.isPackageFile() {
			artifactInfos[ver] = append(artifactInfos[ver], version.ArtifactInfo{
				FileName: item.FileName,
				URL:      item.URL,
				Kind:     item.getKind(),
				OS:       item.getOS(),
				Arch:     item.getArch(),
				Size:     item.Size,
			})
		} else if item.isSHA256File() {
			for index, ppkg := range artifactInfos[ver] {
				if ppkg.FileName == strings.TrimSuffix(item.FileName, ".sha256") {
					artifactInfos[ver][index].Algorithm = string(utils.SHA256)
					artifactInfos[ver][index].ChecksumURL = item.URL
				}
			}
		}
	}

	versions = make([]*version.Version, 0, len(artifactInfos))
	for vname, infos := range artifactInfos {
		v, err := version.NewGoVersion(vname, version.WithArtifacts(infos))
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}
