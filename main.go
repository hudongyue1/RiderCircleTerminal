package main

import (
	"RiderCircleTerminal/controller"
	"RiderCircleTerminal/middleware"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)



/* 待完成

*/

func main() {
	util.InitDB()
	defer func(Db *gorm.DB) {
		err := Db.Close()
		if err != nil {
			panic(err)
		}
	}(util.Db)


	r := gin.Default()

	// 通用接口
	common := r.Group("/common")
	{
		// 发送一张图片
		common.POST("/sendAPhoto", middleware.CommonAuthMiddleware(), controller.SendAPhotoController)

		// 获取一张图片
		common.GET("/getAPhoto", controller.GetAPhotoController)

		// 发送多张图片
		common.POST("/sendSomePhotos", middleware.CommonAuthMiddleware(), controller.SendSomePhotosController)

		// 获取多张图片
		common.GET("/getSomePhotos", controller.GetSomePhotosController)

		// 删除图片
		common.DELETE("/deletePhoto",middleware.CommonAuthMiddleware(), controller.DeletePhotoController)
	}

	// 客户端接口
	client := r.Group("/client")
	{
		// 车圈首页
		circleHome := client.Group("/circleHome")
		{
			// 进入车友圈首页(默认获取最热帖子、问答)，未登录为游客访问
			circleHome.GET("/enterCircleHome", controller.EnterCircleHomeController)

			// 加入车友圈
			circleHome.POST("/joinCircle", middleware.UserAuthMiddleware(), controller.JoinCircleController)
		}

		// 页面"我的"
		myCount := client.Group("/myCount")
		{
			myCount.Use(middleware.UserAuthMiddleware())

			// 进入我的
			myCount.GET("/enterMyCount", controller.EnterMyCountControllerController)
		}

		// 车友圈
		circle := client.Group("/circle")
		{
			// 加入车友圈
			circle.POST("/joinCircle", middleware.UserAuthMiddleware(), controller.JoinCircleController)

			// 进入车友圈(默认获取最热帖子、问答)
			circle.GET("/enterCircle", controller.EnterCircleController)

			// 获取最新的帖子、问答)
			circle.GET("/getNewContentInCircle", controller.GetNewContentInCircleController)
		}

		// 发布
		issue := client.Group("/issue")
		{
			issue.Use(middleware.UserAuthMiddleware())

			// 进入发布页面
			issue.GET("/enterIssue", controller.EnterIssueController)

			// 发布帖子或者问答
			issue.POST("/commitIssue", controller.CommitIssueController)

			// 储存草稿
			issue.POST("/storeDraft", controller.StoreDraftController)

			// 储存草稿
			issue.GET("/getDraft", controller.GetDraftController)

			// 删除草稿
			issue.DELETE("/deleteDraft", controller.DeleteDraftController)
		}

		// 登录
		login := client.Group("/login")
		{
			// 登录
			login.POST("/commitLogin", controller.LoginController)
		}

		// 我的动态
		myActivity := client.Group("/myActivity")
		{
			myActivity.Use(middleware.UserAuthMiddleware())

			// 进入我的动态
			myActivity.GET("/enterMyActivity", controller.EnterMyActivityController)

			// 删除帖子
			myActivity.DELETE("/deletePost", controller.DeletePostController)

			// 删除问答
			myActivity.DELETE("/deleteQuestion", controller.DeleteQuestionController)
		}

		// 我的评论
		myCommentary := client.Group("/myCommentary")
		{
			myCommentary.Use(middleware.UserAuthMiddleware())

			// 进入我的回复
			myCommentary.GET("/enterMyCommentary", controller.EnterMyCommentaryController)

			// 删除回复
			myCommentary.DELETE("/deleteReply", controller.DeleteReplyController)
		}

		// 帖子详情
		post := client.Group("/post")
		{
			// 进入帖子详情
			post.GET("/enterPost", controller.EnterPostController)

			// 发表评论
			post.POST("/commitCommentary", middleware.UserAuthMiddleware(), controller.CommitCommentaryController)

			// 给点赞帖子
			post.POST("/upPost", middleware.UserAuthMiddleware(), controller.UpPostController)

			// 回复评论
			post.POST("/commitReply", middleware.UserAuthMiddleware(), controller.CommitReplyController)
		}

		// 问答详情
		question := client.Group("/question")
		{
			// 进入问答详情
			question.GET("/enterQuestion", controller.EnterQuestionController)

			// 发表回答
			question.POST("/commitAnswer", middleware.UserAuthMiddleware(), controller.CommitAnswerController)

			// 选择回答为最佳回答
			question.POST("/acceptAnswerAsBest", middleware.UserAuthMiddleware(), controller.AcceptAnswerAsBestController)
		}

		// 个人信息
		personalInfo := client.Group("/personalInfo")
		{
			personalInfo.Use(middleware.UserAuthMiddleware())
			// 提交个人信息
			personalInfo.POST("/commitInfo", controller.CommitInfoController)
		}
	}

	// 管理员端接口
	admin := r.Group("/admin")
	{
		// 管理员登录
		admin.POST("/adminLogin", controller.AdminLoginController)

		// 用户管理
		userManage := admin.Group("/userManage")
		{
			userManage.Use(middleware.AdminAuthMiddleware())

			// 获取所有用户信息
			userManage.GET("/getAllUser", controller.GetAllUserController)

			//// 增加或者修改用户
			//userManage.POST("/addOrUpdateUser", controller.AddOrUpdateUser)

			// 删除用户（会删除掉用户加入的车圈信息、发布的帖子和问答，并且用户的评论、回答、回复会设定为用户已注销）
			userManage.DELETE("/deleteUser", controller.DeleteUserController)

			// 搜索用户(精确搜索、模糊搜索)
			userManage.GET("/searchUser", controller.SearchUserController)
		}

		// 车圈管理
		circleManage := admin.Group("/circleManage")
		{
			circleManage.Use(middleware.AdminAuthMiddleware())

			// 获取所有车圈信息
			circleManage.GET("/getAllCircle", controller.GetAllCircleController)

			// 增加或者修改车圈
			circleManage.POST("/addOrUpdateCircle", controller.AddOrUpdateCircleController)

			// 删除车圈
			circleManage.DELETE("/deleteCircle", controller.DeleteCircleController)

			// 搜索车圈(精确搜索、模糊搜索)
			circleManage.GET("/searchCircle", controller.SearchCircleController)
		}

		// 帖子管理
		postManage := admin.Group("/postManage")
		{
			postManage.Use(middleware.AdminAuthMiddleware())

			// 获取所有帖子
			postManage.GET("/getAllPost", controller.GetAllPostController)

			// 删除帖子
			postManage.DELETE("/deletePost", controller.AdminDeletePostController)

			// 搜索帖子
			postManage.GET("/searchPost", controller.SearchPostController)

			// 进入帖子详情
			postManage.GET("/enterPost", controller.EnterPostController)

			// 删除评论
			postManage.DELETE("/deleteCommentary", controller.AdminDeleteCommentaryController)

			// 删除回复
			postManage.DELETE("/deleteReply", controller.AdminDeleteReplyController)
		}

		// 问答管理
		questionManage := admin.Group("/questionManage")
		{
			questionManage.Use(middleware.AdminAuthMiddleware())

			// 获取所有问答
			questionManage.GET("/getAllQuestion", controller.GetAllQuestion)

			// 删除问答
			questionManage.DELETE("/deleteQuestion", controller.AdminDeleteQuestionController)

			// 搜索问答
			questionManage.GET("/searchQuestion", controller.SearchQuestionController)

			// 进入问答详情
			questionManage.GET("/enterQuestion", controller.AdminEnterQuestionController)

			// 删除回答
			questionManage.DELETE("/deleteAnswer", controller.AdminDeleteAnswerController)
		}
	}

	err := r.Run()
	if err != nil {
		panic(err.Error())
	}
}