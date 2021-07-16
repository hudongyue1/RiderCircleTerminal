package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
	"time"
)

type TempAnswer struct {
	QuestionMasterName string `json:"questionMasterName"`
	QuestionID string `json:"questionID"`
	AnswerDescription string `json:"answerDescription"`
}

// CommitAnswerController 提交回答
func CommitAnswerController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempAnswer{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	questionMasterName := json.QuestionMasterName
	questionID := json.QuestionID
	answerDescription := json.AnswerDescription
	if questionMasterName == "" || questionID == "" || answerDescription == ""{
		util.Fail(ctx, nil, "回答失败!")
		return
	}

	answer := model.Answer{
		AnswerID: util.CreateUUID("answer"),
		QuestionID: questionID,
		AnswererName: username,
		AnswerTime: time.Now(),
		AnswerDescription: answerDescription,
		AnswerAcceptance: false,
	}
	if !dao.AddQuestionAnswerNum(questionID) {
		util.Fail(ctx, nil, "回答失败!")
		return
	}

	dao.AddAnswer(&answer)

	data := map[string]interface{} {
		"userName": username,
		"questionMasterName": questionMasterName,
	}

	util.Success(ctx, data, "回答成功!")
}

// AdminDeleteAnswerController 管理员删除回答
func AdminDeleteAnswerController(ctx *gin.Context) {
	answerID := ctx.Query("answerID")

	if answerID == "" {
		util.Fail(ctx, nil, "不存在该回答，无需删除!")
		return
	}

	dao.DeleteAnswer(answerID)

	util.Success(ctx, nil, "成功删除该回答")
}