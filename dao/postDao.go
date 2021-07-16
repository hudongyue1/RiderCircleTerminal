package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddPost 加入Post
func AddPost(post *model.Post, photos *[]string) {
	AddPostPhoto(post.PostID, photos)
	util.Db.Create(&post)
	fmt.Println("增加Post成功！")
}

// AddPostPhoto 给帖子加上PostPhoto
func AddPostPhoto(postID string, photos *[]string) {
	for _, address := range *photos {
		if address != "" {
			postPhoto := model.PostPhoto{PostPhotoID: util.CreateUUID("postPhoto"),
				PostID: postID, PhotoAddress: address}
			util.Db.Create(&postPhoto)
		}
	}
	fmt.Println("增加postPhoto成功！")
}

// DeletePost 用主键删除Post
func DeletePost(postID string) {
	// 删除postPhoto
	util.Db.Where("post_id = ?", postID).Delete(model.PostPhoto{})

	// 删除Post的Commentary
	var commentarys []model.Commentary
	util.Db.Where(&model.Commentary{PostID: postID}).Find(&commentarys)
	for _, commentary := range commentarys {
		DeleteCommentary(commentary.CommentaryID)
	}

	util.Db.Delete(&model.Post{PostID: postID})
	fmt.Println("删除帖子成功！")
}

// SearchPost 用主键查找Post
func SearchPost(postID string) (*model.Post, bool) {
	if postID == "" {
		return nil, true
	}
	post := model.Post{}
	if err := util.Db.Where(&model.Post{PostID: postID}).First(&post).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &post, false
}

// UpdatePost 更新Post
func UpdatePost(post *model.Post) {
	postTemp, _ := SearchPost(post.PostID)
	util.Db.Model(&postTemp).Updates(post)
	fmt.Println("更新帖子信息！")
}

// GetSomeNewPost 获取一些最新的Post
func GetSomeNewPost(num int) *[]model.Post {
	var posts []model.Post
	util.Db.Model(&model.Post{}).Order("post_issue_time").Limit(num).Find(&posts)
	return &posts
}

// GetSomeHotPost 获取一些最热的Post
func GetSomeHotPost(num int) *[]model.Post {
	var posts []model.Post
	util.Db.Model(&model.Post{}).Order("post_commentary_num desc").Limit(num).Find(&posts)
	return &posts
}

// GetSomeNewPostInCircle 获取一些车圈的最新的Post
func GetSomeNewPostInCircle(num int, circleName string) *[]model.Post {
	var posts []model.Post
	util.Db.Where("post_circle_name = ?", circleName).Order("post_issue_time desc").Limit(num).Find(&posts)
	return &posts
}

// GetSomeHotPostInCircle 获取一些车圈的最热的Post
func GetSomeHotPostInCircle(num int, circleName string) *[]model.Post {
	var posts []model.Post
	util.Db.Where("post_circle_name = ?", circleName).Order("post_commentary_num desc").Limit(num).Find(&posts)
	return &posts
}

// GetSomeNewPostForUser 获取一些给用户的最新的Post
func GetSomeNewPostForUser(num int, userName string) *[]model.Post {
	var posts []model.Post
	circles := GetUserAllCircle(userName)
	var circlesData []string
	for _, circle := range *circles {
		circlesData = append(circlesData, circle.CircleName)
	}

	fmt.Println(circlesData)
	util.Db.Where("post_circle_name IN (?)", circlesData).Order("post_issue_time").Limit(num).Find(&posts)

	// 如果用户加入的车圈过少，就推荐一些最新帖子
	if len(posts) <  num {
		temp := GetSomeNewPost(num-len(posts))
		for _, post := range *temp {
			posts = append(posts, post)
		}
	}
	return &posts
}

// GetSomeHotPostForUser 获取一些给用户的最热的Post
func GetSomeHotPostForUser(num int, userName string) *[]model.Post {
	var posts []model.Post
	circles := GetUserAllCircle(userName)
	var circlesData []string
	for _, circle := range *circles {
		circlesData = append(circlesData, circle.CircleName)
	}

	util.Db.Where("post_circle_name IN (?)", circlesData).Order("post_commentary_num desc").Limit(num).Find(&posts)

	// 如果用户加入的车圈过少，就推荐一些其他车圈的最热帖子
	if len(posts) <  num {
		temp := GetSomeHotPost(num-len(posts))
		for _, post := range *temp {
			posts = append(posts, post)
		}
	}

	return &posts
}

// GetPostCommentary 获取帖子的评论
func GetPostCommentary(postID string) *[]model.Commentary {
	var commentaries []model.Commentary
	util.Db.Where(&model.Commentary{PostID: postID}).Find(&commentaries)
	return &commentaries
}

// GetPostPhoto 获取帖子的图片
func GetPostPhoto(postID string) *[]string {
	var photos []model.PostPhoto
	util.Db.Where(&model.PostPhoto{PostID: postID}).Find(&photos)
	var photoData []string
	for _, photo := range photos {
		photoData = append(photoData, photo.PhotoAddress)
	}
	return &photoData
}

// AddPostCommentaryNum 帖子评论数+1
func AddPostCommentaryNum(postID string) {
	post, _ := SearchPost(postID)
	util.Db.Model(&post).Updates(model.Post{PostCommentaryNum: post.PostCommentaryNum+1})
}

// MinusPostCommentaryNum 帖子评论数-1
func MinusPostCommentaryNum(postID string) {
	post, _ := SearchPost(postID)
	util.Db.Model(&post).Update("post_commentary_num", post.PostCommentaryNum-1)
}

// AddPostUpNum 帖子点赞数+1
func AddPostUpNum(postID string, username string) {
	// 增加postUpRelation
	util.Db.Create(&model.PostUpRelation{PostUpRelationID: util.CreateUUID("postUpRelation"),
			PostID: postID,
			UserName: username,
		})

	// 给帖子点赞数+1
	post, _ := SearchPost(postID)
	util.Db.Model(&post).Updates(model.Post{PostUpNum: post.PostUpNum+1})
}

// IsUserUpThePost 用户是否已经给该帖子点赞
func IsUserUpThePost(postID string, username string) bool {
	temp := model.PostUpRelation{}
	if err := util.Db.Where(&model.PostUpRelation{PostID: postID, UserName: username}).First(&temp).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// GetAllPostInDb 获取数据库中所有的帖子
func GetAllPostInDb() *[]model.Post {
	var posts []model.Post
	util.Db.Find(&posts)
	return &posts
}