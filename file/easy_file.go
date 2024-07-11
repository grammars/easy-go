package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var curExeDir string

// GetExeDir 获取当前可执行文件所在的目录
func GetExeDir() string {
	if curExeDir == "" {
		exePath, err := os.Executable()
		if err != nil {
			panic("cannot get executable path: " + err.Error())
		}
		curExeDir = filepath.Dir(exePath)
	}
	return curExeDir
}

// GetPathRelExe 获取相对于可执行文件的绝对路径
func GetPathRelExe(relPath string) string {
	return filepath.Join(GetExeDir(), relPath)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Md5Hex(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件时发生错误", err)
		}
	}(file) // 确保在函数结束时关闭文件
	// 创建MD5哈希对象
	hash := md5.New()
	// 从文件中读取数据并写入哈希对象
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	// 获取哈希值的字节切片
	md5Bytes := hash.Sum(nil)
	// 将字节切片转换为十六进制字符串
	return hex.EncodeToString(md5Bytes), nil
}
