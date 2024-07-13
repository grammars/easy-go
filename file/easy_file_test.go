package file

import (
	"fmt"
	"github.com/grammars/easy-go/sugar"
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
