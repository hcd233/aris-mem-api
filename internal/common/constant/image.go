package constant

const (

	// DefaultMaxImageSize 默认最大图片大小
	//	@update 2026-01-30 12:49:53
	DefaultMaxImageSize = 10 * 1024 * 1024

	// DefaultImageQuality 默认图片质量
	//	@update 2026-01-30 12:46:42
	DefaultImageQuality = 90

	// DefaultThumbnailQuality 默认缩略图质量
	DefaultThumbnailQuality = 25

	// DefaultImageExtension 默认图片扩展名（使用 JPEG 保证纯 Go 实现）
	//	@update 2026-01-30 12:46:44
	DefaultImageExtension = ".jpg"
)
