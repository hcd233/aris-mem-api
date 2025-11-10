package constant

import "time"

const (
	// HeartbeatInterval SSE心跳间隔
	//
	//	@author centonhuang
	//	@update 2025-11-08 04:43:54
	HeartbeatInterval = time.Second * 1

	// AgentChatLockExpire Agent聊天锁过期时间
	//	@update 2025-11-11 02:46:11
	AgentChatLockExpire = time.Second * 10
)
