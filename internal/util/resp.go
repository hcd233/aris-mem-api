// Package util 工具包
package util

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

// WrapSSEResponse 包装SSE响应
//
//	@param ch chan<- *protocol.SSEResponse
//	@param err error
//	@return error
//	@author centonhuang
//	@update 2025-11-11 03:07:17
func WrapSSEResponse(ctx *fiber.Ctx, ch <-chan *protocol.SSEResponse, err error) error {
	logger := logger.WithCtx(ctx.Context())

	if err != nil {
		return WriteErrorResponse(ctx.Response().BodyWriter(), constant.ErrInternalError)
	}
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			resp, ok := <-ch
			if !ok {
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(resp)))
			err = w.Flush()
			if err != nil {
				logger.Warn("[WrapSSEResponse] failed to flush writer", zap.Error(err))
				return
			}
		}
	}))
	return nil
}
