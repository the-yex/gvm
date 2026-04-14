package github

import (
	"fmt"
	"os"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/29 下午2:51
* @Package:
 */
func TestCheck(t *testing.T) {
	if os.Getenv("GVM_NETWORK_TEST") != "1" {
		t.Skip("skipping network-dependent release test (set GVM_NETWORK_TEST=1 to enable)")
	}

	i, y, err := NewReleaseUpdater().CheckForUpdates()
	if err != nil {
		t.Fatalf("check for updates failed: %v", err)
	}
	if !y {
		t.Log("no new release found")
		return
	}

	assert, err := i.FindAsset()
	if err != nil {
		t.Fatalf("find asset failed: %v", err)
	}
	defer assert.Clean()
	os.Setenv("http_proxy", "127.0.0.1:7890")
	os.Setenv("https_proxy", "127.0.0.1:7890")
	fmt.Println(assert.Download())
	fmt.Println(assert.Unzip())
	assert.Install()
	fmt.Println(assert)
}
