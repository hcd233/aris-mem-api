package agent

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

const (
	todoAgentName        = "todoAgent"
	todoAgentDescription = "一个管理待办事项的Agent，擅长理解并且拆解用户描述的事项，转换为待办事项"
)

// NewTodoAgent 创建待办事项Agent
//
//	@param ctx context.Context
//	@param model model.ToolCallingChatModel
//	@param tools []tool.BaseTool
//	@return adk.Agent
//	@return error
//	@author centonhuang
//	@update 2025-11-08 05:13:07
func NewTodoAgent(ctx context.Context, chatModel model.ToolCallingChatModel, tools []tool.BaseTool) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        todoAgentName,
		Description: todoAgentDescription,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
