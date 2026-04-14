package autoindex

import (
	"fmt"
	"os"
	"testing"
	"time"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/12 上午10:02
* @Package:
 */
const OfficialDownloadPageURL = "https://mirrors.ustc.edu.cn/golang/"

func Test_Parse(t *testing.T) {
	if os.Getenv("GVM_NETWORK_TEST") != "1" {
		t.Skip("skipping network-dependent registry parsing test (set GVM_NETWORK_TEST=1 to enable)")
	}

	r, err := NewRegistry(OfficialDownloadPageURL, 10*time.Second)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}
	versions, err := r.AllVersions()
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}
	for _, version := range versions {
		fmt.Println(version.String())
		for _, artifact := range version.Artifacts {
			fmt.Println(artifact.OS, artifact.Arch, artifact.Kind, artifact.FileName, artifact.URL)
			fmt.Println()
		}
	}
}
