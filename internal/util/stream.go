package util

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// AdkIterToChan  将迭代器转换为通道
//
//	@param ctx context.Context
//	@param iter *adk.AsyncIterator[*adk.AgentEvent]
//	@return chan *protocol.SSEResponse
//	@return error
//	@author centonhuang
//	@update 2025-11-11 03:05:06
func AdkIterToChan(ctx context.Context, iter *adk.AsyncIterator[*adk.AgentEvent]) (rsp *huma.StreamResponse) {
	logger := logger.WithCtx(ctx)

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			fCtx := humafiber.Unwrap(ctx)
			fCtx.Set("Content-Type", "text/event-stream")
			fCtx.Set("Cache-Control", "no-cache")
			fCtx.Set("Connection", "keep-alive")
			fCtx.Set("Transfer-Encoding", "chunked")

			fCtx.Response().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				ticker := time.NewTicker(constant.HeartbeatInterval)
				defer ticker.Stop()
				go func() {
					heartBeatCount := 0
					for {
						select {
						case <-ctx.Context().Done():
							return
						case <-ticker.C:
							rsp := &protocol.SSEResponse{
								DataType: enum.SSEDataTypeHeartBeat,
								Data:     strconv.Itoa(heartBeatCount),
							}
							fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
							if err := w.Flush(); err != nil {
								logger.Error("[AdkIterToChan] flush error", zap.Error(err))
								return
							}
							heartBeatCount++
						}
					}
				}()

				for {
					event, ok := iter.Next()
					if !ok {
						logger.Info("[AgentService] reach iter end")
						time.Sleep(constant.HeartbeatInterval)
						return
					}
					if event.Err != nil {
						logger.Error("[AgentService] agent run error", zap.Error(event.Err))
						rsp := &protocol.SSEResponse{
							DataType: enum.SSEDataTypeError,
							Data:     event.Err.Error(),
						}
						fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
						if err := w.Flush(); err != nil {
							logger.Error("[AdkIterToChan] flush error", zap.Error(err))
							return
						}
						return
					}
					message, err := event.Output.MessageOutput.GetMessage()
					if err != nil {
						logger.Error("[AgentService] failed to get message", zap.Error(err))
						rsp := &protocol.SSEResponse{
							DataType: enum.SSEDataTypeError,
							Data:     err.Error(),
						}
						fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
						if err := w.Flush(); err != nil {
							logger.Error("[AdkIterToChan] flush error", zap.Error(err))
							return
						}
						return
					}
					messageData := lo.Must1(sonic.Marshal(message))
					logger.Info("[AgentService] receive event message", zap.ByteString("message", messageData))
					rsp := &protocol.SSEResponse{
						DataType: enum.SSEDataTypeMessage,
						Data:     string(messageData),
					}
					fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
					if err := w.Flush(); err != nil {
						logger.Error("[AdkIterToChan] flush error", zap.Error(err))
						return
					}
				}
			}))
		},
	}
}
