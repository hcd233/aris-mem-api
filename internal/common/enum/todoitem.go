package enum

type (
	// TodoItemStatus 待办事项状态
	//
	//	@author centonhuang
	//	@update 2025-11-07 01:09:42
	TodoItemStatus string

	// TodoItemPriority 待办事项优先级
	//	@update 2025-11-07 01:11:32
	TodoItemPriority string
)

// const
//
//	@param TodoItemStatusPending TodoItemStatus = "pending"
//	@author centonhuang
//	@update 2025-11-07 01:09:46
const (

	// TodoItemStatusPending TodoItemStatus 待定状态
	//	@update 2025-11-07 01:10:14
	TodoItemStatusPending TodoItemStatus = "pending"

	// TodoItemStatusCompleted TodoItemStatus 已完成状态
	//	@update 2025-11-07 01:10:21
	TodoItemStatusCompleted TodoItemStatus = "completed"

	// TodoItemStatusCancelled TodoItemStatus 已取消状态
	//	@update 2025-11-07 01:10:35
	TodoItemStatusCancelled TodoItemStatus = "cancelled"
)

const (
	// TodoItemPriorityLow TodoItemPriority 低优先级
	//	@update 2025-11-07 01:11:55
	TodoItemPriorityLow TodoItemPriority = "low"

	// TodoItemPriorityMedium TodoItemPriority 中优先级
	//	@update 2025-11-07 01:12:02
	TodoItemPriorityMedium TodoItemPriority = "medium"

	// TodoItemPriorityHigh TodoItemPriority 高优先级
	//	@update 2025-11-07 01:12:09
	TodoItemPriorityHigh TodoItemPriority = "high"

	// TodoItemPriorityUrgent TodoItemPriority 紧急优先级
	//	@update 2025-11-07 01:12:16
	TodoItemPriorityUrgent TodoItemPriority = "urgent"
)
