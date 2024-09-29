package str

import "testing"

func TestIsBlank(t *testing.T) {
	if false == IsBlank(" ") {
		t.Error("空格判定blank错误")
	}
	if false == IsBlank("") {
		t.Error("纯空判定blank错误")
	}
	if true == IsBlank("有货") {
		t.Error("纯空判定blank错误")
	}
}

func TestBytesToHex(t *testing.T) {
	t.Logf("16进制显示：%s", BytesToHex([]byte{1, 12, 10}))
}
