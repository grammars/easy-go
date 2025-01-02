package certs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func calcMd5(mac string, expiredTs int64, secret string) string {
	macStd := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(mac), "-", ":"))
	src := fmt.Sprintf("%s@%dS%s", macStd, expiredTs, secret)
	hash := md5.New()
	hash.Write([]byte(src))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func MakeCert(mac string, expiredTime time.Time, secret string) string {
	ts := expiredTime.Unix()
	showTs := TsOffset + ts
	md5Str := calcMd5(mac, ts, secret)
	return fmt.Sprintf("%d%d%d%sef%d01%d2c%d5peg5rOV56C06Kej", random9int(), random9int(), showTs, md5Str, random9int(), random9int(), random9int())
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

func CheckCert(certStr string, secret string) int {
	if certStr == "" || len(certStr) < 62 {
		return CheckError
	}
	tsStrPart := certStr[18:30]
	tsMod, err := strconv.ParseInt(tsStrPart, 10, 64)
	if err != nil {
		return CheckError
	}
	expiredTs := tsMod - TsOffset
	//fmt.Printf("过期的时间戳=%d", expiredTs)
	nowTs := time.Now().Unix()
	if nowTs > expiredTs {
		return CheckExpired
	}
	md5Part := certStr[30:62]

	macList, err := GetMacAddressList()
	if err != nil {
		return CheckError
	} else {
		for _, mac := range macList {
			md5Right := calcMd5(mac, expiredTs, secret)
			if md5Right == md5Part {
				return CheckOk
			}
		}
	}

	return CheckError
}
