package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddDraft 加入草稿
func AddDraft(draft *model.Draft, photos *[]string) {
	AddDraftPhoto(draft.DraftID, photos)
	util.Db.Create(&draft)
	fmt.Println("增加草稿成功！")
}

// AddDraftPhoto 给草稿加上DraftPhoto
func AddDraftPhoto(draftID string, photos *[]string) {
	for _, address := range *photos {
		if address != "" {
			draftPhoto := model.DraftPhoto{DraftPhotoID: util.CreateUUID("draftPhoto"),
				DraftID: draftID, PhotoAddress: address}
			util.Db.Create(&draftPhoto)
		}
	}
	fmt.Println("增加draftPhoto成功！")
}

// DeleteDraft 用主键删除草稿
func DeleteDraft(draftID string) {
	// 删除草稿的draftPhoto
	util.Db.Where("draft_id = ?", draftID).Delete(model.DraftPhoto{})

	util.Db.Delete(&model.Draft{DraftID: draftID})
	fmt.Println("删除草稿成功！")
}

// UpdateDraft 更新Draft
func UpdateDraft(draft *model.Draft, photos *[]string) {
	// 更新circlePhoto
	UpdateDraftPhoto(draft.DraftID, photos)
	draftTemp, _ := SearchDraft(draft.DraftID)
	util.Db.Model(&draftTemp).Updates(draft)
	fmt.Println("更新草稿信息！")
}

// UpdateDraftPhoto 更新DraftPhoto
func UpdateDraftPhoto(draftID string, photos *[]string) {
	util.Db.Where("draft_id = ?", draftID).Delete(model.DraftPhoto{})
	AddDraftPhoto(draftID, photos)
	fmt.Println("更新草稿图片信息！")
}

// SearchDraft 用主键查找草稿
func SearchDraft(draftID string) (*model.Draft, bool) {
	if draftID == "" {
		return nil, true
	}
	draft := model.Draft{}
	if err := util.Db.Where(&model.Draft{DraftID: draftID}).First(&draft).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &draft, false
}

// GetUserDraft 获取用户的草稿
func GetUserDraft(username string) (*model.Draft, bool) {
	draft := model.Draft{}
	if username == "" {
		return nil, true
	}
	if err := util.Db.Where(&model.Draft{UserName: username}).First(&draft).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &draft, false
}

// DeleteUserDraft 删除用户草稿
func DeleteUserDraft(username string) {
	draft, empty := GetUserDraft(username)
	if empty {
		fmt.Println("成功删除用户草稿")
		return
	}
	DeleteDraft(draft.DraftID)
	fmt.Println("成功删除用户草稿")
	return
}

// GetDraftPhoto 获取草稿的图片
func GetDraftPhoto(draftID string) *[]string {
	var photos []model.DraftPhoto
	util.Db.Where(&model.DraftPhoto{DraftID: draftID}).Find(&photos)
	var photoData []string
	for _, photo := range photos {
		photoData = append(photoData, photo.PhotoAddress)
	}
	return &photoData
}