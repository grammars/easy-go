package logger

import (
	"log/slog"
	"testing"
)

func TestSetup(t *testing.T) {
	opt := CreateOption()
	opt.FileEnabled = true
	opt.Setup()
	slog.Info("测试输出", "名称", "Tom", "分数", 99)
}
