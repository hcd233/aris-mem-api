package pool

import (
	"bytes"
	"sync"

	"github.com/alitto/pond/v2"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/dto"
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

	imageUploadPool pond.Pool
}

var (
	poolManager     *PoolManager
	poolManagerOnce sync.Once
)

// InitPoolManager 初始化全局协程池管理器
//
//	@author centonhuang
//	@update 2026-01-31 03:37:28
func InitPoolManager() {
	poolManager = &PoolManager{
		imageObjDAO:     objdao.GetImageObjDAO(),
		imageUploadPool: pond.NewPool(config.PoolWorkers, pond.WithQueueSize(config.PoolQueueSize)),
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

// InitImageUploadPool 初始化图片上传协程池
//
//	@param maxWorkers 最大工作协程数
//	@param queueSize 任务队列大小
//	author centonhuang
//	update 2026-01-31 16:00:00
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

// Stop 停止所有协程池
//
//	author centonhuang
//	update 2026-01-31 16:00:00
func (pm *PoolManager) Stop() {
	if pm.imageUploadPool != nil {
		pm.imageUploadPool.Stop()
	}
}
