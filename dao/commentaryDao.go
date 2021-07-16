package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddCommentary 加入Commentary
func AddCommentary(commentary *model.Commentary) {
	util.Db.Create(&commentary)
	fmt.Println("添加评论成功！")
}

// DeleteCommentary 用主键删除Commentary
func DeleteCommentary(commentaryID string) {
	commentary, empty := SearchCommentary(commentaryID)
	if empty {
		return
	}
	// 首先删除Commentary的Reply
	util.Db.Where("commentary_id = ?", commentaryID).Delete(&model.Reply{})

	// 修改对应Post的commentaryNum
	MinusPostCommentaryNum(commentary.PostID)

	util.Db.Delete(&model.Commentary{CommentaryID: commentaryID})
	fmt.Println("删除评论成功！")
}

// SearchCommentary 用主键查找Commentary
func SearchCommentary(commentaryID string) (*model.Commentary, bool) {
	commentary := model.Commentary{}
	if err := util.Db.Where(&model.Commentary{CommentaryID: commentaryID}).First(&commentary).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &commentary, false
}

// AddCommentaryReplyNum 给评论回复数+1
func AddCommentaryReplyNum(commentaryID string) {
	commentary, _ := SearchCommentary(commentaryID)
	util.Db.Model(&commentary).Updates(model.Commentary{ReplyNum: commentary.ReplyNum+1})
}

// MinusCommentaryReplyNum 给评论回复数-1
func MinusCommentaryReplyNum(commentaryID string) {
	commentary, _ := SearchCommentary(commentaryID)
	util.Db.Model(&commentary).Update("reply_num", commentary.ReplyNum-1)
}

// GetCommentaryReply 获取评论的回复
func GetCommentaryReply(commentaryID string) *[]model.Reply {
	var replys []model.Reply
	util.Db.Where(&model.Reply{CommentaryID: commentaryID}).Find(&replys)
	return &replys
}

// SetCommenterDeleted 设置该评论的评论者已经注销
func SetCommenterDeleted(commentaryID string) {
	commentary, _ := SearchCommentary(commentaryID)
	util.Db.Model(&commentary).Updates(model.Commentary{CommenterName: DELETED_USER_NAME})
}