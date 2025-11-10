// Package util 工具包
package util

import (
	"github.com/hcd233/aris-mem-api/internal/protocol"
)

// WrapHTTPResponse 包装HTTP响应错误
//
//	@param rsp rspT
//	@param err error
//	@return *protocol.HumaHTTPResponse[rspT]
//	@return error
//	@author centonhuang
//	@update 2025-10-31 01:47:14
func WrapHTTPResponse[rspT any](rsp rspT, err error) (*protocol.HTTPResponse[rspT], error) {
	return &protocol.HTTPResponse[rspT]{
		Body: rsp,
	}, err
}
