// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	etool "github.com/cloudwego/eino/components/tool"
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/ai/agent"
	"github.com/hcd233/aris-mem-api/internal/ai/llm"
	"github.com/hcd233/aris-mem-api/internal/ai/tool"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/lock"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// AgentService Agent服务
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type AgentService interface {
	HandleChat(ctx context.Context, req *dto.ChatReq) (rsp *huma.StreamResponse, err error)
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
func (s *agentService) HandleChat(ctx context.Context, req *dto.ChatReq) (rsp *huma.StreamResponse, err error) {
	logger := logger.WithCtx(ctx)

	locker := lock.NewLocker()

	chatModel, err := llm.NewOpenAIChatModel(ctx)
	if err != nil {
		logger.Error("[AgentService] failed to create model", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}

	createTodoItemsTool, err := tool.NewCreateTodoItemsTool(s.todoItemService.CreateTodoItems)
	if err != nil {
		logger.Error("[AgentService] failed to create create todo items tool", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}
	listTodoItemsTool, err := tool.NewListTodoItemsTool(s.todoItemService.ListTodoItems)
	if err != nil {
		logger.Error("[AgentService] failed to create list todo items tool", zap.Error(err))
	}

	todoAgent, err := agent.NewTodoAgent(ctx, chatModel, []etool.BaseTool{createTodoItemsTool, listTodoItemsTool})
	if err != nil {
		logger.Error("[AgentService] failed to create agent", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}

	lockKey := fmt.Sprintf(constant.LockKeyTemplateAgentChat, ctx.Value(constant.CtxKeyUserID).(uint))
	lockValue := ctx.Value(constant.CtxKeyTraceID).(string)
	success, err := locker.Lock(ctx, lockKey, lockValue, constant.AgentChatLockExpire)
	if err != nil {
		logger.Error("[AgentService] lock resource error", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}
	if !success {
		logger.Info("[AgentService] lock resource is already locked", zap.String("lockKey", lockKey))
		return util.WrapErrorSSE(ctx, constant.ErrTooManyRequests), nil
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           todoAgent,
		EnableStreaming: true,
	})

	iter := runner.Query(ctx, req.Body.Message)

	return util.WrapADKIterSSE(ctx, iter), nil
}
