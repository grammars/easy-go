package ego

import "github.com/grammars/easy-go/pkg/file"

func EchoFile() file.MyFile {
	println("快捷文件操作")
	var f = file.MyFile{}
	f.Name = "老🐢王"
	f.Size = 16665
	return f
}
