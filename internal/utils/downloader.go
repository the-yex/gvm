package utils

import (
	"errors"
	"fmt"
	progress2 "github.com/the-yex/gvm/internal/tui/progress"
	"io"
	"io/fs"
	"net/http"
	"os"
)

type SetSize func(int642 int64)

func Download(srcURL string, writer io.Writer, fn SetSize) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, srcURL, nil)
	if err != nil {
		return 0, fmt.Errorf("resource(%s) download failed ==> %s", srcURL, err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("resource(%s) download failed ==> %s", srcURL, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("URL %q is unreachable  ==> %d", srcURL, resp.StatusCode)
	}
	fn(resp.ContentLength)
	return io.Copy(writer, resp.Body)
}

func DownloadFile(srcURL, filename string, flag int, perm fs.FileMode) (int64, error) {
	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return 0, fmt.Errorf("resource(%s) download failed ==> %s", srcURL, err.Error())
	}
	defer f.Close()
	model := progress2.NewModel(nil)
	go func() {
		Download(srcURL, model.MultiWriter(f), model.SetSize)
		model.Quit()
	}()
	model.Start()
	if model.IsCancel() {
		os.Remove(filename)
		return 0, errors.New("download cancel")
	}
	return model.Size(), nil
}
