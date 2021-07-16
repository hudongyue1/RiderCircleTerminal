package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type TempAdmin struct {
	AdminName string `json:"adminName"`
	Password string `json:"password"`
}

// 创建管理员账号
func createAdmin(ctx *gin.Context) {
	// 哈希加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return
	}
	admin := model.Administrator{AdminName: "admin", Password: string(hashedPassword)}
	dao.AddAdmin(&admin)
}

// AdminLoginController 管理员登录
func AdminLoginController(ctx *gin.Context) {
	json := TempAdmin{}
	err := ctx.BindJSON(&json)
	if err != nil {
		util.Fail(ctx, nil, "管理员登录失败!")
		return
	}
	adminName := json.AdminName
	password := json.Password

	fmt.Println("adminName:" + adminName)
	fmt.Println("password:" + password)

	if dao.IsHaveAdmin(adminName) {
		admin, _:= dao.SearchAdmin(adminName)
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password));err != nil {
			util.Fail(ctx, nil, "密码错误，请重试!")
			return
		}
		// 密码正确
		// 发放token
		token, err := util.ReleaseToken(adminName)
		if err != nil {panic(err)}

		data := map[string]interface{}{
			"token": token,
			"adminName": adminName,
		}

		util.Success(ctx, data, "密码正确，登录成功!")
		return
	}else {
		util.Fail(ctx, nil, "不存在此管理员账号!")
	}
}