package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/hcd233/aris-mem-api/internal/config"
)

// NewOpenAIChatModel 创建 OpenAI 聊天模型
//
//	@param ctx
//	@return model.ToolCallingChatModel
//	@return error
//	@author centonhuang
//	@update 2025-11-07 17:14:11
func NewOpenAIChatModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   config.OpenAIModel,
		APIKey:  config.OpenAIAPIKey,
		BaseURL: config.OpenAIBaseURL,
	})
}
