package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
)

// -------------------用户端---------------------

// EnterPostController 进入帖子详情
func EnterPostController(ctx *gin.Context) {
	postID := ctx.Query("postID")
	// 根据postID找到对应post，再用username找到对应发布者
	post, empty := dao.SearchPost(postID)
	if empty {
		util.Fail(ctx, nil, "不存在该帖子，帖子信息获取失败!")
		return
	}

	postIssuer, empty := dao.SearchUserByName(post.PostIssuerName)
	if empty {
		util.Fail(ctx, nil, "帖子信息有误!")
	}

	// 获取帖子图片
	photoData := dao.GetPostPhoto(postID)

	// 获取帖子评论
	commentarys := dao.GetPostCommentary(postID)
	var commentarysData []map[string]interface{}
	for _, commentary := range *commentarys {
		// 获取评论者
		commenter, _:= dao.SearchUserByName(commentary.CommenterName)
		// 获取评论回复
		replys := dao.GetCommentaryReply(commentary.CommentaryID)
		var replysData []map[string]interface{}
		for _, reply := range *replys {
			replyData := map[string]interface{} {
				"replyID": reply.ReplyID,
				"replyerName": reply.ReplyerName,
				"replyDescription": reply.ReplyDescription,
			}
			replysData = append(replysData, replyData)
		}

		commentaryData := map[string]interface{} {
			"commentaryID": commentary.CommentaryID,
			"commenter": map[string]interface{}{
				"commenterName": commenter.UserName,
				"commenterPhoto": commenter.UserPhoto,
			},
			"commentaryDescription": commentary.CommentaryDescription,
			"commentaryTime": commentary.CommentaryTime,
			"replyNum": commentary.ReplyNum,
			"replyArray": replysData,
		}
		commentarysData = append(commentarysData, commentaryData)
	}

	data := map[string]interface{}{
		"postID": postID,
		"postIssuer": map[string]interface{}{
			"postIssuerName": post.PostIssuerName,
			"postIssuerPhoto": postIssuer.UserPhoto,
		},
		"postIssueTime": post.PostIssueTime,
		"postDescription": post.PostDescription,
		"postPhotoArray": photoData,
		"postCircleName": post.PostCircleName,
		"postUpNum": post.PostUpNum,
		"postCommentaryNum": post.PostCommentaryNum,
		"postCommentaryArray": commentarysData,
	}

	util.Success(ctx, data, "获取帖子信息成功!")
}

// UpPostController 给帖子点赞
func UpPostController(ctx *gin.Context) {
	username := ctx.GetString("userName")

 	postID := ctx.Query("postID")
	dao.AddPostUpNum(postID, username)

	util.Success(ctx, nil, "帖子点赞成功!")
}

// DeletePostController 删除帖子
func DeletePostController(ctx *gin.Context) {
	username := ctx.GetString("userName")
	postID := ctx.Query("postID")

	post, empty := dao.SearchPost(postID)
	if empty {
		util.Fail(ctx, nil, "待删除帖子不存在!")
		return
	}
	if post.PostIssuerName != username {
		util.Fail(ctx, nil, "非帖子发布者，无权删除该帖子!")
		return
	}

	dao.DeletePost(postID)


	util.Success(ctx, nil, "帖子删除成功!")
}

// -------------------管理端---------------------

// GetAllPostController 获取所有的帖子
func GetAllPostController(ctx *gin.Context) {
	posts := dao.GetAllPostInDb()
	var postsData []map[string]interface{}
	for _, post := range *posts {
		postPhotoArray := dao.GetPostPhoto(post.PostID)
		postData := map[string]interface{} {
			"postID": post.PostID,
			"postIssuerName": post.PostIssuerName,
			"postIssueTime": post.PostIssueTime,
			"postDescription": post.PostDescription,
			"postPhotoArray": postPhotoArray,
			"postUpNum": post.PostUpNum,
			"postCommentaryNum": post.PostCommentaryNum,
		}
		postsData = append(postsData, postData)
	}

	data := map[string]interface{} {
		"postTotalNum": len(postsData),
		"postArray": postsData,
	}

	util.Success(ctx, data, "成功获取所有的帖子信息!")
}

// AdminDeletePostController 管理员删除帖子
func AdminDeletePostController(ctx *gin.Context) {
	postID := ctx.Query("postID")
	_, empty := dao.SearchPost(postID)
	if empty {
		util.Fail(ctx, nil, "待删除帖子不存在!")
		return
	}
	dao.DeletePost(postID)
	util.Success(ctx, nil, "帖子删除成功!")
}

// SearchPostController 搜索帖子
func SearchPostController(ctx *gin.Context) {
	postID := ctx.Query("postID")

	post, empty := dao.SearchPost(postID)
	if empty {
		util.Fail(ctx, nil, "不存在该帖子")
		return
	}

	postIssuer, _ := dao.SearchUserByName(post.PostIssuerName)
	postPhotoArray := dao.GetPostPhoto(postID)

	data := map[string]interface{} {
		"postID": post.PostID,
		"postIssuer": map[string]interface{} {
			"postIssuerName": postIssuer.UserName,
			"postIssuerPhoto": postIssuer.UserPhoto,
		},
		"postIssueTime": post.PostIssueTime,
		"postDescription": post.PostDescription,
		"postPhotoArray": postPhotoArray,
		"postUpNum": post.PostUpNum,
		"postCommentaryNum": post.PostCommentaryNum,
	}

	util.Success(ctx, data, "返回帖子对应信息")
}
