package util

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
		Body: func(hCtx huma.Context) {
			fCtx := humafiber.Unwrap(hCtx)
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
						case <-ctx.Done():
							return
						case <-ticker.C:
							writeSSEHeartBeatResponse(ctx, w, heartBeatCount)
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
						writeSSEErrorResponse(ctx, w, event.Err)
						return
					}

					messageStream := event.Output.MessageOutput.MessageStream
					if messageStream == nil {
						message, err := event.Output.MessageOutput.GetMessage()
						if err != nil {
							logger.Error("[AgentService] failed to get message", zap.Error(err))
							writeSSEErrorResponse(ctx, w, err)
						}
						messageData := lo.Must1(sonic.Marshal(message))
						logger.Info("[AgentService] receive event message", zap.ByteString("message", messageData))
						writeSSEMessageResponse(ctx, w, message)
						continue
					}
					for {
						message, err := messageStream.Recv()
						if err != nil {
							if err == io.EOF {
								logger.Info("[AgentService] reach message stream end")
								break
							}

							logger.Error("[AgentService] failed to get message", zap.Error(err))
							writeSSEErrorResponse(ctx, w, err)
							return
						}
						writeSSEMessageResponse(ctx, w, message)
					}

				}
			}))
		},
	}
}

func writeSSEMessageResponse(ctx context.Context, w *bufio.Writer, message adk.Message) {
	logger := logger.WithCtx(ctx)
	messageData := lo.Must1(sonic.Marshal(message))
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeMessage,
		Data:     string(messageData),
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteMessageResponse] flush error", zap.Error(err))
	}
}

func writeSSEErrorResponse(ctx context.Context, w *bufio.Writer, err error) {
	logger := logger.WithCtx(ctx)
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeError,
		Data:     err.Error(),
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteErrorResponse] flush error", zap.Error(err))
	}
}

func writeSSEHeartBeatResponse(ctx context.Context, w *bufio.Writer, heartBeatCount int) {
	logger := logger.WithCtx(ctx)
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeHeartBeat,
		Data:     strconv.Itoa(heartBeatCount),
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteHeartBeatResponse] flush error", zap.Error(err))
	}
}
