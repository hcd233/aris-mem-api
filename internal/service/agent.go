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
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/lock"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/util"
	"github.com/samber/lo"
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

	tools := []etool.BaseTool{createTodoItemsTool, listTodoItemsTool}

	todoAgent, err := agent.NewTodoAgent(ctx, chatModel, tools)
	if err != nil {
		logger.Error("[AgentService] failed to create agent", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}

	toolNames := lo.Map(tools, func(tool etool.BaseTool, _ int) string {
		return lo.Must1(tool.Info(ctx)).Name
	})

	logger.Info("[AgentService] init agent", zap.String("llm", config.OpenAIModel), zap.Strings("tools", toolNames))

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

		// CheckPointStore: checkpoint.NewRedisCheckPointStore(),
	})

	iter := runner.Query(ctx, req.Body.Message) // adk.WithCheckPointID(strconv.FormatUint(uint64(userID), 10))

	return util.WrapADKIterSSE(ctx, iter), nil
}
