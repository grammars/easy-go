package certs

import (
	"fmt"
	"github.com/grammars/easy-go/file"
	"log/slog"
	"os"
	"testing"
	"time"
)

const DEMO_SECRET = "123456"

func TestCertMake(t *testing.T) {
	authMac := "  FC-34-97-E2-DE-B0  " // 授权MAC
	expiredStr := "2025-09-01 12:00:00"
	layout := "2006-01-02 15:04:05"

	// 解析日期字符串
	expiredTime, err := time.Parse(layout, expiredStr)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}

	certContent := MakeCert(authMac, expiredTime, DEMO_SECRET, "好家伙", "真厉害")
	certFileName := "../build/auth.cert"
	// 使用os.Create创建文件
	certFile, err := os.Create(certFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(certFile *os.File) {
		err := certFile.Close()
		if err != nil {
			slog.Error("Error closing file", "err", err.Error())
		}
	}(certFile) // 确保在函数结束时关闭文件

	// 使用WriteString方法写入文本
	_, err = certFile.WriteString(certContent)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Text written to file successfully")
}

func TestCertCheck(t *testing.T) {
	txt, err := file.ReadString(`../build/auth.cert`)
	if err != nil {
		slog.Error("Error reading file", "err", err.Error())
		return
	}
	cr := CheckCert(txt, DEMO_SECRET)
	if cr == CheckOk {
		fmt.Println("证书验证成功")
	} else if cr == CheckExpired {
		fmt.Println("证书已过期")
	} else if cr == CheckError {
		fmt.Println("证书验证错误")
	} else {
		fmt.Println("证书验证失败")
	}
}
