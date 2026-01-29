package util

import "regexp"

var tagRegex = regexp.MustCompile(`#([\p{L}\p{N}_]+)`)

// ExtractTags 提取内容中的标签
//
//	@param content
//	@return []string
//	@author centonhuang
//	@update 2026-01-28 21:51:07
func ExtractTags(content string) []string {
	matches := tagRegex.FindAllStringSubmatch(content, -1)
	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}
	return tags
}
