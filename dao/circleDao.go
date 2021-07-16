package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// AddCircle 加入Circle
func AddCircle(circle *model.Circle, photos *[]string) {
	AddCirclePhoto(circle.CircleName, photos)
	util.Db.Create(&circle)
	fmt.Println("增加车圈!")
}

// AddCirclePhoto 给帖子加上CirclePhoto
func AddCirclePhoto(circleName string, photos *[]string) {
	for _, address := range *photos {
		if address == "" {
			address = "/public/photo/default/circlePhoto.png"
		}
		circlePhoto := model.CirclePhoto{CirclePhotoID: util.CreateUUID("circlePhoto"),
			CircleName: circleName, PhotoAddress: address}
		util.Db.Create(&circlePhoto)

	}
	fmt.Println("增加车圈图片信息！")
}

// DeleteCircle 用主键删除Circle
func DeleteCircle(circleName string) {
	// 删除circlePhoto
	util.Db.Where("circle_name = ?", circleName).Delete(&model.CirclePhoto{})

	// 删除车圈的所有posts
	var posts []model.Post
	util.Db.Where(model.Post{PostCircleName: circleName}).Find(&posts)
	for _, post := range posts {
		DeletePost(post.PostID)
	}

	// 删除车圈的所有questions
	var questions []model.Question
	util.Db.Where(model.Question{QuestionCircleName: circleName}).Find(&questions)
	for _, question := range questions {
		DeleteQuestion(question.QuestionID)
	}

	// 删除用户加入车圈
	util.Db.Where("circle_name = ?", circleName).Delete(&model.UserRelation{})

	util.Db.Delete(&model.Circle{CircleName: circleName})
	fmt.Println("删除车圈!")
}

// SearchCircle 用主键查找Circle
func SearchCircle(circleName string) (*model.Circle, bool) {
	circle := model.Circle{}
	if err := util.Db.Where(&model.Circle{CircleName: circleName}).First(&circle).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &circle, false
}

// ObscureSearchCircle 模糊搜索车圈
func ObscureSearchCircle(circleName string) (*[]model.Circle, bool) {
	var circles []model.Circle
	if err := util.Db.Where(fmt.Sprintf("circle_name like %q ", "%" + circleName + "%")).Find(&circles).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &circles, false
}

// UpdateCircle 更新Circle
func UpdateCircle(circle *model.Circle, photos *[]string) {
	// 更新circlePhoto
	UpdateCirclePhoto(circle.CircleName, photos)
	circleTemp, _ := SearchCircle(circle.CircleName)
	util.Db.Model(&circleTemp).Updates(circle)
	fmt.Println("更新车圈信息！")
}

// UpdateCirclePhoto 更新circlePhoto
func UpdateCirclePhoto(circleName string, photos *[]string) {
	util.Db.Where("circle_name = ?", circleName).Delete(model.CirclePhoto{})
	AddCirclePhoto(circleName, photos)
	fmt.Println("更新车圈图片信息！")
}

// GetCirclePhoto 获取车圈的图片
func GetCirclePhoto(circleName string) *[]string {
	var photos []model.CirclePhoto
	util.Db.Where(&model.CirclePhoto{CircleName: circleName}).Find(&photos)

	var photoData []string
	for _, photo := range photos {
		photoData = append(photoData, photo.PhotoAddress)
	}
	return &photoData
}

// UserJoinCircle 用户加入车圈
func UserJoinCircle(username string, circleName string) {
	circle, _ := SearchCircle(circleName)
	util.Db.Model(&circle).Updates(model.Circle{UserNum: circle.UserNum+1})

	userRelationID := util.CreateUUID("userRelation")
	util.Db.Create(&model.UserRelation{UserRelationID: userRelationID, UserName: username, CircleName: circleName})
	fmt.Println("用户加入车圈成功!")
}

// IsUserInCircle 用户是否已加入该车圈
func IsUserInCircle(username string, circleName string) bool {
	temp := model.UserRelation{}
	if err := util.Db.Where(&model.UserRelation{UserName: username, CircleName: circleName}).First(&temp).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// UserQuitCircle 用户退出车圈
func UserQuitCircle(username string, circleName string) {
	util.Db.Where("user_name = ? AND circle_name = ?", username, circleName).Delete(model.UserRelation{})
	fmt.Println("用户退出车圈成功!")
}

// GetSomeUserPhotoInCircle 获取车圈的一些用户的头像
func GetSomeUserPhotoInCircle(num int, circleName string) *[]string {
	var usersPhoto []string
	var userRelations []model.UserRelation
	util.Db.Where(&model.UserRelation{CircleName: circleName}).Limit(num).Find(&userRelations)
	for _, userRelation := range userRelations {
		user, _:= SearchUserByName(userRelation.UserName)
		usersPhoto = append(usersPhoto, user.UserPhoto)
	}
	return &usersPhoto
}

// GetAllUserInCircle 获取车圈的所有用户
func GetAllUserInCircle(circleName string) *[]model.User {
	var users []model.User
	var userRelations []model.UserRelation
	util.Db.Where(&model.UserRelation{CircleName: circleName}).Find(&userRelations)
	for _, userRelation := range userRelations {
		user, _:= SearchUserByName(userRelation.UserName)
		users = append(users, *user)
	}
	return &users
}

// GetCirclePostNum 获取车圈帖子数
func GetCirclePostNum(circleName string) int {
	// 查询车圈中帖子数目
	var postNum int
	util.Db.Model(&model.Post{}).Where(&model.Post{PostCircleName: circleName}).Count(&postNum)
	return postNum
}

// GetCircleQuestionNum 获取车圈问答数
func GetCircleQuestionNum(circleName string) int {
	// 查询车圈中问答数目
	var questionNum int
	util.Db.Model(&model.Question{}).Where(&model.Question{QuestionCircleName: circleName}).Count(&questionNum)
	return questionNum
}

// GetCircleContentNum 获取车圈内容条数
func GetCircleContentNum(circleName string) int {
	// 分别查询车圈中帖子和问答的数目
	return GetCircleQuestionNum(circleName) + GetCirclePostNum(circleName)
}

// GetSomeHotCircle 获取一些热门车圈
func GetSomeHotCircle(num int) *[]model.Circle {
	var circles []model.Circle
	util.Db.Model(&model.Circle{}).Order("user_num desc").Limit(num).Find(&circles)
	return &circles
}

// GetAllCircleInDb 获取数据库中所有的车圈
func GetAllCircleInDb() *[]model.Circle {
	var circles []model.Circle
	util.Db.Find(&circles)
	return &circles
}

// GetCircleActiveUserNum 获取车圈活跃用户数目
func GetCircleActiveUserNum(circleName string) int {
	activeUserNum := 0
	timeToday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	users := GetAllUserInCircle(circleName)
	for _, user := range *users {
		if timeToday.Before(user.ActivatedAt) {
			activeUserNum++
		}
	}
	return activeUserNum
}