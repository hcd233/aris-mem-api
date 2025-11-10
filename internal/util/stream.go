package util

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/samber/lo"
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
func AdkIterToChan(ctx context.Context, iter *adk.AsyncIterator[*adk.AgentEvent]) (ch chan *protocol.SSEResponse) {
	syncCh := make(chan struct{})
	ch = make(chan *protocol.SSEResponse)
	logger := logger.WithCtx(ctx)

	ticker := time.NewTicker(constant.HeartbeatInterval)
	go func() {
		defer ticker.Stop()
		heartBeatCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-syncCh:
				return
			case <-ticker.C:
				ch <- &protocol.SSEResponse{
					DataType: enum.SSEDataTypeHeartBeat,
					Data:     strconv.Itoa(heartBeatCount),
				}
				heartBeatCount++
			}
		}
	}()

	go func() {
		defer func() {
			syncCh <- struct{}{}
			close(ch)
			close(syncCh)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				event, ok := iter.Next()
				if !ok {
					logger.Info("[AgentService] reach iter end")
					return
				}
				if event.Err != nil {
					logger.Error("[AgentService] agent run error", zap.Error(event.Err))
					ch <- &protocol.SSEResponse{
						DataType: enum.SSEDataTypeError,
						Data:     event.Err.Error(),
					}
					return
				}
				message, err := event.Output.MessageOutput.GetMessage()
				if err != nil {
					logger.Error("[AgentService] failed to get message", zap.Error(err))
					ch <- &protocol.SSEResponse{
						DataType: enum.SSEDataTypeError,
						Data:     err.Error(),
					}
					return
				}
				messageData := lo.Must1(sonic.Marshal(message))
				logger.Info("[AgentService] receive event message", zap.ByteString("message", messageData))
				ch <- &protocol.SSEResponse{
					DataType: enum.SSEDataTypeMessage,
					Data:     string(messageData),
				}
			}
		}
	}()
	return ch
}
