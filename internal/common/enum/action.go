package enum

// ActionEntityType entity type for user actions
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionEntityType string

const (
	// ActionEntityArticle article entity
	//	update 2026-01-30 21:00:00
	ActionEntityArticle ActionEntityType = "article"

	// ActionEntityComment comment entity
	//	update 2026-01-30 16:12:00
	ActionEntityComment ActionEntityType = "comment"
)

// ActionType action type for user actions
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionType string

const (
	// ActionTypeLike like action
	//	update 2026-01-30 21:00:00
	ActionTypeLike ActionType = "like"
	// ActionTypeSave save action
	//	update 2026-01-30 21:00:00
	ActionTypeSave ActionType = "save"
)
