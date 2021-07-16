package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddAnswer 加入Answer
func AddAnswer(answer *model.Answer) {
	util.Db.Create(&answer)
	fmt.Println("添加回答成功！")
}

// DeleteAnswer 用主键删除Answer
func DeleteAnswer(answerID string) {
	answer, empty := SearchAnswer(answerID)
	if empty {
		return
	}
	// 如果删除的是最佳答案，需要设置对应Question为未解决
	if answer.AnswerAcceptance {
		SetQuestionIsUnSolved(answer.QuestionID)
	}

	// 修改对应Question的answerNum
	MinusQuestionAnswerNum(answer.QuestionID)

	util.Db.Delete(&model.Answer{AnswerID: answerID})
	fmt.Println("删除回答成功！")
}

// SearchAnswer 用主键查找Answer
func SearchAnswer(answerID string) (*model.Answer, bool) {
	answer := model.Answer{}
	if err := util.Db.Where(&model.Answer{AnswerID: answerID}).First(&answer).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &answer, false
}

// AcceptAnswer 接受某个答案为最佳答案
func AcceptAnswer(answerID string) bool {
	answer, empty := SearchAnswer(answerID)
	if empty { return false }
	util.Db.Model(&answer).Updates(model.Answer{AnswerAcceptance: true})
	return true
}

// SetAnswererDeleted 设置该回答的回答者已经注销
func SetAnswererDeleted(answerID string) {
	answer, _ := SearchAnswer(answerID)
	util.Db.Model(&answer).Updates(model.Answer{AnswererName: DELETED_USER_NAME})
}

