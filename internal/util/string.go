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
