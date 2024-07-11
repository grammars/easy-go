package file

import (
	"fmt"
	ego "github.com/grammars/easy-go"
	"testing"
)

func TestCommon(t *testing.T) {
	fmt.Println("当前可执行文件所在的目录：", ego.File.GetExeDir())
	existsCheck := func(path string) {
		b := ego.File.Exists(path)
		t.Logf("%s %s", path, ego.Return(b, "存在", "不存在"))
	}
	existsCheck("C:\\oem8.log")
	existsCheck("C:\\nothing")
	existsCheck("D:\\Java")
	existsCheck(ego.File.GetPathRelExe("TempInTemp.txt"))
	md5, err := ego.File.Md5Hex("C:\\oem8.log")
	if err != nil {
		t.Errorf("MD5计算失败：%v", err)
	} else {
		t.Logf("MD5计算得到:%s", md5)
	}

}
