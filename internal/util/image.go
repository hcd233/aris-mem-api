package util

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/hcd233/aris-mem-api/internal/common/constant"
	"github.com/hcd233/aris-mem-api/internal/common/enum"
)

// ConvertImageToJPEG 将图片转换为统一的目标格式（JPEG）
// 使用 JPEG 保证纯 Go 实现，避免 CGO 依赖
// 如果原图有透明通道，会在白色背景上合成
//
//	@param imageData 原始图片数据
//	@param mimeType 原始图片的 MIME 类型
//	@return []byte 转换后的 JPEG 图片数据
//	@return error
//	@author centonhuang
//	@update 2026-01-30 00:00:00
func ConvertImageToJPEG(imageData []byte, mimeType string) ([]byte, error) {
	// 验证格式
	if !isValidImageFormat(mimeType) {
		return nil, fmt.Errorf("unsupported image format: %s, only supports jpeg, png, gif", mimeType)
	}

	// 解码图片
	img, err := decodeImage(imageData, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 如果图片有透明通道，需要在白色背景上合成
	img = removeTransparency(img)

	// 编码为 JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: constant.DefaultImageQuality})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to JPEG: %w", err)
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
	case enum.ImageFormatJPEG, enum.ImageFormatPNG, enum.ImageFormatGIF:
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
	default:
		return nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}
}

// removeTransparency 移除图片的透明通道，在白色背景上合成
func removeTransparency(img image.Image) image.Image {
	bounds := img.Bounds()

	// 创建一个新的 RGBA 图片，用白色背景
	dst := image.NewRGBA(bounds)

	// 填充白色背景
	draw.Draw(dst, bounds, &image.Uniform{image.White}, image.Point{}, draw.Src)

	// 将原图绘制到白色背景上
	draw.Draw(dst, bounds, img, bounds.Min, draw.Over)

	return dst
}
