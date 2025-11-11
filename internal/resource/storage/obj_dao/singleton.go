package objdao

import (
	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/resource/storage"
)

// createObjectStorageDAO 创建对象存储DAO
func createObjectStorageDAO(objectType ObjectType) ObjDAO {
	switch storage.GetProvider() {
	case storage.ProviderMinio:
		return &MinioObjDAO{
			ObjectType: objectType,
			BucketName: config.MinioBucketName,
			client:     storage.GetMinioStorage(),
		}
	case storage.ProviderCOS:
		return &CosObjDAO{
			ObjectType: objectType,
			BucketName: config.CosBucketName,
			client:     storage.GetCosClient(),
		}
	default:
		panic("unsupported storage type")
	}
}

// GetImageObjDAO 获取图片对象DAO单例
//
//	return ObjDAO
//	author centonhuang
//	update 2024-10-18 01:10:28
func GetImageObjDAO() ObjDAO {
	return createObjectStorageDAO(ObjectTypeImage)
}

// GetThumbnailObjDAO 获取缩略图对象DAO单例
//
//	return ObjDAO
//	author centonhuang
//	update 2024-10-18 01:09:59
func GetThumbnailObjDAO() ObjDAO {
	return createObjectStorageDAO(ObjectTypeThumbnail)
}
