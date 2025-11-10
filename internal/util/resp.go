// Package util 工具包
package util

import (
	"io"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
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

// WriteErrorResponse 写入错误响应
//
//	@param ctx
//	@param err
//	@return error
//	@author centonhuang
//	@update 2025-11-10 20:55:14
func WriteErrorResponse(bodyWriter io.Writer, err *model.Error) error {
	_, writeErr := bodyWriter.Write(lo.Must1(sonic.Marshal(&dto.CommonRsp{Error: err})))
	return writeErr
}
