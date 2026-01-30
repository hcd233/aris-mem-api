package constant

const (

	// LockKeyTemplateAgentChat Agent聊天锁键模板
	//	@update 2025-11-11 17:10:45
	LockKeyTemplateAgentChat = "agent:chat:%v"

	// LockKeyTemplateMiddleware 中间件锁键模板
	//	@update 2025-11-11 17:23:31
	LockKeyTemplateMiddleware = "%s:%s:%v"

	// CacheKeyTemplateCosTempSecret COS临时密钥缓存键模板
	//	@update 2026-01-31 05:58:39
	CacheKeyTemplateCosTempSecret = "cos_temp_secret:%d"
)
