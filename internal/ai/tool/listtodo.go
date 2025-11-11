package tool

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/common/model"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/protocol/dto"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	listTodoItemsToolName        = "listTodoItems"
	listTodoItemsToolDescription = "查询待办事项列表"
)

// ListTodoItemsHandler 查询待办事项列表处理器
//
//	@param ctx
//	@param req
//	@return output
//	@return err
//	@author centonhuang
//	@update 2025-11-11 15:04:44
type ListTodoItemsHandler func(ctx context.Context, req *dto.ListTodoItemsReq) (output *dto.ListTodoItemsRsp, err error)

// ListTodoItemsInput TodoItem 待办事项实体
//
//	@author centonhuang
//	@update 2025-11-11 15:27:24
type ListTodoItemsInput struct {
	Page      int       `json:"page" jsonschema:"description=分页，默认从1开始"`
	PageSize  int       `json:"pageSize" jsonschema:"description=每页大小，用户不指定默认10条，最大50条"`
	Query     string    `json:"query" jsonschema:"description=查询关键词，用户不指定默认为空"`
	Sort      enum.Sort `json:"sort" jsonschema:"description=排序方式，用户不指定默认升序"`
	SortField string    `json:"sortField" jsonschema:"description=排序字段，用户不指定默认id,enum=id,enum=createdAt,enum=updatedAt"`
}

func (l *ListTodoItemsInput) toDto() *dto.ListTodoItemsReq {
	return &dto.ListTodoItemsReq{
		CommonParam: model.CommonParam{
			PageParam: model.PageParam{
				Page:     l.Page,
				PageSize: l.PageSize,
			},
			QueryParam: model.QueryParam{
				Query: l.Query,
			},
			SortParam: model.SortParam{
				Sort: l.Sort,
			},
		},
		SortField: l.SortField,
	}
}

// NewListTodoItemsTool 创建查询待办事项列表工具
//
//	@param handler
//	@return tool.InvokableTool
//	@return error
//	@author centonhuang
//	@update 2025-11-11 15:33:17
func NewListTodoItemsTool(handler ListTodoItemsHandler) (tool.InvokableTool, error) {
	return utils.InferTool(
		listTodoItemsToolName,
		listTodoItemsToolDescription,
		func(ctx context.Context, input *ListTodoItemsInput) (output *CommonToolOutput, err error) {
			logger := logger.WithCtx(ctx)
			logger.Info("[ListTodoItemsTool] list todo items", zap.Any("input", input))

			output = &CommonToolOutput{}
			req := input.toDto()

			rsp, err := handler(ctx, req)
			if err != nil {
				logger.Error("[ListTodoItemsTool] list todo items error", zap.Error(err))
				output.Result = "查询待办事项列表失败，出现系统内部错误"
				return output, nil
			}

			output.Result = fmt.Sprintf(`查询待办事项列表成功
## 数据:
%s
## 分页信息:
第%d页，每页%d条，共%d条
`, lo.Must1(sonic.MarshalIndent(rsp.TodoItems, "", "  ")), rsp.PageInfo.Page, rsp.PageInfo.PageSize, rsp.PageInfo.Total)
			return output, nil
		},
	)
}
