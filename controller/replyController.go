package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
	"time"
)

type TempReply struct {
	CommentaryMasterName string `json:"commentaryMasterName"`
	CommentaryID string `json:"commentaryID"`
	ReplyDescription string `json:"replyDescription"`
}

// CommitReplyController 提交回复
func CommitReplyController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempReply{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	commentaryMasterName := json.CommentaryMasterName
	commentaryID := json.CommentaryID
	replyDescription := json.ReplyDescription

	if commentaryID == "" {
		util.Fail(ctx, nil, "缺少commentaryID，回复失败!")
		return
	}

	reply := model.Reply{
		ReplyID: util.CreateUUID("reply"),
		CommentaryID: commentaryID,
		ReplyerName: username,
		ReplyDescription: replyDescription,
		ReplyTime: time.Now(),
	}

	dao.AddCommentaryReplyNum(commentaryID)
	dao.AddReply(&reply)

	data := map[string]interface{} {
		"userName": username,
		"commentaryMasterName": commentaryMasterName,
	}

	util.Success(ctx, data, "回复成功!")
}

// DeleteReplyController 删除回复
func DeleteReplyController(ctx *gin.Context) {
	username := ctx.GetString("userName")
	replyID := ctx.Query("replyID")

	reply, emptyReply := dao.SearchReply(replyID)
	// 找到对应commenterName
	commentary, emptyCommentary := dao.SearchCommentary(reply.CommentaryID)
	commenterName := commentary.CommenterName

	if emptyReply || emptyCommentary {
		util.Fail(ctx, nil, "待删除回复不存在!")
		return
	}
	if commenterName != username {
		util.Fail(ctx, nil, "非评论者，无权删除该回复!")
		return
	}
	dao.DeleteReply(replyID)

	util.Db.Model(&commentary).Updates(model.Commentary{ReplyNum: commentary.ReplyNum-1})
	util.Success(ctx, nil, "删除回复成功!")
}


// AdminDeleteReplyController 管理员删除回复
func AdminDeleteReplyController(ctx *gin.Context) {
	replyID := ctx.Query("replyID")

	if replyID == "" {
		util.Fail(ctx, nil, "不存在该回复，无需删除!")
		return
	}

	dao.DeleteReply(replyID)
	util.Success(ctx, nil, "删除回复成功!")
}