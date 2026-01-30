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

	// DefaultMaxCompressedImageSize 默认最大压缩图片大小
	//	@update 2026-01-31 04:07:31
	DefaultMaxCompressedImageSize = 2 * 1024 * 1024
)

// 图片上传协程池配置键
const (
	// ConfigKeyImageUploadPoolMaxWorkers 图片上传协程池最大工作协程数配置键
	//	@update 2026-01-31 16:00:00
	ConfigKeyImageUploadPoolMaxWorkers = "image.upload.pool.max.workers"

	// ConfigKeyImageUploadPoolQueueSize 图片上传协程池任务队列大小配置键
	//	@update 2026-01-31 16:00:00
	ConfigKeyImageUploadPoolQueueSize = "image.upload.pool.queue.size"
)
