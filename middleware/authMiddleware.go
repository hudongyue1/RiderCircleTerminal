package middleware

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

// UserAuthMiddleware 用户身份验证
func UserAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// token为空
		if tokenString ==""{
			util.Fail(ctx, nil, "没有token!")
			ctx.Abort()
			return
		}

		// 处理带头的token
		if strings.HasPrefix(tokenString,"Bearer ") {
			tokenString = tokenString[7:]
		}

		token,claims,err := util.ParseToken(tokenString)

		if err != nil || !token.Valid {
			util.Fail(ctx, nil, "token错误，请登录")
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的UserName
		username := claims.UserName

		var user model.User
		if err := util.Db.Where(&model.User{UserName: username}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
			fmt.Println("该用户尚未注册")
			util.Fail(ctx, nil, "不存在该用户，请注册")
			ctx.Abort()
			return
		}

		fmt.Println("用户token和身份验证成功")

		//用户存在 将username写入
		ctx.Set("userName", username)
		ctx.Next()
	}
}

// CommonAuthMiddleware 用户或者管理员身份验证
func CommonAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// token为空
		if tokenString ==""{
			util.Fail(ctx, nil, "没有token!")
			ctx.Abort()
			return
		}

		// 处理带头的token
		if strings.HasPrefix(tokenString,"Bearer ") {
			tokenString = tokenString[7:]
		}

		token,claims,err := util.ParseToken(tokenString)

		if err != nil || !token.Valid {
			util.Fail(ctx, nil, "token错误，请登录")
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的UserName
		username := claims.UserName

		var user model.User
		var admin model.Administrator
		if err := util.Db.Where(&model.User{UserName: username}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
			if err := util.Db.Where(model.Administrator{AdminName: username}).First(&admin).Error; gorm.IsRecordNotFoundError(err) {
				util.Fail(ctx, nil, "不存在该管理员，请注册")
				ctx.Abort()
				return
			}
			fmt.Println("管理员token和身份验证成功")
			ctx.Set("admin", username)
			ctx.Next()
		}else {
			fmt.Println("用户token和身份验证成功")
			//用户存在 将username写入
			ctx.Set("userName", username)
			ctx.Next()
		}
	}
}

// AdminAuthMiddleware 用户或者管理员身份验证
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// token为空
		if tokenString ==""{
			util.Fail(ctx, nil, "没有token!")
			ctx.Abort()
			return
		}

		// 处理带头的token
		if strings.HasPrefix(tokenString,"Bearer ") {
			tokenString = tokenString[7:]
		}

		token,claims,err := util.ParseToken(tokenString)

		if err != nil || !token.Valid {
			util.Fail(ctx, nil, "token错误，请登录")
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的UserName
		adminName := claims.UserName

		var admin model.Administrator

		if err := util.Db.Where(model.Administrator{AdminName: adminName}).First(&admin).Error; gorm.IsRecordNotFoundError(err) {
			util.Fail(ctx, nil, "不存在该管理员!")
			ctx.Abort()
			return
		}
		fmt.Println("管理员token和身份验证成功")
		ctx.Set("admin", adminName)
		ctx.Next()
	}
}