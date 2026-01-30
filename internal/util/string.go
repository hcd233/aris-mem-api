package util

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/mozillazg/go-pinyin"
)

// ToDataURL 将文件转换为 data URL
//
//	@param contentType 文件类型
//	@param bytes
//	@return string
//	@author centonhuang
//	@update 2025-11-13 17:49:49
func ToDataURL(contentType string, bytes []byte) string {
	base64Data := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
}

// DecodeBase64OrDataURL 解码 base64 字符串或 Data URL
// 支持两种格式：
// 1. 纯 base64 字符串
// 2. Data URL 格式 (data:image/png;base64,...)
//
//	@param input base64 字符串或 Data URL
//	@return []byte 解码后的字节数组
//	@return string MIME 类型（如果可以提取）
//	@return error
//	@author centonhuang
//	@update 2026-01-30 00:00:00
func DecodeBase64OrDataURL(input string) (bytes []byte, mimeType string, err error) {
	// 检查是否是 Data URL 格式
	if strings.HasPrefix(input, "data:") {
		// 查找 base64, 标记
		parts := strings.SplitN(input, ",", 2)
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid data URL format")
		}

		// 提取 MIME 类型 (data:image/png;base64 -> image/png)
		header := parts[0]
		if strings.Contains(header, ";") {
			mimeType = strings.TrimPrefix(strings.Split(header, ";")[0], "data:")
		}

		// 使用第二部分（纯 base64 数据）
		input = parts[1]
	}

	// 解码 base64
	bytes, err = base64.StdEncoding.DecodeString(input)

	return bytes, mimeType, nil
}

// GenerateSlug 生成Slug
//
//	@param s
//	@return string
//	@author centonhuang
//	@update 2026-01-28 21:47:50
func GenerateSlug(s string) string {
	var parts []string
	var buffer []rune

	a := pinyin.NewArgs()

	for _, r := range s {
		if (unicode.IsLetter(r) || unicode.IsDigit(r)) && r < 128 { // ASCII letter (English)
			buffer = append(buffer, r)
		} else if unicode.Is(unicode.Han, r) { // Chinese character
			// Flush buffer if it contains English letters
			if len(buffer) > 0 {
				parts = append(parts, strings.ToLower(string(buffer)))
				buffer = buffer[:0]
			}

			// Process Chinese character
			py := pinyin.LazyPinyin(string(r), a)
			if len(py) > 0 {
				parts = append(parts, py[0])
			}
		}
		// Ignore all other characters (special symbols, spaces, etc.)
	}

	// Flush remaining buffer
	if len(buffer) > 0 {
		parts = append(parts, strings.ToLower(string(buffer)))
	}

	parts = append(parts, uuid.New().String()[:8])

	return strings.Join(parts, "-")
}
