package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)

	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = secureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

// GenerateRandomHex 生成指定长度的随机十六进制字符串
func GenerateRandomHex(length int) string {
	bytes := secureRandomBytes(length)
	return hex.EncodeToString(bytes)
}

// secureRandomBytes 生成指定长度的安全随机字节
func secureRandomBytes(length int) []byte {
	if length <= 0 {
		return []byte{}
	}
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 如果安全随机数生成失败，使用备用方案
		for i := range b {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
			b[i] = letterBytes[n.Int64()]
		}
	}
	return b
}

// GenerateRandomNumber 生成指定长度的随机数字字符串
func GenerateRandomNumber(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		result[i] = digits[n.Int64()]
	}

	return string(result)
}
