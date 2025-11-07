// Package util 工具包
package util

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
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
func WrapHTTPResponse[rspT any](rsp rspT, err error) (*protocol.HTTPResponse[rspT], huma.StatusError) {
	if statusErr := transformError(err); statusErr != nil {
		return nil, statusErr
	}
	return &protocol.HTTPResponse[rspT]{
		Body: rsp,
	}, nil
}

func transformError(err error) (statusErr huma.StatusError) {
	switch err {
	case constant.ErrDataNotExists: // 404
		statusErr = huma.Error404NotFound(err.Error())
	case constant.ErrDataExists, constant.ErrBadRequest, constant.ErrInsufficientQuota: // 400
		statusErr = huma.Error400BadRequest(err.Error())
	case constant.ErrUnauthorized: // 401
		statusErr = huma.Error401Unauthorized(err.Error())
	case constant.ErrNoPermission: // 403
		statusErr = huma.Error403Forbidden(err.Error())
	case constant.ErrTooManyRequests: // 429
		statusErr = huma.Error429TooManyRequests(err.Error())
	case constant.ErrInternalError: // 500
		statusErr = huma.Error500InternalServerError(err.Error())
	case constant.ErrNoImplement: // 501
		statusErr = huma.Error501NotImplemented(err.Error())
	case nil:
		statusErr = nil
	default:
		statusErr = huma.Error500InternalServerError("Unknown error: " + err.Error())
	}
	return
}
