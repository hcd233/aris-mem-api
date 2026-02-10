package middleware

import (
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/fiber/v2"
)

// FgprofMiddleware fgprof中间件
//
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-09-25 21:17:02
func FgprofMiddleware() fiber.Handler {
	// go tool pprof -http=:8081 http://xxxx/debug/fgprof
	return fgprof.New()
}
