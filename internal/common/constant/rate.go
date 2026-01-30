package constant

import "time"

const (

	// PeriodOAuth2Callback OAuth2回调限频周期
	//	@update 2025-11-12 11:27:05
	PeriodOAuth2Callback = 4 * time.Second
	// LimitOAuth2Callback OAuth2回调限频
	//	@update 2025-11-12 11:26:56
	LimitOAuth2Callback = 16

	// PeriodImageUpload Image上传限频周期
	//	@update 2026-01-31 16:00:00
	PeriodImageUpload = 1 * time.Minute
	// LimitImageUpload Image上传限频
	//	@update 2026-01-31 16:00:00
	LimitImageUpload = 16

	// PeriodActionDo Action执行限频周期
	//	@update 2026-01-31 16:00:00
	PeriodActionDo = 4 * time.Second
	// LimitActionDo Action执行限频
	//	@update 2026-01-31 16:00:00
	LimitActionDo = 16

	// PeriodActionUndo Action撤销限频周期
	//	@update 2026-01-31 16:00:00
	PeriodActionUndo = 4 * time.Second
	// LimitActionUndo Action撤销限频
	//	@update 2026-01-31 16:00:00
	LimitActionUndo = 16
)
