package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type TempUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	UserPhoto string `json:"userPhoto"`
	UserDescription string `json:"userDescription"`
}


// -------------------用户端---------------------

// CommitInfoController 提交用户信息
func CommitInfoController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempUser{}
	if err := ctx.BindJSON(&json); err != nil {

		return
	}

	photo := json.UserPhoto
	if photo == "" {
		photo = DEFAULT_USER_PHOTO
	}
	description := json.UserDescription

	fmt.Println("用户名:" + username)
	fmt.Println("图片路径:" + photo)
	fmt.Println("个人描述:" + json.UserDescription)

	user := model.User{UserName: username, UserPhoto: photo, UserDescription: description, UpdatedAt: time.Now()}
	dao.UpdateUser(&user)

	data := map[string]interface{}{
		"userName": username,
	}

	util.Success(ctx, data, "修改个人信息成功!")
}

// LoginController Login 用户登录
func LoginController(ctx *gin.Context) {
	json := TempUser{}
	err := ctx.BindJSON(&json)
	if err != nil {
		panic(err)
	}

	username := json.UserName
	password := json.Password

	fmt.Println("userName:" + username)
	fmt.Println("password:" + password)

	// 检验是否已经存在该用户
	if dao.IsHaveUser(username) { // 用户已存在，登录验证
		fmt.Println("用户登录：")
		// 验证密码
		user, _:= dao.SearchUserByName(username)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));err != nil {
			util.Fail(ctx, nil, "密码错误，请重试!")
			return
		}
		// 密码正确
		// 发放token
		token, err := util.ReleaseToken(username)
		if err != nil {panic(err)}

		// 修改用户activeAt
		dao.SetUserActiveAtNow(username)

		data := map[string]interface{}{
			"token": token,
			"userName": username,
		}

		util.Success(ctx, data, "密码正确，登录成功!")
		return
	} else { // 用户不存在，注册用户
		fmt.Println("用户注册：")
		timeNow := time.Now()
		// 加密用户密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
			return
		}
		user := model.User{UserName: username, Password: string(hashedPassword),ActivatedAt: timeNow, CreatedAt: timeNow, UpdatedAt: timeNow}
		dao.AddUser(&user)
		// 发放token
		token, err := util.ReleaseToken(username)
		if err != nil {panic(err)}

		// 修改用户activeAt
		dao.SetUserActiveAtNow(username)

		data := map[string]interface{}{
			"token": token,
			"userName": username,
		}

		util.Success(ctx, data, "注册成功，自动登录!")
		return
	}
}

// EnterMyCountControllerController 进入我的页面，返回用户信息
func EnterMyCountControllerController (ctx *gin.Context) {

	username := ctx.GetString("userName")
	user, _:= dao.SearchUserByName(username)

	data := map[string]interface{}{
		"userName": username,
		"userPhoto": user.UserPhoto,
		"userDescription": user.UserDescription,
	}

	util.Success(ctx, data, "获取个人信息成功!")
}

// EnterMyCommentaryController 进入我的评论页面，返回给用户的评论
func EnterMyCommentaryController (ctx *gin.Context) {
	username := ctx.GetString("userName")

	replys := dao.GetUserAllReply(username)

	var replysData []map[string]interface{}
	for _, reply := range *replys {
		replyer, _:= dao.SearchUserByName(reply.ReplyerName)
		replyData := map[string]interface{} {
			"postID": reply.ReplyID,
			"replyID": reply.ReplyID,
			"replyer": map[string]interface{} {
				"replyerName": reply.ReplyerName,
				"replyerPhoto": replyer.UserPhoto,
			},
			"replyDescription": reply.ReplyDescription,
			"replyTime": reply.ReplyTime,
		}
		replysData = append(replysData, replyData)
	}

	data := map[string]interface{} {
		"userName": username,
		"replyArray": replysData,
	}

	util.Success(ctx, data, "成功获取用户收到的所有评论!")
}

// EnterMyActivityController 进入我的动态页面，返回用户的所有帖子和问题
func EnterMyActivityController (ctx *gin.Context) {
	username := ctx.GetString("userName")

	fmt.Println("进入MyActicvity!")
	fmt.Println("userName:" + username)

	// 获取用户拥有的所有帖子
	posts := dao.GetUserAllPost(username)
	var postsData []map[string]interface{}
	for _, post := range *posts {
		// 获取帖子图片
		photoData := dao.GetPostPhoto(post.PostID)

		postIssuer, _:= dao.SearchUserByName(post.PostIssuerName)
		postData := map[string]interface{} {
			"postID": post.PostID,
			"postIssueTime": post.PostIssueTime,
			"postIssuer": map[string]interface{} {
				"postIssuerName": postIssuer.UserName,
				"postIssuerPhoto": postIssuer.UserPhoto,
			},
			"postDescription": post.PostDescription,
			"postPhotoArray": photoData,
			"postUpNum": post.PostUpNum,
			"postCommentaryNum": post.PostCommentaryNum,
		}
		postsData = append(postsData, postData)
	}

	// 获取用户拥有的所有问答
	questions := dao.GetUserAllQuestion(username)
	var questionsData []map[string]interface{}
	for _, question := range *questions {
		// 获取问答图片
		questionPhotoData := dao.GetQuestionPhoto(question.QuestionID)

		questionIssuer, _:= dao.SearchUserByName(question.QuestionIssuerName)
		questionData := map[string]interface{} {
			"questionID": question.QuestionID,
			"questionIssueTime": question.QuestionIssueTime,
			"questionIssuer": map[string]interface{} {
				"questionIssuerName": questionIssuer.UserName,
				"questionIssuerPhoto": questionIssuer.UserPhoto,
			},
			"questionDescription": question.QuestionDescription,
			"questionSolved": question.QuestionSolved,
			"questionPhotoArray": questionPhotoData,
			"questionAnswerNum": question.QuestionAnswerNum,
		}
		questionsData = append(questionsData, questionData)
	}

	data := map[string]interface{} {
		"userName": username,
		"postArray": postsData,
		"questionArray": questionsData,
	}

	util.Success(ctx, data, "成功获取用户拥有的帖子和问答!")
}

// -------------------管理端---------------------

// GetAllUserController 获取所有的用户信息
func GetAllUserController(ctx *gin.Context) {
	users := dao.GetAllUserInDb()
	var usersData []map[string]interface{}
	for _, user := range *users {
		userData := map[string]interface{} {
			"userName": user.UserName,
			"userPhoto": user.UserPhoto,
			"userDescription": user.UserDescription,
			"userCreatedAt": user.CreatedAt,
		}
		usersData = append(usersData, userData)
	}

	data := map[string]interface{} {
		"userTotalNum": len(usersData),
		"userArray": usersData,
	}

	util.Success(ctx, data, "成功获取所有的用户信息!")
}

// DeleteUserController 删除用户
func DeleteUserController(ctx *gin.Context) {
	username := ctx.Query("userName")
	if !dao.IsHaveUser(username) {
		util.Fail(ctx, nil, "不存在该用户，无需删除!")
		return
	}

	dao.DeleteUserInDb(username)
	util.Success(ctx, nil, "成功删除该用户!")
}

// SearchUserController 查询用户
func SearchUserController(ctx *gin.Context) {
	username := ctx.Query("userName")

	// 首先是精确搜索，如果用户不存在，进行模糊搜索
	user, empty := dao.SearchUserByName(username)

	if empty { // 模糊搜索，返回相似用户
		obscureUsers, empty := dao.ObscureSearchUser(username)
		if empty {
			util.Fail(ctx, nil, "没有相关用户，请尝试别的关键字!")
			return
		}
		var obscureUsersData []map[string]interface{}
		for _, obscureUser := range *obscureUsers {
			obscureUserData := map[string]interface{} {
				"userName": obscureUser.UserName,
				"userPhoto": obscureUser.UserPhoto,
				"userDescription": obscureUser.UserDescription,
				"userCreateTime": obscureUser.CreatedAt,
				"userPostNum": len(*dao.GetUserAllPost(obscureUser.UserName)),
				"userQuestionNum": len(*dao.GetUserAllQuestion(obscureUser.UserName)),
				"userJoinCircleNum": len(*dao.GetUserAllCircle(obscureUser.UserName)),
			}
			obscureUsersData = append(obscureUsersData, obscureUserData)
		}
		data := map[string]interface{} {
			"likeUserNum": len(obscureUsersData),
			"userArray": obscureUsersData,
		}
		util.Success(ctx, data, "找到相关用户!")
		return
	} else { // 精确搜索，返回用户详细信息
		// 获取用户发布的所有帖子和问答

		var usersData []map[string]interface{}
		userData := map[string]interface{} {
			"userName": user.UserName,
			"userPhoto": user.UserPhoto,
			"userDescription": user.UserDescription,
			"userCreateTime": user.CreatedAt,
			"userPostNum": len(*dao.GetUserAllPost(user.UserName)),
			"userQuestionNum": len(*dao.GetUserAllQuestion(user.UserName)),
			"userJoinCircleNum": len(*dao.GetUserAllCircle(user.UserName)),
		}
		usersData = append(usersData, userData)

		data := map[string]interface{} {
			"userArray": usersData,
		}

		util.Success(ctx, data, "成功查询到该用户!")
		return
	}
}