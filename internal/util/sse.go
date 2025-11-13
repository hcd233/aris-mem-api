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
	"github.com/hcd233/aris-mem-api/internal/common/model"
	"github.com/hcd233/aris-mem-api/internal/lock"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// WrapADKIterSSE 将Adk迭代器转换为SSE响应
//
//	@param ctx context.Context
//	@param iter *adk.AsyncIterator[*adk.AgentEvent]
//	@return rsp
//	@author centonhuang
//	@update 2025-11-11 17:43:42
func WrapADKIterSSE(ctx context.Context, iter *adk.AsyncIterator[*adk.AgentEvent]) (rsp *huma.StreamResponse) {
	logger := logger.WithCtx(ctx)
	locker := lock.NewLocker()

	return &huma.StreamResponse{
		Body: func(hCtx huma.Context) {
			fCtx := humafiber.Unwrap(hCtx)
			fCtx.Set("Content-Type", "text/event-stream")
			fCtx.Set("Cache-Control", "no-cache")
			fCtx.Set("Connection", "keep-alive")
			fCtx.Set("Transfer-Encoding", "chunked")

			fCtx.Response().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				lockKey := fmt.Sprintf(constant.LockKeyTemplateAgentChat, ctx.Value(constant.CtxKeyUserID).(uint))
				lockValue := ctx.Value(constant.CtxKeyTraceID).(string)
				success, err := locker.Lock(ctx, lockKey, lockValue, constant.AgentChatLockExpire)
				if err != nil {
					logger.Error("[AgentService] lock resource error", zap.Error(err))
					writeSSEErrorResponse(ctx, w, constant.ErrInternalError)
					return
				}
				if !success {
					logger.Info("[AgentService] lock resource is already locked", zap.String("lockKey", lockKey))
					writeSSEErrorResponse(ctx, w, constant.ErrResourceLocked)
					return
				}
				defer func() {
					lockKey := fmt.Sprintf(constant.LockKeyTemplateAgentChat, ctx.Value(constant.CtxKeyUserID))
					lockValue := ctx.Value(constant.CtxKeyTraceID).(string)
					err := locker.Unlock(ctx, lockKey, lockValue)
					if err != nil {
						logger.Error("[StreamADKIter] unlock resource error", zap.Error(err))
					}
				}()

				writeSSENoneResponse(ctx, w, enum.SSEStatusStart)
				defer writeSSENoneResponse(ctx, w, enum.SSEStatusEnd)

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
						writeSSEErrorResponse(ctx, w, constant.ErrInternalError)
						return
					}

					messageStream := event.Output.MessageOutput.MessageStream
					if messageStream == nil {
						message, err := event.Output.MessageOutput.GetMessage()
						if err != nil {
							logger.Error("[AgentService] failed to get message", zap.Error(err))
							writeSSEErrorResponse(ctx, w, constant.ErrInternalError)
						}
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
							writeSSEErrorResponse(ctx, w, constant.ErrInternalError)
							return
						}
						writeSSEMessageResponse(ctx, w, message)
					}

				}
			}))
		},
	}
}

// WrapErrorSSE 包装错误响应
//
//	@param ctx
//	@param err
//	@return rsp
//	@author centonhuang
//	@update 2025-11-11 17:46:36
func WrapErrorSSE(ctx context.Context, err *model.Error) (rsp *huma.StreamResponse) {
	return &huma.StreamResponse{
		Body: func(hCtx huma.Context) {
			fCtx := humafiber.Unwrap(hCtx)
			fCtx.Set("Content-Type", "text/event-stream")
			fCtx.Set("Cache-Control", "no-cache")
			fCtx.Set("Connection", "keep-alive")
			fCtx.Set("Transfer-Encoding", "chunked")

			fCtx.Response().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				writeSSEErrorResponse(ctx, w, err)
			}))
		},
	}
}

func writeSSEMessageResponse(ctx context.Context, w *bufio.Writer, message adk.Message) {
	logger := logger.WithCtx(ctx)
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeMessage,
		Status:   enum.SSEStatusStreaming,
		Data:     message,
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteMessageResponse] flush error", zap.Error(err))
	}
}

func writeSSEErrorResponse(ctx context.Context, w *bufio.Writer, err *model.Error) {
	logger := logger.WithCtx(ctx)
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeError,
		Status:   enum.SSEStatusError,
		Data:     &dto.CommonRsp{Error: err},
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
		Status:   enum.SSEStatusStreaming,
		Data:     strconv.Itoa(heartBeatCount),
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteHeartBeatResponse] flush error", zap.Error(err))
	}
}

func writeSSENoneResponse(ctx context.Context, w *bufio.Writer, status enum.SSEStatus) {
	logger := logger.WithCtx(ctx)
	rsp := &protocol.SSEResponse{
		DataType: enum.SSEDataTypeNone,
		Status:   status,
		Data:     nil,
	}
	fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
	if err := w.Flush(); err != nil {
		logger.Error("[WriteNoneResponse] flush error", zap.Error(err))
	}
}
