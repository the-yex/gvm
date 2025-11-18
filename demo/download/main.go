package main

import (
	"fmt"
	"github.com/the-yex/gvm/internal/utils"
	"os"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/26 下午4:40
* @Package:
 */
func main() {
	size, err := utils.DownloadFile("https://golang.google.cn/dl/go1.24.6.darwin-arm64.tar.gz", "./go1.24.6.darwin-arm64.tar.gz", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(size)
}
