package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddAdmin 加入管理员
func AddAdmin(admin *model.Administrator) {
	util.Db.Create(&admin)
	fmt.Println("增加管理员成功！")
}

// IsHaveAdmin 数据库中是否有该管理员
func IsHaveAdmin(adminName string) bool {
	admin := model.Administrator{}
	if err := util.Db.Where(&model.Administrator{AdminName: adminName}).First(&admin).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Println("不存在该管理员")
		return false
	}
	fmt.Println("存在该管理员")
	return true
}

// SearchAdmin 用主键查找管理员
func SearchAdmin(adminName string) (*model.Administrator, bool){
	admin := model.Administrator{}
	if err := util.Db.Where(&model.Administrator{AdminName: adminName}).First(&admin).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &admin, false
}