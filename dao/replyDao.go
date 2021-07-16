package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddReply 加入Reply
func AddReply(reply *model.Reply) {
	util.Db.Create(&reply)
	fmt.Println("添加回复成功！")
}

// DeleteReply 用主键删除Reply
func DeleteReply(replyID string) {
	reply, empty := SearchReply(replyID)
	if empty {
		return
	}
	util.Db.Delete(&model.Reply{ReplyID: replyID})

	// 修改对应评论ReplyNum
	MinusCommentaryReplyNum(reply.CommentaryID)

	fmt.Println("删除回复成功！")
}

// SearchReply 用主键查找Reply
func SearchReply(replyID string) (*model.Reply, bool) {
	reply := model.Reply{}
	if err := util.Db.Where(&model.Reply{ReplyID: replyID}).First(&reply).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &reply, false
}

// SetReplyerDeleted 设置该回复的回复者已经注销
func SetReplyerDeleted(replyID string) {
	reply, _ := SearchReply(replyID)
	util.Db.Model(&reply).Updates(model.Reply{ReplyerName: DELETED_USER_NAME})
}