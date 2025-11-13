package dao

import (
	"github.com/hcd233/aris-mem-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserDAO 用户DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type UserDAO struct {
	baseDAO[model.User]
}

// GetByName 通过用户名获取用户
//
//	receiver dao *UserDAO
//	param db *gorm.DB
//	param name string
//	param fields []string
//	return user *model.User
//	return err error
//	author centonhuang
//	update 2024-10-17 05:18:46
func (dao *UserDAO) GetByName(db *gorm.DB, name string, fields []string) (user *model.User, err error) {
	err = db.Select(fields).Where(model.User{Name: name}).First(&user).Error
	return
}

// GetByGoogleBindID 通过Google绑定ID获取用户
//
//	@receiver dao *UserDAO
//	@param db
//	@param googleBindID
//	@param fields
//	@param preloads
//	@return user
//	@return err
//	@author centonhuang
//	@update 2025-11-13 10:41:10
func (dao *UserDAO) GetByGoogleBindID(db *gorm.DB, googleBindID string, fields []string) (user *model.User, err error) {
	err = db.Select(fields).Where(model.User{GoogleBindID: googleBindID}).First(&user).Error
	return
}

// GetByGithubBindID 通过Github绑定ID获取用户
//
//	@receiver dao *UserDAO
//	@param db
//	@param email
//	@param fields
//	@param preloads
//	@return user
//	@return err
//	@author centonhuang
//	@update 2025-11-13 10:41:17
func (dao *UserDAO) GetByGithubBindID(db *gorm.DB, githubBindID string, fields []string) (user *model.User, err error) {
	err = db.Select(fields).Where(model.User{GithubBindID: githubBindID}).First(&user).Error
	return
}
