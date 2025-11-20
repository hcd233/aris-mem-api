package enum

// ObjectStoragePlatform 存储提供商
//
//	@author centonhuang
//	@update 2025-11-20 14:44:28
type ObjectStoragePlatform string

const (

	// ObjectStoragePlatformMinio ObjectStoragePlatform Minio存储
	//	@update 2025-11-20 14:44:17 by centonhuang
	ObjectStoragePlatformMinio ObjectStoragePlatform = "minio"

	// ObjectStoragePlatformCOS ObjectStoragePlatform 腾讯云COS存储
	//	@update 2025-11-20 14:44:17 by centonhuang
	ObjectStoragePlatformCOS ObjectStoragePlatform = "cos"
)
