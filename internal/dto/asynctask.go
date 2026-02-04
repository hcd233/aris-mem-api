package dto

import "context"

// EmailSendTask 邮件通知任务
//
//	author centonhuang
//	update 2026-02-04 16:30:00
type EmailSendTask struct {
	Ctx      context.Context
	Emails   []string
	Subject  string
	HTMLBody string
}

// ImageUploadTask 图片上传任务
//
//	@author centonhuang
//	@update 2026-01-31 16:00:00
type ImageUploadTask struct {
	Ctx       context.Context
	ImageName string
	ImageData []byte
}
