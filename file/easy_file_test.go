package file

import (
	"fmt"
	"github.com/grammars/easy-go/sugar"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCommon(t *testing.T) {
	fmt.Println("当前可执行文件所在的目录：", GetExeDir())
	existsCheck := func(path string) {
		b := Exists(path)
		t.Logf("%s %s", path, sugar.ReturnIf(b, "存在", "不存在"))
	}
	existsCheck("C:\\oem8.log")
	existsCheck("C:\\nothing")
	existsCheck("D:\\Java")
	existsCheck(GetPathRelExe("TempInTemp.txt"))
	md5, err := Md5Hex("C:\\oem8.log")
	if err != nil {
		t.Errorf("MD5计算失败：%v", err)
	} else {
		t.Logf("MD5计算得到:%s", md5)
	}
}

func TestConvertAbsPath(t *testing.T) {
	var do = func(path string, expectPath string) {
		abs := ConvAbsPathRelExe(path)
		if abs != expectPath {
			t.Errorf("转化绝对路径❌ 输入路径=%s 得到绝对路径=%s 预期=%s", path, abs, expectPath)
		} else {
			t.Logf("转化绝对路径✅ 得到绝对路径=%s ", abs)
		}
	}
	currentOS := runtime.GOOS
	switch currentOS {
	case "windows":
		do("D:\\Java", "D:\\Java")
		do("D:\\Go\\xyz", "D:\\Go\\xyz")
	case "linux":
		do("/home/demo/", "/home/demo/")
		do("/home/lab/juice", "/home/lab/juice")
	}
	do("application.yml", filepath.Join(GetExeDir(), "application.yml"))
}
