package enum

// DialogStatus 对话状态
//
//	@author centonhuang
//	@update 2025-11-13 19:10:00
type DialogStatus string

const (
	// DialogStatusCompleted DialogStatus 已完成状态
	//	@author centonhuang
	//	@update 2025-11-13 19:10:00
	DialogStatusCompleted DialogStatus = "completed"

	// DialogStatusCancelled DialogStatus 已取消状态
	//	@author centonhuang
	//	@update 2025-11-13 19:10:00
	DialogStatusCancelled DialogStatus = "cancelled"

	// DialogStatusError DialogStatus 错误状态
	//	@update 2025-11-13 19:39:28
	DialogStatusError DialogStatus = "error"
)
