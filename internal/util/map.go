package util

import "reflect"

// HasNonZeroValue 检查是否无空值
//
//	@param fields map[string]interface{}
//	@return bool
//	@author centonhuang
//	@update 2026-01-29 01:57:16
func HasNonZeroValue(fields map[string]interface{}) bool {
	for _, value := range fields {
		if !reflect.ValueOf(value).IsZero() {
			return true
		}
	}
	return false
}
