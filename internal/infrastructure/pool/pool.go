// Package pool 协程池管理器
//
//	author centonhuang
//	update 2026-02-04 16:10:57
package pool

import (
	"bytes"

	"github.com/alitto/pond/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/smtp"
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"go.uber.org/zap"
)

// PoolManager 全局协程池管理器
//
//	author centonhuang
//	update 2026-01-31 16:00:00
type PoolManager struct {
	imageObjDAO objdao.ObjDAO
	emailClient *smtp.Client

	imageUploadPool       pond.Pool
	emailNotificationPool pond.Pool
}

var poolManager *PoolManager

// InitPoolManager 初始化全局协程池管理器
//
//	@author centonhuang
//	@update 2026-01-31 03:37:28
func InitPoolManager() {
	poolManager = &PoolManager{
		imageObjDAO:           objdao.GetImageObjDAO(),
		emailClient:           smtp.GetEmailClient(),
		imageUploadPool:       pond.NewPool(config.PoolWorkers, pond.WithQueueSize(config.PoolQueueSize)),
		emailNotificationPool: pond.NewPool(config.PoolWorkers, pond.WithQueueSize(config.PoolQueueSize)),
	}
}

// GetPoolManager 获取全局协程池管理器实例
//
//	return *PoolManager
//	author centonhuang
//	update 2026-01-31 16:00:00
func GetPoolManager() *PoolManager {
	return poolManager
}

// StopPoolManager 停止全局协程池管理器
//
//	@author centonhuang
//	@update 2026-01-31 03:47:43
func StopPoolManager() {
	if poolManager != nil {
		poolManager.Stop()
	}
}

// SubmitImageUploadTask InitImageUploadPool 初始化图片上传协程池
//
//	@receiver pm *PoolManager
//	@param task
//	@return error
//	@author centonhuang
//	@update 2026-02-04 16:10:57
func (pm *PoolManager) SubmitImageUploadTask(task *dto.ImageUploadTask) error {
	logger := logger.WithCtx(task.Ctx)
	return pm.imageUploadPool.Go(func() {
		userID := task.Ctx.Value(constant.CtxKeyUserID).(uint)
		err := pm.imageObjDAO.UploadObject(task.Ctx, userID, task.ImageName, int64(len(task.ImageData)), bytes.NewReader(task.ImageData))
		if err != nil {
			logger.Error("[PoolManager] async upload image failed", zap.Error(err), zap.String("imageName", task.ImageName))
			return
		}
		logger.Info("[PoolManager] async upload image success", zap.String("imageName", task.ImageName))
	})
}

// SubmitEmailSendTask 提交新用户注册邮件通知任务
//
//	param task *dto.EmailNotificationTask 邮件通知任务
//	author centonhuang
//	update 2026-02-04 16:30:00
func (pm *PoolManager) SubmitEmailSendTask(task *dto.EmailSendTask) error {
	logger := logger.WithCtx(task.Ctx)
	return pm.emailNotificationPool.Go(func() {
		// 发送邮件
		err := pm.emailClient.SendBatchHTMLMail(task.Emails, task.Subject, task.HTMLBody)
		if err != nil {
			logger.Error("[PoolManager] failed to send email notification",
				zap.Error(err),
				zap.Strings("recipients", task.Emails),
				zap.String("subject", task.Subject))
			return
		}

		logger.Info("[PoolManager] email notification sent successfully",
			zap.Strings("recipients", task.Emails),
			zap.String("subject", task.Subject))
	})
}

// Stop 停止所有协程池
//
//	author centonhuang
//	update 2026-01-31 16:00:00
func (pm *PoolManager) Stop() {
	if pm.imageUploadPool != nil {
		pm.imageUploadPool.Stop()
	}
	if pm.emailNotificationPool != nil {
		pm.emailNotificationPool.Stop()
	}
}
