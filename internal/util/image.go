package util

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// ConvertImageToWebp 将图片转换为统一的目标格式（WebP）
// WebP 支持透明通道，无需额外处理
//
//	@param imageData 原始图片数据
//	@param mimeType 原始图片的 MIME 类型
//	@return []byte 转换后的 WebP 图片数据
//	@return error
//	@author centonhuang
//	@update 2026-01-30 00:00:00
func ConvertImageToWebp(imageData []byte, mimeType string) ([]byte, error) {
	// 验证格式
	if !isValidImageFormat(mimeType) {
		return nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}

	// 解码图片
	img, err := decodeImage(imageData, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 编码为 WebP
	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{
		Lossless: false,
		Quality:  constant.DefaultImageQuality,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to WebP: %w", err)
	}

	return buf.Bytes(), nil
}

// isValidImageFormat 验证图片格式是否支持
//
//	@param mimeType MIME 类型
//	@return bool
//	@author centonhuang
//	@update 2026-01-30 12:44:05
func isValidImageFormat(mimeType string) bool {
	switch enum.ImageFormat(mimeType) {
	case enum.ImageFormatJPEG, enum.ImageFormatPNG, enum.ImageFormatGIF, enum.ImageFormatWebP:
		return true
	default:
		return false
	}
}

// decodeImage 根据 MIME 类型解码图片
func decodeImage(data []byte, mimeType string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch enum.ImageFormat(mimeType) {
	case enum.ImageFormatJPEG:
		return jpeg.Decode(reader)
	case enum.ImageFormatPNG:
		return png.Decode(reader)
	case enum.ImageFormatGIF:
		return gif.Decode(reader)
	case enum.ImageFormatWebP:
		return webp.Decode(reader)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}
}
