// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"
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
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/hcd233/aris-mem-api/internal/resource/database"
	"github.com/hcd233/aris-mem-api/internal/resource/database/dao"
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
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
	dialogDAO       *dao.DialogDAO
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
		dialogDAO:       dao.GetDialogDAO(),
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
	db := database.GetDBInstance(ctx)

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

	commonParam := &dao.CommonParam{
		PageParam: dao.PageParam{
			Page:     1,
			PageSize: constant.DialogHistoryTurn,
		},
		SortParam: dao.SortParam{
			Sort:      enum.SortDesc,
			SortField: "id",
		},
	}

	historyDialog, _, err := s.dialogDAO.Paginate(db, &model.Dialog{UserID: userID, Status: enum.DialogStatusCompleted}, []string{"id", "input_messages", "output_messages"}, commonParam)
	if err != nil {
		logger.Error("[AgentService] failed to get history dialog", zap.Error(err))
		return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
	}

	messages := lo.Flatten(lo.Map(lo.Reverse(historyDialog), func(dialog *model.Dialog, _ int) []*schema.Message {
		return append(dialog.InputMessages, dialog.OutputMessages...)
	}))

	inputContent := []schema.MessageInputPart{}
	if bodyData.Content != "" {
		inputContent = append(inputContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: bodyData.Content,
		})
	}

	if audio := bodyData.Audio; audio.Size > 0 {
		if !strings.HasSuffix(audio.Filename, ".wav") {
			logger.Warn("[AgentService] audio file is not a wav file", zap.String("filename", audio.Filename))
			return util.WrapErrorSSE(ctx, constant.ErrBadRequest), nil
		}
		file := audio.File
		defer file.Close()

		err = s.audioObjDAO.UploadObject(ctx, userID, audio.Filename, audio.Size, file)
		if err != nil {
			logger.Error("[AgentService] failed to upload audio", zap.String("filename", audio.Filename), zap.Error(err))
			return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		}

		presignedURL, err := s.audioObjDAO.PresignObject(ctx, userID, audio.Filename)
		if err != nil {
			logger.Error("[AgentService] failed to presign audio", zap.String("filename", audio.Filename), zap.Error(err))
			return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		}

		_ = lo.Must1(file.Seek(0, io.SeekStart))

		// audioBytes, err := io.ReadAll(file)
		// if err != nil {
		// 	logger.Error("[AgentService] failed to read audio", zap.String("filename", bodyData.Audio.Filename), zap.Error(err))
		// 	return util.WrapErrorSSE(ctx, constant.ErrInternalError), nil
		// }

		// audioBase64 := base64.StdEncoding.EncodeToString(audioBytes)

		inputContent = append(inputContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeAudioURL,
			Audio: &schema.MessageInputAudio{
				MessagePartCommon: schema.MessagePartCommon{
					Base64Data: lo.ToPtr(presignedURL.String()), // NOTE: EINO OpenAI Lib only supports base64 data, not URL. It is a trick
					// URL: lo.ToPtr(presignedURL.String()),
					// Base64Data: lo.ToPtr(fmt.Sprintf("data:%s;base64,%s", bodyData.Audio.ContentType, audioBase64)),
					MIMEType: audio.ContentType,
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

	userMessage := &schema.Message{
		Role:                  schema.User,
		UserInputMultiContent: inputContent,
	}

	messages = append(messages, userMessage)

	iter := runner.Run(ctx, messages) // adk.WithCheckPointID(strconv.FormatUint(uint64(userID), 10))

	return util.WrapADKIterSSE(ctx, iter, userMessage), nil
}
