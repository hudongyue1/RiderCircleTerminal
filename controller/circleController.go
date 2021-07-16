package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type TempCircle struct {
	CircleName string `json:"circleName"`
	CircleMasterName string `json:"circleMasterName"`
	CircleDescription string `json:"circleDescription"`
	CirclePhotoArray []string `json:"circlePhotoArray"`
}

// -------------------客户端---------------------

// EnterCircleHomeController 进入车圈首页（默认获取最热帖子和问答）
func EnterCircleHomeController(ctx *gin.Context) {
	// 获取token
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println("token:" + tokenString)
	if strings.HasPrefix(tokenString,"Bearer ") {
		tokenString = tokenString[7:]
	}
	token,claims,errToken := util.ParseToken(tokenString)

	// 用户未登录则为空
	loginName := ""

	var posts *[]model.Post
	var questions *[]model.Question
	var circles *[]model.Circle

	// 获取帖子和问答数
	sumNum := 10
	numOfPosts := int(sumNum/2)
	numOfQuestions := sumNum - numOfPosts

	// 如果token为空或者无效token，则为游客访问
	if tokenString == "" || errToken != nil || !token.Valid {
		fmt.Println("游客访问!")
		// 游客访问
		// 获取一些最热的帖子
		posts = dao.GetSomeHotPost(numOfPosts)
		// 获取一些最热的问答
		questions = dao.GetSomeHotQuestion(numOfQuestions)
		// 获取一些热门车圈
		circles = dao.GetSomeHotCircle(5)
	} else {
		// 用户访问
		username := claims.UserName
		loginName = username

		fmt.Println("用户"+ username +"访问!")

		// 获取一些用户相关的最热帖子
		posts = dao.GetSomeHotPostForUser(numOfPosts, username)
		// 获取一些用户相关的最热问答
		questions = dao.GetSomeHotQuestionForUser(numOfQuestions, username)
		// 获取一些热门车圈
		circles = dao.GetSomeHotCircle(5)
	}

	data := handleContentInCircleHome(posts, questions, circles, loginName)
	util.Success(ctx, data, "成功获取热门车圈、帖子和问答!")
}

// EnterCircleController 进入车圈(默认获取车圈最热帖子、问答)
func EnterCircleController(ctx *gin.Context) {
	// 获取token
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println("token:" + tokenString)
	if strings.HasPrefix(tokenString,"Bearer ") {
		tokenString = tokenString[7:]
	}
	_,claims,_ := util.ParseToken(tokenString)

	loginName := claims.UserName
	circleName := ctx.Query("circleName")
	fmt.Println("circleName:" + circleName)

	// 一次获取
	sumNum := 10
	numOfPosts := int(sumNum/2)
	numOfQuestions := sumNum - numOfPosts

	// 获取一些最热的帖子
	posts := dao.GetSomeHotPostInCircle(numOfPosts, circleName)
	// 获取一些最热的问答
	questions := dao.GetSomeHotQuestionInCircle(numOfQuestions, circleName)

	data := handleContentInCircle(posts, questions, circleName, loginName)

	util.Success(ctx, data, "欢迎来到车圈"+ circleName +"成功获取热门帖子和问答!")
}

// GetNewContentInCircleController 获取车圈最新的帖子、问答
func GetNewContentInCircleController(ctx *gin.Context) {
	// 获取token
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println("token:" + tokenString)
	if strings.HasPrefix(tokenString,"Bearer ") {
		tokenString = tokenString[7:]
	}
	_,claims,_ := util.ParseToken(tokenString)

	// 用户未登录则为空
	loginName := claims.UserName
	circleName := ctx.Query("circleName")

	// 一次获取
	sumNum := 10
	numOfPosts := int(sumNum/2)
	numOfQuestions := sumNum - numOfPosts

	// 获取一些最新的帖子
	posts := dao.GetSomeNewPostInCircle(numOfPosts, circleName)
	// 获取一些最新的问答
	questions := dao.GetSomeNewQuestionInCircle(numOfQuestions, circleName)

	data := handleContentInCircle(posts, questions, circleName, loginName)

	util.Success(ctx, data, "欢迎来到车圈"+ circleName +"成功获取最新帖子和问答!")
}

// JoinCircleController 用户加入车友圈
func JoinCircleController(ctx *gin.Context) {
	username := ctx.GetString("userName")
	circleName := ctx.Query("circleName")
	if dao.IsUserInCircle(username, circleName) {
		util.Fail(ctx, nil, "用户已经在该车圈中!")
		return
	}

	dao.UserJoinCircle(username, circleName)

	util.Success(ctx, nil, "用户加入车圈成功")
}

// handleContentInCircleHome 处理车圈首页查询结果
func handleContentInCircleHome(posts *[]model.Post, questions *[]model.Question, circles *[]model.Circle, loginName string) map[string]interface{} {
	var circlesData []map[string]interface{}
	var postsData []map[string]interface{}
	var questionsData []map[string]interface{}

	// 处理post数据
	for _, post := range *posts {
		// 获取帖子图片
		postPhotoData := dao.GetPostPhoto(post.PostID)

		// 是否已经点赞
		var isUped = false
		if loginName != "" && dao.IsUserUpThePost(post.PostID, loginName) {
			isUped = true
		}

		postIssuer, _:= dao.SearchUserByName(post.PostIssuerName)
		postData := map[string]interface{} {
			"postID": post.PostID,
			"isUped": isUped,
			"postIssuer": map[string]interface{} {
				"postIssuerName": postIssuer.UserName,
				"postIssuerPhoto": postIssuer.UserPhoto,
			},
			"postIssueTime": post.PostIssueTime,
			"postDescription": post.PostDescription,
			"postPhotoArray": postPhotoData,
			"postCircleName": post.PostCircleName,
			"postUpNum": post.PostUpNum,
			"postCommentaryNum": post.PostCommentaryNum,
		}
		postsData = append(postsData, postData)
	}

	// 处理question数据
	for _, question := range *questions {
		// 获取问答图片
		questionPhotoData := dao.GetQuestionPhoto(question.QuestionID)
		questionIssuer, _:= dao.SearchUserByName(question.QuestionIssuerName)
		// 处理最佳答案
		var acceptAnswerData map[string]interface{}
		if question.QuestionSolved == true {
			acceptAnswer := dao.GetAcceptAnswer(&question)
			acceptAnswerer, _ := dao.SearchUserByName(acceptAnswer.AnswererName)
			acceptAnswerData = map[string]interface{} {
				"answerID": acceptAnswer.AnswerID,
				"answerer": map[string]interface{} {
					"answererName": acceptAnswerer.UserName,
					"answererPhoto": acceptAnswerer.UserPhoto,
				},
				"answerTime": acceptAnswer.AnswerTime,
				"answerDescription": acceptAnswer.AnswerDescription,
			}
		}

		questionData := map[string]interface{} {
			"questionID": question.QuestionID,
			"questionIssuer": map[string]interface{} {
				"questionIssuerName": questionIssuer.UserName,
				"questionIssuerPhoto": questionIssuer.UserPhoto,
			},
			"questionIssueTime": question.QuestionIssueTime,
			"questionDescription": question.QuestionDescription,
			"questionCircleName": question.QuestionCircleName,
			"questionSolved": question.QuestionSolved,
			"questionPhotoArray": questionPhotoData,
			"questionAnswerNum": question.QuestionAnswerNum,
			"acceptAnswer": acceptAnswerData,
		}
		questionsData = append(questionsData, questionData)
	}

	// 处理circle数据
	for _, circle := range *circles {
		// 获取车圈图片
		circlePhotoData := dao.GetCirclePhoto(circle.CircleName)
		var isJoined = false
		if loginName != "" && dao.IsUserInCircle(loginName, circle.CircleName) {
			isJoined = true
		}
		circleMaster, _ := dao.SearchUserByName(circle.CircleMasterName)
		usersPhoto := dao.GetSomeUserPhotoInCircle(3, circle.CircleName)

		circleData := map[string]interface{} {
			"circleName": circle.CircleName,
			"isJoined": isJoined,
			"circleMasterPhoto": circleMaster.UserPhoto,
			"circleUserPhoto": usersPhoto,
			"circleActiveUserNum": dao.GetCircleActiveUserNum(circle.CircleName),
			"circleUserNum": circle.UserNum,
			"circlePhotoArray": circlePhotoData,
		}
		circlesData = append(circlesData, circleData)
	}

	data := map[string]interface{} {
		"hotCircleArray": circlesData,
		"postArray": postsData,
		"questionArray": questionsData,
	}
	return data
}

// handleContentInCircle 处理车圈查询结果
func handleContentInCircle(posts *[]model.Post, questions *[]model.Question, circleName string, loginName string) map[string]interface{} {
	var postsData []map[string]interface{}
	var questionsData []map[string]interface{}
	circle, _ := dao.SearchCircle(circleName)

	// 处理post数据
	for _, post := range *posts {
		// 获取帖子图片
		postPhotoData := dao.GetPostPhoto(post.PostID)

		// 是否已经点赞
		var isUped = false
		if loginName != "" && dao.IsUserUpThePost(post.PostID, loginName) {
			isUped = true
		}

		postIssuer, _:= dao.SearchUserByName(post.PostIssuerName)
		postData := map[string]interface{} {
			"postID": post.PostID,
			"isUped": isUped,
			"postIssuer": map[string]interface{} {
				"postIssuerName": postIssuer.UserName,
				"postIssuerPhoto": postIssuer.UserPhoto,
			},
			"postIssueTime": post.PostIssueTime,
			"postDescription": post.PostDescription,
			"postPhotoArray": postPhotoData,
			"postUpNum": post.PostUpNum,
			"postCommentaryNum": post.PostCommentaryNum,
		}
		postsData = append(postsData, postData)
	}

	// 处理question数据
	for _, question := range *questions {
		// 获取问答图片
		questionPhotoData := dao.GetQuestionPhoto(question.QuestionID)
		questionIssuer, _:= dao.SearchUserByName(question.QuestionIssuerName)
		// 获取最佳答案
		var acceptAnswerData map[string]interface{}
		if question.QuestionSolved == true {
			acceptAnswer := dao.GetAcceptAnswer(&question)
			acceptAnswerer, _ := dao.SearchUserByName(acceptAnswer.AnswererName)
			acceptAnswerData = map[string]interface{} {
				"answerID": acceptAnswer.AnswerID,
				"answerer": map[string]interface{} {
					"answererName": acceptAnswerer.UserName,
					"answererPhoto": acceptAnswerer.UserPhoto,
				},
				"answerTime": acceptAnswer.AnswerTime,
				"answerDescription": acceptAnswer.AnswerDescription,
			}
		}

		questionData := map[string]interface{} {
			"questionID": question.QuestionID,
			"questionIssuer": map[string]interface{} {
				"questionIssuerName": questionIssuer.UserName,
				"questionIssuerPhoto": questionIssuer.UserPhoto,
			},
			"questionIssueTime": question.QuestionIssueTime,
			"questionDescription": question.QuestionDescription,
			"questionCircleName": question.QuestionCircleName,
			"questionSolved": question.QuestionSolved,
			"questionPhotoArray": questionPhotoData,
			"questionAnswerNum": question.QuestionAnswerNum,
			"acceptAnswer": acceptAnswerData,
		}
		questionsData = append(questionsData, questionData)
	}

	circleMaster, _ := dao.SearchUserByName(circle.CircleMasterName)
	circleBackground := (*dao.GetCirclePhoto(circle.CircleName))[0]
	usersPhoto := dao.GetSomeUserPhotoInCircle(3, circleName)

	var isJoined = false
	if loginName != "" && dao.IsUserInCircle(loginName, circle.CircleName) {
		isJoined = true
	}

	data := map[string]interface{} {
		"circleName": circle.CircleName,
		"circleDescription": circle.CircleDescription,
		"circleUserNum": circle.UserNum,
		"isJoined": isJoined,
		"circleMasterPhoto": circleMaster.UserPhoto,
		"circleBackGround": circleBackground,
		"circleUserPhoto": usersPhoto,
		"questionArray": questionsData,
		"postArray": postsData,
	}

	return data
}

// -------------------管理端---------------------

// AddOrUpdateCircleController 创建或者修改车圈
func AddOrUpdateCircleController(ctx *gin.Context) {
	json := TempCircle{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	circleName := json.CircleName
	circleMasterName := json.CircleMasterName
	circleDescription := json.CircleDescription
	photos := json.CirclePhotoArray

	var circle model.Circle
	fmt.Println(photos)

	// 判断车圈是否已经存在
	tempCircle, empty := dao.SearchCircle(circleName)

	// 判断车圈主是否存在
	if !dao.IsHaveUser(circleMasterName) {
		util.Fail(ctx, nil, "不存在选为车圈主的用户!")
		return
	}

	if empty { // 创建车圈
		circle = model.Circle{CircleName: circleName, CircleMasterName: circleMasterName,
			UserNum: 0, CircleDescription: circleDescription}

		dao.AddCircle(&circle, &photos)
		dao.UserJoinCircle(circleMasterName, circleName)

		data := map[string]interface{} {
			"circleMasterName": circleMasterName,
			"circleName": circleName,
		}
		util.Success(ctx, data, "创建车圈成功!")
		return
	} else { // 修改车圈
		circle = model.Circle{CircleName: circleName, CircleMasterName: circleMasterName,
			UserNum: tempCircle.UserNum, CircleDescription: circleDescription}

		dao.UpdateCircle(&circle, &photos)
		if !dao.IsUserInCircle(circleMasterName, circleName) {
			dao.UserJoinCircle(circleMasterName, circleName)
		}
		data := map[string]interface{} {
			"circleMasterName": circleMasterName,
			"circleName": circleName,
		}
		util.Success(ctx, data, "修改车圈成功!")
	}
}

// GetAllCircleController 获取所有的车圈
func GetAllCircleController(ctx *gin.Context) {
	circles := dao.GetAllCircleInDb()
	var circlesData []map[string]interface{}
	for _, circle := range *circles {
		circlePhotoArray := dao.GetCirclePhoto(circle.CircleName)
		circleData := map[string]interface{} {
			"circleName": circle.CircleName,
			"circleMasterName": circle.CircleMasterName,
			"circleDescription": circle.CircleDescription,
			"circlePhotoArray": circlePhotoArray,
			"circleCreateTime": circle.CreatedAt,
			"circleUserNum": circle.UserNum,
			"circleActiveUserNum": dao.GetCircleActiveUserNum(circle.CircleName),
		}
		circlesData = append(circlesData, circleData)
	}

	data := map[string]interface{} {
		"circleTotalNum": len(circlesData),
		"circleArray": circlesData,
	}

	util.Success(ctx, data, "成功获取所有的车圈信息!")
}

// DeleteCircleController 删除车圈
func DeleteCircleController(ctx *gin.Context) {
	circleName := ctx.Query("circleName")

	_, empty := dao.SearchCircle(circleName)
	if empty {
		util.Fail(ctx, nil, "不存在该车圈!")
		return
	}

	dao.DeleteCircle(circleName)
	util.Success(ctx, nil, "删除车圈成功!")
}

// SearchCircleController 查询车圈
func SearchCircleController(ctx *gin.Context) {
	circleName := ctx.Query("circleName")

	// 首先进行精确搜索，如果不存在，则进行模糊搜索
	circle, empty := dao.SearchCircle(circleName)

	if empty { // 模糊搜索
		obscureCircles, empty := dao.ObscureSearchCircle(circleName)
		if empty {
			util.Fail(ctx, nil, "没有相关车圈，请尝试别的关键字!")
			return
		}
		var obscureCirclesData []map[string]interface{}
		for _, obscureCircle := range *obscureCircles {
			obscureCirclePhoto := dao.GetCirclePhoto(obscureCircle.CircleName)
			obscureCircleData := map[string]interface{} {
				"circleName": obscureCircle.CircleName,
				"circleMasterName": obscureCircle.CircleMasterName,
				"circleDescription": obscureCircle.CircleDescription,
				"circlePhotoArray": obscureCirclePhoto,
				"circleCreateTime": obscureCircle.CreatedAt,
				"circleUserNum": obscureCircle.UserNum,
				"circleActiveUserNum": dao.GetCircleActiveUserNum(obscureCircle.CircleName),
			}
			obscureCirclesData = append(obscureCirclesData, obscureCircleData)
		}
		data := map[string]interface{} {
			"likeCircleNum": len(obscureCirclesData),
			"circleArray": obscureCirclesData,
		}
		util.Success(ctx, data, "找到相关车圈!")
		return
	} else { // 精确搜索
		circlePhotoArray := dao.GetCirclePhoto(circleName)

		var circlesData []map[string]interface{}

		circleData := map[string]interface{} {
			"circleName": circle.CircleName,
			"circleMasterName": circle.CircleMasterName,
			"circleDescription": circle.CircleDescription,
			"circlePhotoArray": circlePhotoArray,
			"circleCreateTime": circle.CreatedAt,
			"circleUserNum": circle.UserNum,
			"circleActiveUserNum": dao.GetCircleActiveUserNum(circle.CircleName),
		}
		circlesData = append(circlesData, circleData)

		data := map[string]interface{} {
			"circleArray": circlesData,
		}
		util.Success(ctx, data, "返回车圈对应信息!")
	}
}