// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	etool "github.com/cloudwego/eino/components/tool"
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/ai/agent"
	"github.com/hcd233/aris-mem-api/internal/ai/llm"
	"github.com/hcd233/aris-mem-api/internal/ai/tool"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// AgentService Agent服务
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type AgentService interface {
	HandleChat(ctx context.Context, req *dto.ChatReq) (*huma.StreamResponse, error)
}

type agentService struct {
	todoItemService TodoItemService
}

// NewAgentService 创建Agent服务
//
//	return AgentService
//	author centonhuang
//	update 2025-01-05 21:00:00
func NewAgentService() AgentService {
	return &agentService{
		todoItemService: NewTodoItemService(),
	}
}

// HandleChat 处理聊天
//
//	receiver s *agentService
//	param ctx context.Context
//	param req *dto.ChatReq
//	return *huma.StreamResponse, error
func (s *agentService) HandleChat(ctx context.Context, req *dto.ChatReq) (*huma.StreamResponse, error) {
	logger := logger.WithCtx(ctx)

	chatModel, err := llm.NewOpenAIChatModel(ctx)
	if err != nil {
		logger.Error("[AgentService] failed to create model", zap.Error(err))
		return nil, err
	}

	createTodoItemsTool, err := tool.NewCreateTodoItemsTool(s.todoItemService.CreateTodoItems)
	todoAgent, err := agent.NewTodoAgent(ctx, chatModel, []etool.BaseTool{createTodoItemsTool})
	if err != nil {
		logger.Error("[AgentService] failed to create agent", zap.Error(err))
		return nil, err
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           todoAgent,
		EnableStreaming: true,
	})

	iter := runner.Query(ctx, req.Body.Message)

	return &huma.StreamResponse{
		Body: func(hctx huma.Context) {
			hctx.SetHeader("Content-Type", "text/event-stream")
			hctx.SetHeader("Cache-Control", "no-cache")
			hctx.SetHeader("Connection", "keep-alive")

			ctx := hctx.BodyWriter().(*fasthttp.RequestCtx)
			ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				ticker := time.NewTicker(constant.HeartbeatInterval)
				defer ticker.Stop()
				go func() {
					heartBeatCount := 0
					for {
						select {
						case <-ctx.Done():
							return
						case <-ticker.C:
							data := lo.Must1(sonic.Marshal(protocol.SSEResponse{
								DataType: enum.SSEDataTypeHeartBeat,
								Data:     strconv.Itoa(heartBeatCount),
							}))
							heartBeatCount++
							fmt.Fprintf(w, "id: %d\n", heartBeatCount)
							fmt.Fprintf(w, "event: heartbeat\n")
							fmt.Fprintf(w, "data: %s\n\n", data)
							w.Flush()
						}
					}
				}()

				for {
					event, ok := iter.Next()
					if !ok {
						logger.Info("[AgentService] reach iter end")
						return
					}
					if event.Err != nil {
						logger.Error("[AgentService] agent run error", zap.Error(event.Err))
						data := lo.Must1(sonic.Marshal(protocol.SSEResponse{
							DataType: enum.SSEDataTypeError,
							Data:     event.Err.Error(),
						}))
						fmt.Fprintf(w, "event: error\n")
						fmt.Fprintf(w, "data: %s\n\n", data)
						w.Flush()
						return
					}
					message, err := event.Output.MessageOutput.GetMessage()
					if err != nil {
						logger.Error("[AgentService] failed to get message", zap.Error(err))
						data := lo.Must1(sonic.Marshal(protocol.SSEResponse{
							DataType: enum.SSEDataTypeError,
							Data:     err.Error(),
						}))
						fmt.Fprintf(w, "event: error\n")
						fmt.Fprintf(w, "data: %s\n\n", data)
						w.Flush()
						return
					}
					messageData := lo.Must1(sonic.Marshal(message))
					logger.Info("[AgentService] receive event message", zap.ByteString("message", messageData))
					data := lo.Must1(sonic.Marshal(protocol.SSEResponse{
						DataType: enum.SSEDataTypeMessage,
						Data:     string(messageData),
					}))
					fmt.Fprintf(w, "event: message\n")
					fmt.Fprintf(w, "data: %s\n\n", data)
					w.Flush()
				}
			}))
		},
	}, nil
}
