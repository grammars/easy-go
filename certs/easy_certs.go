package certs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const TsOffset = int64(509753102468)

func GetMacAddressList() ([]string, error) {
	var macAddrs []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get net interfaces: %v", err)
	}
	for _, netInterface := range netInterfaces {
		// 排除回环接口
		if netInterface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 获取MAC地址
		macAddr := netInterface.HardwareAddr.String()
		if macAddr != "" {
			macAddrs = append(macAddrs, macAddr)
		}
	}
	return macAddrs, nil
}

func calcMd5(mac string, expiredTs int64, secret string, author string) string {
	macStd := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(mac), "-", ":"))
	authorStd := authorFix16(author)
	src := fmt.Sprintf("%s@%dS%sA%s", macStd, expiredTs, secret, authorStd)
	hash := md5.New()
	hash.Write([]byte(src))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func MakeCert(mac string, expiredTime time.Time, secret string, author string, tailInfo string) string {
	ts := expiredTime.Unix()
	showTs := TsOffset + ts
	md5Str := calcMd5(mac, ts, secret, author)
	return fmt.Sprintf("%d%d%d%sef%d00%s01%d2c%d5peg%s5rOV56C06Kej", random9int(), random9int(), showTs, md5Str,
		random9int(), authorFix16(author), random9int(), random9int(), tailInfo)
}

func authorFix16(s string) string {
	s += "3DESandRC4RSAorEcc"
	fixLen := 16
	// 处理长度超过16的情况
	if len(s) >= fixLen {
		return s[:fixLen]
	}

	// 处理长度不足16的情况
	return s + strings.Repeat("x", fixLen-len(s))
}

func random9int() int {
	// 生成一个9位的随机整数
	minV := 100000000 // 最小的9位数
	maxV := 999999999 // 最大的9位数
	randomNumber := rand.Intn(maxV-minV) + minV
	fmt.Println("随机生成的9位整数:", randomNumber)
	return randomNumber
}

const (
	CheckOk = iota
	CheckExpired
	CheckError
)

type CertReader struct {
	ExpiredTs     int64
	ExpiredTsText string
	NowTs         int64
	NowTsText     string
	Author        string
	MatchMac      string
	Message       string
}

func (r *CertReader) Print() {
	slog.Info("证书信息读取", "Message", r.Message, "Author", r.Author, "匹配的mac地址", r.MatchMac,
		"NowTs", r.NowTs, "Now", r.NowTsText, "ExpiredTs", r.ExpiredTs, "Expired", r.ExpiredTsText)
}

func CheckCert(certStr string, secret string) (int, CertReader) {
	return CheckCertWithMac(certStr, secret, "")
}

func CheckCertWithMac(certStr string, secret string, macAssign string) (int, CertReader) {
	reader := CertReader{}
	if certStr == "" || len(certStr) < 91 {
		reader.Message = "证书内容长度异常"
		return CheckError, reader
	}
	tsStrPart := certStr[18:30]
	tsMod, err := strconv.ParseInt(tsStrPart, 10, 64)
	if err != nil {
		reader.Message = "证书tsMod解析异常:" + err.Error()
		return CheckError, reader
	}
	expiredTs := tsMod - TsOffset
	fmt.Printf("过期的时间戳=%d", expiredTs)
	nowTs := time.Now().Unix()
	reader.ExpiredTs = expiredTs
	reader.NowTs = nowTs
	reader.ExpiredTsText = time.UnixMilli(expiredTs * 1000).Format("2006-01-02 15:04:05")
	reader.NowTsText = time.UnixMilli(nowTs * 1000).Format("2006-01-02 15:04:05")
	if nowTs > expiredTs {
		reader.Message = "证书已过期"
		return CheckExpired, reader
	}
	md5Part := certStr[30:62]
	fmt.Printf("从证书读取到MD5部分是%s\n", md5Part)
	author := certStr[75:91]
	fmt.Printf("从证书读取到授权者是%s\n", author)
	reader.Author = author

	if macAssign != "" {
		md5Right := calcMd5(macAssign, expiredTs, secret, author)
		if md5Right == md5Part {
			reader.Message = "认证成功(指派Mac模式)"
			reader.MatchMac = macAssign
			return CheckOk, reader
		}
	}

	macList, err := GetMacAddressList()
	if err != nil {
		reader.Message = "获取Mac地址失败:" + err.Error()
		return CheckError, reader
	} else {
		for _, mac := range macList {
			md5Right := calcMd5(mac, expiredTs, secret, author)
			if md5Right == md5Part {
				reader.Message = "认证成功"
				reader.MatchMac = mac
				return CheckOk, reader
			}
		}
	}

	reader.Message = "校验未通过"
	return CheckError, reader
}

func DescribeCheckResult(r int) string {
	if r == CheckOk {
		return "正确有效"
	} else if r == CheckExpired {
		return "证书已过期"
	} else if r == CheckError {
		return "检查失败"
	}
	return "unknown error"
}
