package util

import "regexp"

var tagRegex = regexp.MustCompile(`#(\w+)`)

// ExtractTags 提取内容中的标签
//	@param content 
//	@return []string 
//	@author centonhuang 
//	@update 2026-01-28 21:51:07 
func ExtractTags(content string) []string {
	return tagRegex.FindAllString(content, -1)
}