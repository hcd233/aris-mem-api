package router

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
	"github.com/hcd233/aris-mem-api/internal/handler"
	"github.com/hcd233/aris-mem-api/internal/protocol"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
)

// initHealthRouter 初始化健康检查路由
//
//	@param healthGroup
//	@author centonhuang
//	@update 2025-11-07 14:59:06
func initHealthRouter(healthGroup huma.API) {
	pingHandler := handler.NewPingHandler()

	huma.Register(healthGroup, huma.Operation{
		OperationID: "healthCheck",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "HealthCheck",
		Description: "Check the server health",
		Tags:        []string{"health"},
	}, pingHandler.HandlePing)
}

func initSSEHealthRouter(sseGroup fiber.Router) {
	sseGroup.Get("ssehealth", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Transfer-Encoding", "chunked")

		ctx.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for i := 0; i < 3; i++ {
				rsp := &protocol.SSEResponse{
					DataType: enum.SSEDataTypeHeartBeat,
					Data:     strconv.Itoa(i),
				}
				fmt.Fprintf(w, "data: %s\n\n", lo.Must1(sonic.Marshal(rsp)))
				err := w.Flush()
				if err != nil {
					return
				}
				time.Sleep(constant.HeartbeatInterval)
			}
		}))
		return nil
	})
}
