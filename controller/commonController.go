package controller

import (
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var (
	UPLOAD_PATH_BASE = "./public/photo"
	DEFAULT_USER_PHOTO = "./public/photo/default/userPhoto.png"
	DEFAULT_CIRCLE_PHOTO = "./public/photo/default/circlePhoto.png"
	DEFAULT_POST_PHOTO = "./public/photo/default/postPhoto.png"
	DEFAULT_QUESTION_PHOTO = "./public/photo/default/questionPhoto.png"
)


// SendAPhotoController 发送一张图片
func SendAPhotoController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	// 图片所属类型 user、circle、post、question
	belong := ctx.PostForm("belong")

	// 获取图片
	form, _ := ctx.MultipartForm()
	file := form.File["file"]
	if len(file) == 0 {
		util.Fail(ctx, nil, "没有要上传的图片")
		return
	}
	photo := file[0]

	// 判断图片所属类型
	if belong != "user" && belong != "circle" && belong != "post" && belong != "question" {
		util.Fail(ctx, nil, "图片所属类型出错，图片上传失败!")
		return
	}

	// 类型检验
	exam := strings.ToLower(path.Ext(photo.Filename))
	if (exam != ".png" && exam != ".jpg" && exam != ".jpeg") || photo == nil {
		util.Fail(ctx, nil, "不支持该格式图片，图片上传失败!")
		return
	}
	log.Println(photo.Filename)

	// 构造不重复的图片名
	dst := UPLOAD_PATH_BASE + "/" + belong + "/" + util.CreateUUID(belong + "Photo") + exam

	// 上传文件到指定的目录
	if err := ctx.SaveUploadedFile(photo, dst); err != nil {
		util.Fail(ctx, nil, "图片路径有问题，图片上传失败!")
		return
	}

	data := map[string]interface{} {
		"userName": username,
		"photo": dst,
		"photoName": photo.Filename,
		"belong": belong,
	}

	util.Success(ctx, data, "图片上传成功!")
}


// SendSomePhotosController 发送一些图片
func SendSomePhotosController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	// 图片所属类型 user、circle、post、question
	belong := ctx.PostForm("belong")

	// 获取图片
	form, _ := ctx.MultipartForm()
	photos := form.File["file"]

	// 判断图片所属类型
	if belong != "user" && belong != "circle" && belong != "post" && belong != "question" {
		util.Fail(ctx, nil, "图片所属类型出错，图片上传失败!")
		return
	}

	// 记录出错图片数
	failCount := 0
	// 记录图片路径
	var photosData []string
	for _, photo := range photos {
		// 类型检验
		exam := strings.ToLower(path.Ext(photo.Filename))
		if (exam != ".png" && exam != ".jpg" && exam != ".jpeg") || photo == nil {
			failCount++
			continue
		}
		log.Println(photo.Filename)

		// 构造不重复的图片名
		dst := UPLOAD_PATH_BASE + "/" + belong + "/" + util.CreateUUID(belong + "Photo") + exam

		// 上传文件到指定的目录
		if err := ctx.SaveUploadedFile(photo, dst); err != nil {
			failCount++
			continue
		}
		photosData = append(photosData, dst)
	}

	data := map[string]interface{} {
		"userName": username,
		"photoArray": photosData,
		"belong": belong,
	}

	msg := fmt.Sprintf("%d张图片，上传失败!  ", failCount) +
		fmt.Sprintf("%d张图片，上传成功!", len(photos)-failCount)
	util.Success(ctx, data, msg)
}

// GetAPhotoController 获取一张图片
func GetAPhotoController(ctx *gin.Context){
	photoUrl := ctx.Query("photo")

	if photoUrl == "" {
		util.Fail(ctx, nil, "图片路径为空!")
		return
	}

	fmt.Println("photoName:" + photoUrl)

	//获取文件名称带后缀
	fileNameWithSuffix := path.Base(photoUrl)
	//获取文件的后缀
	fileType := path.Ext(fileNameWithSuffix)
	if fileType != ".jpg" && fileType != ".jpeg" && fileType != ".png" {
		util.Fail(ctx, nil, "不支持该格式数据!")
		return
	}

	photo, _ := ioutil.ReadFile(photoUrl)
	_, err := ctx.Writer.WriteString(string(photo))
	if err != nil {
		util.Fail(ctx, nil, "获取图片失败!")
		return
	}

	util.Success(ctx, nil, "获取图片成功!")
}


type TempPhotoArray struct {
	PhotoArray []string `json:"photoArray"`
}

// GetSomePhotosController 获取一些图片
func GetSomePhotosController(ctx *gin.Context){
	json := TempPhotoArray{}
	err := ctx.BindJSON(&json)
	if err != nil {
		return
	}

	photoUrlArray := json.PhotoArray
	for _, photoUrl := range photoUrlArray {
		//获取文件名称带后缀
		fileNameWithSuffix := path.Base(photoUrl)
		//获取文件的后缀
		fileType := path.Ext(fileNameWithSuffix)

		if fileType != ".jpg" && fileType != ".jpeg" && fileType != ".png" {
			util.Fail(ctx, nil, "不支持该格式数据!")
			return
		}

		photo, _ := ioutil.ReadFile(photoUrl)
		_, err := ctx.Writer.WriteString(string(photo))
		if err != nil {
			util.Fail(ctx, nil, "获取图片失败!")
			return
		}
	}

	util.Success(ctx, nil, "获取图片成功!")
}

// DeletePhotoController 删除图片
func DeletePhotoController(ctx *gin.Context) {
	photoUrl := ctx.Query("photo")

	// 默认路径不可以删除
	if photoUrl == DEFAULT_USER_PHOTO || photoUrl == DEFAULT_CIRCLE_PHOTO || photoUrl == DEFAULT_POST_PHOTO || photoUrl == DEFAULT_QUESTION_PHOTO {
		return
	}

	err := os.Remove(photoUrl)
	if err != nil {
		util.Fail(ctx, nil, "删除图片失败!")
		return
	}
	util.Success(ctx, nil, "删除图片成功!")
}