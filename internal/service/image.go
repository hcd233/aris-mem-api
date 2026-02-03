package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/dto"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/cache"
	"github.com/hcd233/aris-mem-api/internal/infrastructure/pool"
	objdao "github.com/hcd233/aris-mem-api/internal/infrastructure/storage/obj_dao"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/hcd233/aris-mem-api/internal/util"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"go.uber.org/zap"
)

// ImageService 图片服务
//
//	author centonhuang
//	@update 2026-01-31 16:00:00
type ImageService interface {
	UploadImage(ctx context.Context, req *dto.UploadImageReq) (rsp *dto.UploadImageRsp, err error)
	GetCosTempCredential(ctx context.Context, req *dto.EmptyReq) (rsp *dto.GetCosTempCredentialRsp, err error)
}

type imageService struct {
	imageObjDAO objdao.ObjDAO
	stsClient   *sts.Client
}

// NewImageService 创建图片服务
//
//	return ImageService
//	author centonhuang
//	@update 2026-01-31 16:00:00
func NewImageService() ImageService {
	stsClient := sts.NewClient(
		config.CosSecretID,
		config.CosSecretKey,
		nil,
	)

	return &imageService{
		imageObjDAO: objdao.GetImageObjDAO(),
		stsClient:   stsClient,
	}
}

// UploadImage 上传图片
// 使用全局协程池异步处理上传任务
//
//	return *UploadImageRsp
//	author centonhuang
//	@update 2026-01-31 16:00:00
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

// GetCosTempCredential 获取COS临时密钥
//
//	@param ctx context.Context
//	@param req *dto.GetCosTempCredentialReq
//	@return rsp *dto.GetCosTempCredentialRsp
//	@return err error
//	author centonhuang
//	@update 2026-01-31 18:00:00
func (s *imageService) GetCosTempCredential(ctx context.Context, _ *dto.EmptyReq) (*dto.GetCosTempCredentialRsp, error) {
	rsp := &dto.GetCosTempCredentialRsp{}

	logger := logger.WithCtx(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	redisClient := cache.GetRedisClient()

	// 检查缓存中是否已有临时密钥
	cacheKey := fmt.Sprintf(constant.CacheKeyTemplateCosTempSecret, userID)
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中，直接返回
		var credential dto.CosTempCredential
		if err := json.Unmarshal([]byte(cachedData), &credential); err == nil {
			// 检查是否即将过期（剩余时间小于1/4时重新申请）
			if credential.ExpiredTime-time.Now().Unix() > int64(1.0/4*float64(config.CosSTSDuration)) {
				logger.Info("[ImageService] return cached credential", zap.Uint("userID", userID))
				rsp.CosTempCredential = &credential
				return rsp, nil
			}
		}
	}

	// 构建COS资源路径
	// 格式: qcs::cos:<region>:uid/<appid>:<bucketname>/<path>
	resource := fmt.Sprintf("qcs::cos:%s:uid/%s:%s-%s/%s",
		config.CosRegion,
		config.CosAppID,
		config.CosBucketName,
		config.CosAppID,
		fmt.Sprintf("user-%d/image/*", userID),
	)

	// 构建策略
	policy := &sts.CredentialPolicy{
		Version: "2.0",
		Statement: []sts.CredentialPolicyStatement{
			{
				// 允许的权限
				Action: []string{
					// 简单上传
					"name/cos:PutObject",
					"name/cos:PostObject",
					// 分块上传
					"name/cos:InitiateMultipartUpload",
					"name/cos:ListMultipartUploads",
					"name/cos:ListParts",
					"name/cos:UploadPart",
					"name/cos:CompleteMultipartUpload",
					"name/cos:AbortMultipartUpload",
				},
				Effect: "allow",
				Resource: []string{
					resource,
				},
				// 条件限制（可选）
				Condition: map[string]map[string]interface{}{
					"string_equal": {
						// 限制上传的文件类型
						"cos:content-type": []string{
							"image/jpeg",
							"image/jpg",
							"image/png",
							"image/gif",
							"image/webp",
						},
					},
					"numeric_less_than_equal": {
						// 限制文件大小（10MB）
						"cos:content-length": 10 * 1024 * 1024,
					},
				},
			},
		},
	}

	// 申请临时密钥
	opt := &sts.CredentialOptions{
		Policy:          policy,
		Region:          config.CosRegion,
		DurationSeconds: int64(config.CosSTSDuration),
	}

	// 调用STS接口
	credential, err := s.stsClient.GetCredential(opt)
	if err != nil {
		logger.Error("[ImageService] failed to get credential", zap.Error(err))
		rsp.Error = constant.ErrInternalError
		return rsp, nil
	}

	// 构建响应
	rsp.CosTempCredential = &dto.CosTempCredential{
		SecretID:     credential.Credentials.TmpSecretID,
		SecretKey:    credential.Credentials.TmpSecretKey,
		SessionToken: credential.Credentials.SessionToken,
		ExpiredTime:  int64(credential.ExpiredTime),
		Expiration:   credential.Expiration,
		StartTime:    int64(credential.StartTime),
		RequestID:    credential.RequestId,
		BucketName:   config.CosBucketName,
		Region:       config.CosRegion,
		AppID:        config.CosAppID,
	}

	// 缓存临时密钥（有效期减去1分钟作为缓存时间，预留刷新时间）
	cacheDuration := time.Duration(config.CosSTSDuration-60) * time.Second
	if cacheData, err := json.Marshal(rsp.CosTempCredential); err == nil {
		redisClient.Set(ctx, cacheKey, cacheData, cacheDuration)
	}

	logger.Info("[ImageService] credential generated",
		zap.Uint("userID", userID),
		zap.Int("expiredTime", credential.ExpiredTime),
	)

	return rsp, nil
}
