package official

import (
	"fmt"
	"testing"
	"time"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/12 上午10:02
* @Package:
 */
const OfficialDownloadPageURL = "https://golang.google.cn/dl/"

func Test_Parse(t *testing.T) {
	r, err := NewRegistry(OfficialDownloadPageURL, 10*time.Second)
	if nil != err {
		panic(err)
	}
	versions, err := r.ArchivedVersions()
	if err != nil {
		panic(err)
	}
	for _, version := range versions {
		fmt.Println(version.String())
	}
}
