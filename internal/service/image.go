package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/dto"
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	"go.uber.org/zap"
)

// ImageService 图片服务
//
//	author centonhuang
//	update 2026-01-31 14:00:00
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
//	update 2026-01-31 14:00:00
func NewImageService() ImageService {
	return &imageService{
		imageObjDAO: objdao.GetImageObjDAO(),
	}
}

// UploadImage 上传图片
//
//	return *UploadImageRsp
//	author centonhuang
//	update 2026-01-31 14:00:00
func (s *imageService) UploadImage(ctx context.Context, req *dto.UploadImageReq) (*dto.UploadImageRsp, error) {
	rsp := &dto.UploadImageRsp{}

	logger := logger.WithCtx(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if req.RawBody.Size > constant.DefaultMaxImageSize {
		logger.Error("[ImageService] image size exceeds limit", zap.Uint("userID", userID), zap.Int64("size", req.RawBody.Size))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	// 解码 base64 或 Data URL
	file, err := req.RawBody.Open()
	if err != nil {
		logger.Error("[ImageService] failed to decode image", zap.Error(err), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}
	defer file.Close()

	ext := filepath.Ext(req.RawBody.Filename)
	ext = strings.ToLower(ext)
	if ext == "" {
		ext = constant.DefaultImageExtension
	}

	imageData, err := io.ReadAll(file)
	if err != nil {
		logger.Error("[ImageService] failed to read image", zap.Error(err), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}
	// 验证并转换图片格式为统一的 JPEG 格式
	imageData, err = util.ConvertImageToJPEG(imageData, ext)
	if err != nil {
		logger.Error("[ImageService] failed to convert image format",
			zap.Error(err),
			zap.Uint("userID", userID),
			zap.String("ext", ext))
		rsp.Error = constant.ErrInvalidFile
		return rsp, nil
	}

	md5Hash := md5.Sum(imageData)
	md5Str := hex.EncodeToString(md5Hash[:])

	// 生成图片文件名
	imageName := fmt.Sprintf("atc-img-%s%s", md5Str[:8], constant.DefaultImageExtension)

	// 上传图片到对象存储
	err = s.imageObjDAO.UploadObject(ctx, userID, imageName, int64(len(imageData)), bytes.NewReader(imageData))
	if err != nil {
		logger.Error("[ImageService] 上传图片失败", zap.Error(err), zap.Uint("userID", userID))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	rsp.ImageName = imageName
	return rsp, nil
}
