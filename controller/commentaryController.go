package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
	"time"
)

type TempCommentary struct {
	PostMasterName string `json:"postMasterName"`
	PostID string `json:"postID"`
	CommentaryDescription string `json:"commentaryDescription"`
}

// CommitCommentaryController 提交评论
func CommitCommentaryController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempCommentary{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	postMasterName := json.PostMasterName
	postID := json.PostID
	commentaryDescription := json.CommentaryDescription

	commentary := model.Commentary{
		CommentaryID: util.CreateUUID("commentary"),
		PostID: postID,
		CommenterName: username,
		CommentaryDescription: commentaryDescription,
		CommentaryTime: time.Now(),
		ReplyNum: 0,
	}
	dao.AddPostCommentaryNum(postID)
	dao.AddCommentary(&commentary)

	data := map[string]interface{} {
		"userName": username,
		"postMasterName": postMasterName,
	}

	util.Success(ctx, data, "评论成功!")
}

// AdminDeleteCommentaryController 管理员删除评论
func AdminDeleteCommentaryController(ctx *gin.Context) {
	commentaryID := ctx.Query("commentaryID")

	if commentaryID == "" {
		util.Fail(ctx, nil, "不存在该评论，无需删除!")
		return
	}

	dao.DeleteCommentary(commentaryID)
	util.Success(ctx, nil, "成功删除该评论")

}