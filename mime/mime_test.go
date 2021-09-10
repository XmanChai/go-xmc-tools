package mime

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetFileType(t *testing.T) {
	fs, err := os.Open("gopher")
	if err != nil {
		t.Logf("open file error:%v", err)
	}
	src, err := ioutil.ReadAll(fs)
	t.Log(GetFileType(src[:10]))
}
