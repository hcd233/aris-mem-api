package constant

import "time"

const (
	// HeartbeatInterval SSE心跳间隔
	//
	//	@author centonhuang
	//	@update 2025-11-08 04:43:54
	HeartbeatInterval = 1 * time.Second

	// AgentChatLockExpire Agent聊天锁过期时间
	//	@update 2025-11-11 02:46:11
	AgentChatLockExpire = 5 * time.Minute

	// PresignObjectExpire 预签名对象过期时间
	//	@update 2025-11-12 19:20:26
	PresignObjectExpire = 5 * time.Minute
)
