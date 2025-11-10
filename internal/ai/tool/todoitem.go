package tool

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	createTodoItemsToolName        = "createTodoItems"
	createTodoItemsToolDescription = "基于用户输入创建一系列待办事项，返回创建结果"
)

// CreateTodoItemsHandler 创建待办事项处理器
//
//	@author centonhuang
//	@update 2025-11-08 17:37:08
type CreateTodoItemsHandler func(ctx context.Context, req *dto.CreateTodoItemsReq) (output *dto.EmptyRsp, err error)

// TodoItem 待办事项实体
//
//	@author centonhuang
//	@update 2025-11-07 17:31:26
type TodoItem struct {
	Name     string                `json:"name" jsonschema:"description=待办事项名字，10个字以内"`
	Summary  string                `json:"summary" jsonschema:"description=待办事项摘要，100个字以内"`
	Content  string                `json:"content" jsonschema:"description=待办事项具体内容，1000个字以内"`
	Priority enum.TodoItemPriority `json:"priority" jsonschema:"description=待办事项优先级,enum=low,enum=medium,enum=high,enum=urgent"`
}

func (t *TodoItem) toDto() *dto.TodoItem {
	return &dto.TodoItem{
		Name:     t.Name,
		Summary:  t.Summary,
		Content:  t.Content,
		Priority: t.Priority,
	}
}

// CreateTodoItemsInput 创建待办事项输入
//
//	@author centonhuang
//	@update 2025-11-07 17:33:12
type CreateTodoItemsInput struct {
	TodoItems []*TodoItem `json:"todoItems" jsonschema:"description=待办事项列表"`
}

// NewCreateTodoItemsTool 创建创建待办事项工具
//
//	创建创建待办事项工具
//	@return tool.InvokableTool 创建待办事项工具
//	@return error
//	@author centonhuang
//	@update 2025-11-07 17:36:16
func NewCreateTodoItemsTool(handler CreateTodoItemsHandler) (tool.InvokableTool, error) {
	return utils.InferTool(
		createTodoItemsToolName,
		createTodoItemsToolDescription,
		func(ctx context.Context, input *CreateTodoItemsInput) (output *CommonToolOutput, err error) {
			logger := logger.WithCtx(ctx)
			logger.Info("[CreateTodoItemsTool] create todo items", zap.Any("input", input))

			output = &CommonToolOutput{}

			req := &dto.CreateTodoItemsReq{
				Body: &dto.CreateTodoItemsReqBody{
					TodoItems: lo.Map(input.TodoItems, func(item *TodoItem, _ int) *dto.TodoItem {
						return item.toDto()
					}),
				},
			}

			_, err = handler(ctx, req)
			if err != nil {
				logger.Error("[CreateTodoItemsTool] create todo items error", zap.Error(err))
				output.Result = "创建待办事项失败，出现系统内部错误"
				return output, nil
			}
			logger.Info("[CreateTodoItemsTool] create todo items success")
			output.Result = "创建待办事项成功"
			return output, nil
		},
	)
}
