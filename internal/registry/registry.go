package registry

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/registry/autoindex"
	"github.com/the-yex/gvm/internal/registry/fancyindex"
	"github.com/the-yex/gvm/internal/registry/official"
	"github.com/the-yex/gvm/internal/version"
	"maps"
	"slices"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/9 下午2:21
* @Package:
 */

// MirrorType 表示 Golang 下载镜像类型
type MirrorType string

const (
	Official   MirrorType = "Official"   // https://go.dev/dl/
	CNOfficial MirrorType = "CNOfficial" // https://golang.google.cn/dl/
	Aliyun     MirrorType = "Aliyun"     // https://mirrors.aliyun.com/golang/
	HUST       MirrorType = "HUST"       // https://mirrors.hust.edu.cn/golang/
	NJU        MirrorType = "NJU"        // https://mirrors.nju.edu.cn/golang/
	USTC       MirrorType = "USTC"       // https://mirrors.ustc.edu.cn/golang/
)

type Mirror struct {
	Name string
	URL  string
}

var Mirrors = map[string]MirrorType{
	"https://go.dev/dl/":                  Official,
	"https://golang.google.cn/dl/":        CNOfficial,
	"https://mirrors.aliyun.com/golang/":  Aliyun,
	"https://mirrors.hust.edu.cn/golang/": HUST,
	"https://mirrors.nju.edu.cn/golang/":  NJU,
	"https://mirrors.ustc.edu.cn/golang/": USTC,
}

type Registry interface {
	// StableVersions 返回所有稳定版本
	StableVersions() (versions []*version.Version, err error)
	// UnstableVersions 返回所有不稳定版本
	UnstableVersions() (versions []*version.Version, err error)
	// ArchivedVersions 返回所有归档版本
	ArchivedVersions() (versions []*version.Version, err error)
	// AllVersions 返回所有版本（稳定版、不稳定版和归档版）
	AllVersions() (versions []*version.Version, err error)
}

func NewRegistry() (Registry, error) {
	mirrorUrl := viper.GetString(consts.CONFIG_MIRROR)
	mirror, exist := Mirrors[mirrorUrl]
	if !exist {
		supported := slices.SortedStableFunc(maps.Keys(Mirrors), func(s string, s2 string) int {
			return strings.Compare(s, s2)
		})
		return nil, fmt.Errorf(
			"无效的配置 URL: %q\n支持的 URL 列表如下:\n  %s",
			mirrorUrl,
			strings.Join(supported, "\n  "),
		)
	}
	switch mirror {
	case Official, CNOfficial:
		return official.NewRegistry(mirrorUrl)
	case USTC:
		return autoindex.NewRegistry(mirrorUrl)
	default:
		return fancyindex.NewRegistry(mirrorUrl)
	}
}
