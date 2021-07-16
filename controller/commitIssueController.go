package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
	"time"
)

type TempIssue struct {
	CircleName string `json:"circleName"`
	Description string `json:"description"`
	PhotoArray []string `json:"photoArray"`
	Choose int `json:"choose"`
}

func CommitIssueController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempIssue{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	circleName := json.CircleName
	description := json.Description
	photoArray := json.PhotoArray
	choose := json.Choose

	// 用户尚未加入车圈
	if dao.IsUserInCircle(username, circleName) == false {
		util.Fail(ctx, nil, "用户尚未加入该车圈，发布失败!")
		return
	}

	if choose == 0 { // 发帖子
		postID := util.CreateUUID("post")
		post := model.Post{PostID: postID, PostIssuerName: username,
			PostIssueTime: time.Now(), PostDescription: description,
			PostCircleName: circleName, PostUpNum: 0, PostCommentaryNum: 0}
		dao.AddPost(&post, &photoArray)

		// 删除草稿
		dao.DeleteUserDraft(username)

		data := map[string]interface{}{
			"userName": username,
			"postID": postID,
		}

		util.Success(ctx, data, "帖子发布成功")
	}else if choose == 1 { // 发问答
		questionID := util.CreateUUID("question")
		question := model.Question{QuestionID: questionID,
			QuestionIssuerName: username, QuestionIssueTime: time.Now(),
			QuestionDescription: description, QuestionCircleName: circleName,
			QuestionSolved: false, QuestionAnswerNum: 0}
		dao.AddQuestion(&question, &photoArray)

		// 删除草稿
		dao.DeleteUserDraft(username)

		data := map[string]interface{}{
			"userName": username,
			"questionID": questionID,
		}

		util.Success(ctx, data, "问答发布成功")
	}
}

// EnterIssueController 进入发布页面
func EnterIssueController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	circles := dao.GetUserAllCircle(username)

	var circlesData []map[string]interface{}
	for _, circle := range *circles {
		var circlePhoto string
		circlesPhoto := *dao.GetCirclePhoto(circle.CircleName)
		// 车圈没有图片
		if len(circlesPhoto) == 0 {
			circlePhoto = DEFAULT_CIRCLE_PHOTO
		}else {
			// 车圈第一张图为车圈头像
			circlePhoto = circlesPhoto[0]
		}

		circleData := map[string]interface{} {
			"circlePhoto": circlePhoto,
			"circleName": circle.CircleName,
			"circleUserNum": circle.UserNum,
			"circleContentNum": dao.GetCircleContentNum(circle.CircleName),
		}
		circlesData = append(circlesData, circleData)
	}

	data := map[string]interface{}{
		"userName": username,
		"circleArray": circlesData,
	}

	util.Success(ctx, data, "成功获取用户加入的所有车圈!")
}