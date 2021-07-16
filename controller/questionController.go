package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// -------------------客户端---------------------

// EnterQuestionController 进入问答详情
func EnterQuestionController(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println("token:" + tokenString)
	if strings.HasPrefix(tokenString,"Bearer ") {
		tokenString = tokenString[7:]
	}
	token,claims,errToken := util.ParseToken(tokenString)

	loginName := ""

	if tokenString == "" || errToken != nil || !token.Valid {
		// 游客访问
	} else {
		// 用户访问
		loginName = claims.UserName
	}

	questionID := ctx.Query("questionID")
	fmt.Println("questionID:")
	fmt.Println(questionID)

	// 根据questionID找到对应question，再用username找到对应发布者
	question, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "不存在该问答，问答信息获取失败!")
		return
	}

	questionIssuer, empty := dao.SearchUserByName(question.QuestionIssuerName)
	if empty {
		util.Fail(ctx, nil, "该问答发起人不存在，问答信息有误!")
		return
	}

	var answersData []map[string]interface{}
	// 获取问答的回答
	answers := dao.GetQuestionAnswer(questionID)

	// 获取问答图片
	photoData := dao.GetQuestionPhoto(questionID)

	// 获取最佳答案
	var acceptAnswerData map[string]interface{}
	var acceptAnswer *model.Answer
	if question.QuestionSolved == true { // 如果存在最佳答案
		acceptAnswer = dao.GetAcceptAnswer(question)
		acceptAnswerer, _ := dao.SearchUserByName(acceptAnswer.AnswererName)
		acceptAnswerData = map[string]interface{} {
			"answerID": acceptAnswer.AnswerID,
			"answerer": map[string]interface{} {
				"answererName": acceptAnswerer.UserName,
				"answererPhoto": acceptAnswerer.UserPhoto,
			},
			"answerTime": acceptAnswer.AnswerTime,
			"answerDescription": acceptAnswer.AnswerDescription,
			"answerAcceptance": true,
		}
		answersData = append(answersData, acceptAnswerData)
	}

	for _, answer := range *answers {
		if acceptAnswer != nil && answer.AnswerID == acceptAnswer.AnswerID {
			continue
		}
		// 获取回答者
		answerer, _:= dao.SearchUserByName(answer.AnswererName)
		answerData := map[string]interface{} {
			"answerID": answer.AnswerID,
			"answerer": map[string]interface{}{
				"answererName": answerer.UserName,
				"answererPhoto": answerer.UserPhoto,
			},
			"answerTime": answer.AnswerTime,
			"answerDescription": answer.AnswerDescription,
			"answerAcceptance": answer.AnswerAcceptance,
		}
		answersData = append(answersData, answerData)
	}

	// 访问该问答是否为发布者本人
	var questionSelf = false
	if loginName != "" && loginName == question.QuestionIssuerName {
		questionSelf = true
	}

	data := map[string]interface{}{
		"questionID": questionID,
		"questionIssuer": map[string]interface{}{
			"questionIssuerName": question.QuestionIssuerName,
			"questionIssuerPhoto": questionIssuer.UserPhoto,
		},
		"questionIssueTime": question.QuestionIssueTime,
		"questionDescription": question.QuestionDescription,
		"questionPhotoArray": photoData,
		"questionCircleName": question.QuestionCircleName,
		"questionSolved": question.QuestionSolved,
		"questionSelf": questionSelf,
		"questionAnswerNum": question.QuestionAnswerNum,
		"questionAnswerArray": answersData,
	}
	util.Success(ctx, data, "获取问答信息成功!")
}

// AcceptAnswerAsBestController 接受答案为最佳回答
func AcceptAnswerAsBestController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	questionID := ctx.Query("questionID")
	answerID := ctx.Query("answerID")
	if questionID == "" || answerID == "" {
		util.Fail(ctx, nil, "questionID和answerID缺失!")
		return
	}

	// 判定username是否为问答发起人
	question, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "不存在该问答!")
		return
	}
	if question.QuestionIssuerName != username {
		util.Fail(ctx, nil, "该用户没有权限选择此问题的最佳回答!")
	}

	if !dao.SetQuestionIsSolved(answerID, questionID) {
		util.Fail(ctx, nil, "设定最佳答案失败!")
		return
	}
	util.Success(ctx, nil, "已接该回答为最佳回答!")
}

// DeleteQuestionController 删除问答
func DeleteQuestionController(ctx *gin.Context) {
	username := ctx.GetString("userName")
	questionID := ctx.Query("questionID")

	question, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "待删除问答不存在!")
		return
	}
	if question.QuestionIssuerName != username {
		util.Fail(ctx, nil, "非问答发布者，无权删除该问答!")
		return
	}
	dao.DeleteQuestion(questionID)
	util.Success(ctx, nil, "删除问答成功!")
}

// -------------------管理端---------------------

// GetAllQuestion 获取所有的问答
func GetAllQuestion(ctx *gin.Context) {
	questions := dao.GetAllQuestionInDb()
	var questionsData []map[string]interface{}
	for _, question := range *questions {
		questionPhotoArray := dao.GetQuestionPhoto(question.QuestionID)
		questionData := map[string]interface{} {
			"questionID": question.QuestionID,
			"questionIssuerName": question.QuestionIssuerName,
			"questionIssueTime": question.QuestionIssueTime,
			"questionDescription": question.QuestionDescription,
			"questionPhotoArray": questionPhotoArray,
			"questionAnswerNum": question.QuestionAnswerNum,
			"questionSolved": question.QuestionSolved,
		}
		questionsData = append(questionsData, questionData)
	}

	data := map[string]interface{} {
		"questionTotalNum": len(questionsData),
		"questionArray": questionsData,
	}

	util.Success(ctx, data, "成功获取所有的问答信息!")
}

// AdminDeleteQuestionController 管理员删除问答
func AdminDeleteQuestionController(ctx *gin.Context) {
	questionID := ctx.Query("questionID")
	if questionID == "" {
		util.Fail(ctx, nil, "待删除问答不存在!")
		return
	}
	_, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "待删除问答不存在!")
		return
	}
	dao.DeleteQuestion(questionID)
	util.Success(ctx, nil, "问答删除成功!")
}

// SearchQuestionController 搜索帖子
func SearchQuestionController(ctx *gin.Context) {
	questionID := ctx.Query("questionID")

	question, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "不存在该问答")
		return
	}

	questionIssuer, _ := dao.SearchUserByName(question.QuestionIssuerName)
	questionPhotoArray := dao.GetQuestionPhoto(questionID)

	data := map[string]interface{} {
		"questionID": question.QuestionID,
		"questionIssuer": map[string]interface{} {
			"questionIssuerName": questionIssuer.UserName,
			"questionIssuerPhoto": questionIssuer.UserPhoto,
		},
		"questionIssueTime": question.QuestionIssueTime,
		"questionDescription": question.QuestionDescription,
		"postPhotoArray": questionPhotoArray,
		"questionAnswerNum": question.QuestionAnswerNum,
		"questionSolved": question.QuestionSolved,
	}

	util.Success(ctx, data, "返回问答对应信息")
}


// AdminEnterQuestionController 管理员进入问答详情
func AdminEnterQuestionController(ctx *gin.Context) {
	questionID := ctx.Query("questionID")
	fmt.Println("questionID:")
	fmt.Println(questionID)

	// 根据questionID找到对应question，再用username找到对应发布者
	question, empty := dao.SearchQuestion(questionID)
	if empty {
		util.Fail(ctx, nil, "不存在该问答，问答信息获取失败!")
		return
	}

	questionIssuer, empty := dao.SearchUserByName(question.QuestionIssuerName)
	if empty {
		util.Fail(ctx, nil, "该问答发起人不存在，问答信息有误!")
		return
	}

	var answersData []map[string]interface{}
	// 获取问答的回答
	answers := dao.GetQuestionAnswer(questionID)

	// 获取问答图片
	photoData := dao.GetQuestionPhoto(questionID)

	// 获取最佳答案
	var acceptAnswerData map[string]interface{}
	var acceptAnswer *model.Answer
	if question.QuestionSolved == true { // 如果存在最佳答案
		acceptAnswer = dao.GetAcceptAnswer(question)
		acceptAnswerer, _ := dao.SearchUserByName(acceptAnswer.AnswererName)
		acceptAnswerData = map[string]interface{} {
			"answerID": acceptAnswer.AnswerID,
			"answerer": map[string]interface{} {
				"answererName": acceptAnswerer.UserName,
				"answererPhoto": acceptAnswerer.UserPhoto,
			},
			"answerTime": acceptAnswer.AnswerTime,
			"answerDescription": acceptAnswer.AnswerDescription,
			"answerAcceptance": true,
		}
		answersData = append(answersData, acceptAnswerData)
	}

	for _, answer := range *answers {
		if acceptAnswer != nil && answer.AnswerID == acceptAnswer.AnswerID {
			continue
		}
		// 获取回答者
		answerer, _:= dao.SearchUserByName(answer.AnswererName)
		answerData := map[string]interface{} {
			"answerID": answer.AnswerID,
			"answerer": map[string]interface{}{
				"answererName": answerer.UserName,
				"answererPhoto": answerer.UserPhoto,
			},
			"answerTime": answer.AnswerTime,
			"answerDescription": answer.AnswerDescription,
			"answerAcceptance": answer.AnswerAcceptance,
		}
		answersData = append(answersData, answerData)
	}

	data := map[string]interface{}{
		"questionID": questionID,
		"questionIssuer": map[string]interface{}{
			"questionIssuerName": question.QuestionIssuerName,
			"questionIssuerPhoto": questionIssuer.UserPhoto,
		},
		"questionIssueTime": question.QuestionIssueTime,
		"questionDescription": question.QuestionDescription,
		"questionPhotoArray": photoData,
		"questionCircleName": question.QuestionCircleName,
		"questionSolved": question.QuestionSolved,
		"questionAnswerNum": question.QuestionAnswerNum,
		"questionAnswerArray": answersData,
	}
	util.Success(ctx, data, "获取问答信息成功!")
}