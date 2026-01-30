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
func ConvertImageToJPEG(imageData []byte, ext string) ([]byte, error) {
	// 解码图片
	img, err := decodeImage(imageData, ext)
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

// decodeImage 根据 MIME 类型解码图片
func decodeImage(data []byte, ext string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(reader)
	case ".png":
		return png.Decode(reader)
	case ".gif":
		return gif.Decode(reader)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
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

// CompressImageToSize 将图片压缩至指定大小以内
// 优先降低质量，如果仍超过限制则按比例缩小尺寸
//
//	@param imageData 原始图片数据
//	@param targetSize 目标大小（字节）
//	@return []byte 压缩后的 JPEG 图片数据
//	@return error
//	@author centonhuang
//	@update 2026-01-31 16:00:00
func CompressImageToSize(imageData []byte, targetSize int) ([]byte, error) {
	// 如果已经小于目标大小，直接返回
	if len(imageData) <= targetSize {
		return imageData, nil
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 移除透明通道
	img = removeTransparency(img)

	// 第一步：尝试降低 JPEG 质量
	for quality := constant.DefaultImageQuality; quality >= 60; quality -= 10 {
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("failed to encode image: %w", err)
		}

		if buf.Len() <= targetSize {
			return buf.Bytes(), nil
		}

	}

	// 第二步：降低质量后仍超过限制，开始缩小尺寸
	// 使用二分查找找到合适的缩放比例
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	low, high := 0.1, 1.0
	var bestResult []byte

	for i := 0; i < 10; i++ { // 最多迭代10次
		scale := (low + high) / 2

		newWidth := int(float64(origWidth) * scale)
		newHeight := int(float64(origHeight) * scale)

		// 确保最小尺寸
		if newWidth < 100 || newHeight < 100 {
			break
		}

		// 缩放图片
		resizedImg := resizeImage(img, newWidth, newHeight)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 75})
		if err != nil {
			return nil, fmt.Errorf("failed to encode resized image: %w", err)
		}

		if buf.Len() <= targetSize {
			bestResult = buf.Bytes()
			low = scale
		} else {
			high = scale
		}
	}

	if bestResult != nil {
		return bestResult, nil
	}

	// 如果二分查找失败，使用最小比例
	scale := 0.1
	newWidth := int(float64(origWidth) * scale)
	newHeight := int(float64(origHeight) * scale)

	if newWidth < 100 {
		newWidth = 100
	}
	if newHeight < 100 {
		newHeight = 100
	}

	resizedImg := resizeImage(img, newWidth, newHeight)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 60})
	if err != nil {
		return nil, fmt.Errorf("failed to encode final resized image: %w", err)
	}

	return buf.Bytes(), nil
}

// resizeImage 按比例缩放图片
func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 创建新的图片
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 简单的最近邻插值算法
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * origWidth / newWidth
			srcY := y * origHeight / newHeight
			resized.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return resized
}
