package enum

// Platform 存储提供商
type ObjectStoragePlatform string

const (
	// PlatformMinio Minio存储
	ObjectStoragePlatformMinio ObjectStoragePlatform = "minio"
	// PlatformCOS 腾讯云COS存储
	ObjectStoragePlatformCOS ObjectStoragePlatform = "cos"
)
