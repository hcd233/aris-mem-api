// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/adk"
	etool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-mem-api/internal/ai/agent"
	"github.com/hcd233/aris-mem-api/internal/ai/llm"
	"github.com/hcd233/aris-mem-api/internal/ai/tool"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	objdao "github.com/hcd233/aris-mem-api/internal/resource/storage/obj_dao"
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
	audioObjDAO     objdao.ObjDAO
}

// NewAgentService 创建Agent服务
//
//	return AgentService
//	author centonhuang
//	update 2025-01-05 21:00:00
func NewAgentService() AgentService {
	return &agentService{
		todoItemService: NewTodoItemService(),
		audioObjDAO:     objdao.GetAudioObjDAO(),
	}
}

// HandleChat 处理聊天
//
//	receiver s *agentService
//	param ctx context.Context
//	param req *dto.ChatReq
//	return *huma.StreamResponse, error
func (s *agentService) HandleChat(ctx context.Context, req *dto.ChatReq) (rsp *huma.StreamResponse, err error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)

	bodyData := req.RawBody.Data()

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

	inputContent := []schema.MessageInputPart{}
	if bodyData.Content != "" {
		inputContent = append(inputContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: bodyData.Content,
		})
	}

	if bodyData.Audio.Size > 0 {
		if !strings.HasSuffix(bodyData.Audio.Filename, ".wav") {
			logger.Warn("[AgentService] audio file is not a wav file", zap.String("filename", bodyData.Audio.Filename))
			return util.WrapErrorSSE(ctx, constant.ErrBadRequest), nil
		}
		file := bodyData.Audio.File
		defer file.Close()

		err = s.audioObjDAO.UploadObject(ctx, userID, bodyData.Audio.Filename, bodyData.Audio.Size, file)
		if err != nil {
			logger.Error("[AgentService] failed to upload audio", zap.String("filename", bodyData.Audio.Filename), zap.Error(err))
			return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		}

		// presignedURL, err := s.audioObjDAO.PresignObject(ctx, userID, bodyData.Audio.Filename)
		// if err != nil {
		// 	logger.Error("[AgentService] failed to presign audio", zap.String("filename", bodyData.Audio.Filename), zap.Error(err))
		// 	return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		// }

		_ = lo.Must1(file.Seek(0, io.SeekStart))

		audioBytes, err := io.ReadAll(file)
		if err != nil {
			logger.Error("[AgentService] failed to read audio", zap.String("filename", bodyData.Audio.Filename), zap.Error(err))
			return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		}

		audioBase64 := base64.StdEncoding.EncodeToString(audioBytes)

		inputContent = append(inputContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeAudioURL,
			Audio: &schema.MessageInputAudio{
				MessagePartCommon: schema.MessagePartCommon{
					// URL: lo.ToPtr(presignedURL.String()),
					Base64Data: lo.ToPtr(fmt.Sprintf("data:%s;base64,%s", bodyData.Audio.ContentType, audioBase64)),
					MIMEType:   bodyData.Audio.ContentType,
				},
			},
		})
	}

	if len(bodyData.Images) > 0 {
		logger.Warn("[AgentService] images are not supported yet", zap.Strings("filenames", lo.Map(bodyData.Images, func(image huma.FormFile, _ int) string {
			return image.Filename
		})))
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           todoAgent,
		EnableStreaming: true,

		// CheckPointStore: checkpoint.NewRedisCheckPointStore(),
	})

	iter := runner.Run(ctx, []adk.Message{
		{
			Role:                  schema.User,
			UserInputMultiContent: inputContent,
		},
	},
	) // adk.WithCheckPointID(strconv.FormatUint(uint64(userID), 10))

	return util.WrapADKIterSSE(ctx, iter), nil
}
