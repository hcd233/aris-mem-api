package enum

type (
	// Permission string 权限
	//	update 2024-09-21 01:34:29
	Permission string
)

const (

	// PermissionReader general permission
	//	update 2024-06-22 10:05:15
	PermissionReader Permission = "reader"

	// PermissionCreator creator permission
	//	update 2024-06-22 10:05:17
	PermissionCreator Permission = "creator"

	// PermissionAdmin admin permission
	//	update 2024-06-22 10:05:17
	PermissionAdmin Permission = "admin"
)

// GetLevel 获取权限等级
//
//	@param p Permission
//	@return int8
//	@author centonhuang
//	@update 2025-11-07 15:05:26
func (p Permission) GetLevel() int8 {
	return permissionLevelMapping[p]
}

// PermissionLevelMapping 权限等级映射
//
//	update 2024-09-21 01:34:29
var (
	permissionLevelMapping = map[Permission]int8{
		PermissionReader:  1,
		PermissionCreator: 2,
		PermissionAdmin:   3,
	}
)
