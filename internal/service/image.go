package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/pool"
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// ImageService 图片服务
//
//	author centonhuang
//	update 2026-01-31 16:00:00
type ImageService interface {
	UploadImage(ctx context.Context, req *dto.UploadImageReq) (rsp *dto.UploadImageRsp, err error)
}

type imageService struct {
	imageObjDAO objdao.ObjDAO
}

// NewImageService 创建图片服务
//
//	return ImageService
//	author centonhuang
//	update 2026-01-31 16:00:00
func NewImageService() ImageService {
	return &imageService{
		imageObjDAO: objdao.GetImageObjDAO(),
	}
}

// UploadImage 上传图片
// 使用全局协程池异步处理上传任务
//
//	return *UploadImageRsp
//	author centonhuang
//	update 2026-01-31 16:00:00
func (s *imageService) UploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.UploadImageRsp, error) {
	rsp := &dto.UploadImageRsp{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	logger := logger.WithCtx(ctx)

	image := req.RawBody.Data().Image

	if image.Size > constant.DefaultMaxImageSize {
		logger.Error("[ImageService] image size exceeds limit", zap.Int64("size", image.Size), zap.Int64("expectedSize", constant.DefaultMaxImageSize))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	ext := filepath.Ext(image.Filename)
	ext = strings.ToLower(ext)
	if ext == "" {
		ext = constant.DefaultImageExtension
	}

	imageData, err := io.ReadAll(image.File)
	if err != nil {
		logger.Error("[ImageService] failed to read image", zap.Error(err))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}
	// 验证并转换图片格式为统一的 JPEG 格式
	imageData, err = util.ConvertImageToJPEG(imageData, ext)
	if err != nil {
		logger.Error("[ImageService] failed to convert image format",
			zap.Error(err),
			zap.String("ext", ext))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	// 压缩图片至2MB以内
	imageData, err = util.CompressImageToSize(imageData, constant.DefaultMaxCompressedImageSize)
	if err != nil {
		logger.Error("[ImageService] failed to compress image",
			zap.Error(err))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	logger.Info("[ImageService] image compressed successfully",
		zap.Int64("originalSize", image.Size),
		zap.Int("compressedSize", len(imageData)))

	md5Hash := md5.Sum(imageData)
	md5Str := hex.EncodeToString(md5Hash[:])

	// 生成图片文件名
	imageName := fmt.Sprintf("atc-img-%s%s", md5Str[:8], constant.DefaultImageExtension)

	exists, err := s.imageObjDAO.CheckObject(ctx, userID, imageName)
	if err != nil {
		logger.Error("[ImageService] failed to check object", zap.Error(err))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}
	if exists {
		logger.Info("[ImageService] image already exists", zap.String("imageName", imageName))
		rsp.ImageName = imageName
		return rsp, nil
	}

	// 创建上传任务
	task := &dto.ImageUploadTask{
		Ctx:       util.CopyContextValues(ctx),
		ImageName: imageName,
		ImageData: imageData,
	}

	// 提交到全局协程池异步上传
	poolMgr := pool.GetPoolManager()
	err = poolMgr.SubmitImageUploadTask(task)
	if err != nil {
		logger.Error("[ImageService] failed to submit image upload task", zap.Error(err), zap.String("imageName", imageName))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	rsp.ImageName = imageName
	return rsp, nil
}
