package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// DELETED_USER_NAME 已注销用户
var DELETED_USER_NAME = "该用户已注销"

// AddUser 加入User
func AddUser(user *model.User) {
	util.Db.Create(&user)
	fmt.Println("增加用户成功！")
}

// DeleteUserInDb 用主键删除User，需要同时删去用户发布的帖子、问答、草稿，加入的车圈，并将该用户的评论、回复设为该用户已经注销
func DeleteUserInDb(username string) {
	// 删除userRelation
	circles := GetUserAllCircle(username)
	for _, circle := range *circles {
		UserQuitCircle(username, circle.CircleName)
	}

	// 删除用户的Post
	posts := GetUserAllPost(username)
	for _, post := range *posts {
		DeletePost(post.PostID)
	}

	// 删除用户的Question
	questions := GetUserAllQuestion(username)
	for _, question := range *questions {
		DeleteQuestion(question.QuestionID)
	}

	// 删除用户的Draft
	draft, empty := GetUserDraft(username)
	if !empty {
		DeleteDraft(draft.DraftID)
	}

	// 修改用户的commentary，将commenter设为已注销
	commentarys := GetUserAllCommentary(username)
	for _, commentary := range *commentarys {
		SetCommenterDeleted(commentary.CommentaryID)
	}

	// 修改用户的reply，将replyer设为已注销
	replys := GetUserAllReply(username)
	for _, reply := range *replys {
		SetReplyerDeleted(reply.ReplyID)
	}

	// 修改用户的answer，将answerer设为已注销
	answers := GetUserAllAnswer(username)
	for _, answer := range *answers {
		SetAnswererDeleted(answer.AnswerID)
	}

	util.Db.Delete(&model.User{UserName: username})
	fmt.Println("删除用户成功！")
}

// SearchUserByName 用主键查找User
func SearchUserByName(username string) (*model.User, bool){
	user := model.User{}
	if err := util.Db.Where(&model.User{UserName: username}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &user, false
}

// ObscureSearchUser 模糊搜索用户
func ObscureSearchUser(username string) (*[]model.User, bool) {
	var users []model.User
	if err := util.Db.Where(fmt.Sprintf("user_name like %q ", "%" + username + "%")).Find(&users).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &users, false
}

// UpdateUser 更新User
func UpdateUser(user *model.User) {
	userTemp, _ := SearchUserByName(user.UserName)
	util.Db.Model(&userTemp).Updates(user)
	fmt.Println("更新用户信息！")
}

// IsHaveUser 数据库中是否有该用户
func IsHaveUser(username string) bool {
	user := model.User{}
	if err := util.Db.Where(&model.User{UserName: username}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Println("该用户尚未注册")
		return false
	}
	fmt.Println("user.UserName:" + user.UserName)
	fmt.Println("已存在该用户")
	return true
}

// GetUserAllPost 获取用户所有的帖子
func GetUserAllPost(username string) *[]model.Post {
	var posts []model.Post
	util.Db.Where(&model.Post{PostIssuerName: username}).Find(&posts)
	return &posts
}

// GetUserAllQuestion 获取用户所有的问答
func GetUserAllQuestion(username string) *[]model.Question {
	var questions []model.Question
	util.Db.Where(&model.Question{QuestionIssuerName: username}).Find(&questions)
	return &questions
}

// GetUserAllCommentary 获取用户所有的评论
func GetUserAllCommentary(username string) *[]model.Commentary {
	var commentaries []model.Commentary
	util.Db.Where(&model.Commentary{CommenterName: username}).Find(&commentaries)
	return &commentaries
}

// GetUserAllReply 获取用户的所有回复
func GetUserAllReply(username string) *[]model.Reply {
	commentaries := GetUserAllCommentary(username)
	var replys []model.Reply
	for _, commentary := range *commentaries {
		var tempReplys []model.Reply
		util.Db.Where(&model.Reply{CommentaryID: commentary.CommentaryID}).Find(&tempReplys)
		for _, reply := range tempReplys {
			replys = append(replys, reply)
		}
	}
	return &replys
}

// GetUserAllAnswer 获取用户所有的回答
func GetUserAllAnswer(username string) *[]model.Answer {
	var answers []model.Answer
	util.Db.Where(&model.Answer{AnswererName: username}).Find(&answers)
	return &answers
}

// GetUserAllCircle 获取用户加入的所有车圈
func GetUserAllCircle(username string) *[]model.Circle {
	var circles []model.Circle
	var userRelations []model.UserRelation
	util.Db.Where(&model.UserRelation{UserName: username}).Find(&userRelations)
	for _, userRelation := range userRelations {
		circle, _:= SearchCircle(userRelation.CircleName)
		circles = append(circles, *circle)
	}
	return &circles
}

// GetAllUserInDb 获取数据库中所有的用户
func GetAllUserInDb() *[]model.User {
	var users []model.User
	util.Db.Find(&users)
	return &users
}

// SetUserActiveAtNow 修改用户的activeAt时间为现在
func SetUserActiveAtNow(username string) {
	userTemp, _ := SearchUserByName(username)
	util.Db.Model(&userTemp).Updates(model.User{ActivatedAt: time.Now()})
}


