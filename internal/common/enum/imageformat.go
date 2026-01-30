package enum

// ImageFormat 支持的图片格式
type ImageFormat string

const (

	// ImageFormatJPEG ImageFormat JPEG图片格式
	//	@update 2026-01-30 12:38:48
	ImageFormatJPEG ImageFormat = "image/jpeg"

	// ImageFormatPNG ImageFormat PNG图片格式
	//	@update 2026-01-30 12:38:48
	ImageFormatPNG ImageFormat = "image/png"

	// ImageFormatGIF ImageFormat GIF图片格式
	//	@update 2026-01-30 12:38:48
	ImageFormatGIF ImageFormat = "image/gif"

	// ImageFormatWebP ImageFormat WebP图片格式
	//	@update 2026-01-30 12:38:48
	ImageFormatWebP ImageFormat = "image/webp"
)
