// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"

	"github.com/cloudwego/eino/adk"
	etool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/hcd233/aris-mem-api/internal/ai/agent"
	"github.com/hcd233/aris-mem-api/internal/ai/llm"
	"github.com/hcd233/aris-mem-api/internal/ai/tool"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"go.uber.org/zap"
)

// AgentService Agent服务
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type AgentService interface {
	HandleChat(ctx context.Context, req *dto.ChatReq, sender sse.Sender)
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
//	param sender sse.Sender
func (s *agentService) HandleChat(ctx context.Context, req *dto.ChatReq, sender sse.Sender) {
	logger := logger.WithCtx(ctx)

	chatModel, err := llm.NewOpenAIChatModel(ctx)
	if err != nil {
		logger.Error("[AgentService] failed to create model", zap.Error(err))
		return
	}

	createTodoItemsTool, err := tool.NewCreateTodoItemsTool(s.todoItemService.CreateTodoItems)
	todoAgent, err := agent.NewTodoAgent(ctx, chatModel, []etool.BaseTool{createTodoItemsTool})
	if err != nil {
		logger.Error("[AgentService] failed to create agent", zap.Error(err))
		return
	}

	iter := todoAgent.Run(ctx, &adk.AgentInput{
		Messages: []*schema.Message{
			schema.UserMessage(req.Body.Message),
		},
		EnableStreaming: true,
	})
	for {
		event, ok := iter.Next()
		if !ok {
			logger.Info("[AgentService] reach iter end")
			sender.Data(protocol.SSEResponse{
				DataType: enum.SSEDataTypeText,
				Data:     "",
			})
			break
		}
		sr := event.Output.MessageOutput.MessageStream
		for {
			message, err := sr.Recv()
			if err != nil {
				logger.Error("[AgentService] failed to recv message", zap.Error(err))
				sender.Data(protocol.SSEResponse{
					DataType: enum.SSEDataTypeError,
					Data:     err.Error(),
				})
				break
			}
			sender.Data(protocol.SSEResponse{
				DataType: enum.SSEDataTypeText,
				Data:     message.Content,
			})
		}
	}
}
