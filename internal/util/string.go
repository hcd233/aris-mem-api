package util

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/mozillazg/go-pinyin"
	"github.com/samber/lo"
	"golang.org/x/net/html"
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

// ToThumbnailURL 将图片 URL 转换为缩略图 URL
//
//	@param imageURL
//	@return string
//	@author centonhuang
//	@update 2026-01-30 15:05:04
func ToThumbnailURL(imageURL string) string {
	u := lo.Must1(url.Parse(imageURL))
	q := u.RawQuery
	if q != "" {
		u.RawQuery = fmt.Sprintf("%s&imageView2/1/q/%d", q, constant.DefaultThumbnailQuality)
	} else {
		u.RawQuery = fmt.Sprintf("imageView2/1/q/%d", constant.DefaultThumbnailQuality)
	}
	return u.String()
}

// ExtractTextFromHTML 从 HTML 中提取纯文本
//
//	@param s HTML 内容
//	@return string 提取的纯文本
//	@author centonhuang
//	@update 2026-01-31 10:00:00
func ExtractTextFromHTML(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return ""
	}
	var text strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)
	return strings.Join(strings.Fields(text.String()), " ")
}
