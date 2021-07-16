package controller

import (
	"RiderCircleTerminal/dao"
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"github.com/gin-gonic/gin"
)

type TempDraft struct {
	Choose int `json:"choose"`
	Description string `json:"description"`
	CircleName string `json:"circleName"`
	PhotoArray []string `json:"photoArray"`
}

// StoreDraftController 保存或修改草稿
func StoreDraftController(ctx *gin.Context) {
	username := ctx.GetString("userName")

	json := TempDraft{}
	err := ctx.BindJSON(&json)
	if err != nil {
		return
	}

	choose := json.Choose
	description := json.Description
	circleName := json.CircleName
	photoArray := json.PhotoArray

	draft, empty := dao.GetUserDraft(username)
	if !empty { // 更新草稿
		draft := model.Draft{
			DraftID: draft.DraftID,
			UserName: draft.UserName,
			Choose: choose,
			Description: description,
			CircleName: circleName,
		}

		dao.UpdateDraft(&draft, &photoArray)

		data := map[string]interface{}{
			"userName": username,
			"draftID": draft.DraftID,
		}
		util.Success(ctx, data, "成功修改草稿")
		return
	} else { // 新建草稿
		draftID := util.CreateUUID("draft")

		draft := model.Draft{
			DraftID:    draftID,
			UserName:    username,
			Choose:      choose,
			Description: description,
			CircleName:  circleName,
		}
		dao.AddDraft(&draft, &photoArray)

		data := map[string]interface{}{
			"userName": username,
			"draftID": draftID,
		}
		util.Success(ctx, data, "成功保存草稿")
	}
}

// GetDraftController 获取草稿
func GetDraftController(ctx *gin.Context)  {
	username := ctx.GetString("userName")
	draft, empty := dao.GetUserDraft(username)
	if empty {
		util.Fail(ctx, nil, "该用户没有草稿!")
		return
	}

	photoArray := dao.GetDraftPhoto(draft.DraftID)

	data := map[string]interface{} {
		"userName": draft.UserName,
		"circleName": draft.CircleName,
		"description": draft.Description,
		"photoArray": photoArray,
		"choose": draft.Choose,
	}
	util.Success(ctx, data, "成功获取用户草稿!")
}


// DeleteDraftController 删除草稿
func DeleteDraftController(ctx *gin.Context)  {
	username := ctx.GetString("userName")

	dao.DeleteUserDraft(username)
	util.Success(ctx, nil, "成功删除草稿!")
}