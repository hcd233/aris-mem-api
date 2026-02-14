package enum

// NotificationStatus 通知状态
//
//	@author centonhuang
//	@update 2026-02-03 22:30:00
type NotificationStatus string

const (
	// NotificationStatusUnread 未读
	//	@author centonhuang
	//	@update 2026-02-03 22:30:00
	NotificationStatusUnread NotificationStatus = "unread"

	// NotificationStatusRead 已读
	//	@author centonhuang
	//	@update 2026-02-03 22:30:00
	NotificationStatusRead NotificationStatus = "read"
)

// NotificationType 通知类型
//
//	@author centonhuang
//	@update 2026-02-03 22:30:00
type NotificationType string

const (

	// NotificationTypeLike 点赞
	//	@update 2026-02-03 19:18:32
	NotificationTypeLike NotificationType = "like"

	// NotificationTypeSave 收藏
	//	@update 2026-02-03 19:18:31
	NotificationTypeSave NotificationType = "save"

	// NotificationTypeComment 评论
	//	@update 2026-02-03 19:18:28
	NotificationTypeComment NotificationType = "comment"

	// NotificationTypeAt @
	//	@update 2026-02-12 16:00:00
	NotificationTypeAt NotificationType = "at"
)

// NotificationEntityType 通知实体类型
//
//	@author centonhuang
//	@update 2026-02-03 19:17:47
type NotificationEntityType string

const (
	// NotificationEntityTypeArticle 文章
	//	@author centonhuang
	//	@update 2026-02-03 22:30:00
	NotificationEntityTypeArticle NotificationEntityType = "article"

	// NotificationEntityTypeComment 评论
	//	@author centonhuang
	//	@update 2026-02-03 22:30:00
	NotificationEntityTypeComment NotificationEntityType = "comment"
)

// NotificationCategory 通知分类
//
//	@author centonhuang
//	@update 2026-02-12 16:00:00
type NotificationCategory string

const (
	// NotificationCategoryLikeAndSave 点赞和收藏
	//	@author centonhuang
	//	@update 2026-02-12 16:00:00
	NotificationCategoryLikeAndSave NotificationCategory = "likeAndSave"

	// NotificationCategoryCommentAndAt 评论和@
	//	@author centonhuang
	//	@update 2026-02-12 16:00:00
	NotificationCategoryCommentAndAt NotificationCategory = "commentAndAt"
)
