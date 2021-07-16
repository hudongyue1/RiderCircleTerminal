package dao

import (
	"RiderCircleTerminal/model"
	"RiderCircleTerminal/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

// AddQuestion AddPost 加入Post
func AddQuestion(question *model.Question, photos *[]string) {
	AddQuestionPhoto(question.QuestionID, photos)
	util.Db.Create(&question)
	fmt.Println("增加question成功！")
}

// AddQuestionPhoto 给问答加上QuestionPhoto
func AddQuestionPhoto(questionID string, photos *[]string) {
	for _, address := range *photos {
		if address != "" {
			postPhoto := model.QuestionPhoto{QuestionPhotoID: util.CreateUUID("questionPhoto"),
				QuestionID: questionID, PhotoAddress: address}
			util.Db.Create(&postPhoto)
		}
	}
	fmt.Println("增加questionPhoto成功！")
}

// DeleteQuestion 用主键删除Question
func DeleteQuestion(questionID string) {
	// 删除questionPhoto
	util.Db.Where("question_id = ?", questionID).Delete(model.QuestionPhoto{})

	// 首先删除Question的Answer
	var answers []model.Answer
	util.Db.Where(&model.Answer{QuestionID: questionID}).Find(&answers)
	for _, answer := range answers {
		util.Db.Delete(&model.Answer{AnswerID: answer.AnswerID})
	}

	util.Db.Delete(&model.Question{QuestionID: questionID})
	fmt.Println("删除问答成功！")
}

// SearchQuestion 用主键查找Question
func SearchQuestion(questionID string) (*model.Question, bool) {
	//err := util.Db.Where(&model.User{UserName: username}).First(&user).Error; gorm.IsRecordNotFoundError(err)
	question := model.Question{}
	if questionID == "" {
		return nil, true
	}
	if err := util.Db.Where(&model.Question{QuestionID: questionID}).First(&question).Error; gorm.IsRecordNotFoundError(err) {
		return nil, true
	}
	return &question, false
}

// UpdateQuestion 更新Question
func UpdateQuestion(question *model.Question) {
	questionTemp, _ := SearchQuestion(question.QuestionID)
	util.Db.Model(&questionTemp).Updates(question)
	fmt.Println("更新帖子信息！")
}

// GetAcceptAnswer 获取最佳答案a
func GetAcceptAnswer(question *model.Question) *model.Answer {
	answer := model.Answer{}
	util.Db.Where("answer_acceptance = true and question_id = ?", question.QuestionID).Find(&answer)
	return &answer
}

// GetSomeNewQuestion 获取一些最新的Question
func GetSomeNewQuestion(num int) *[]model.Question {
	var questions []model.Question
	util.Db.Model(&model.Question{}).Order("question_issue_time").Limit(num).Find(&questions)
	return &questions
}

// GetSomeHotQuestion 获取一些最热的Question
func GetSomeHotQuestion(num int) *[]model.Question {
	var questions []model.Question
	util.Db.Model(&model.Question{}).Order("question_answer_num desc").Limit(num).Find(&questions)
	return &questions
}

// GetSomeNewQuestionInCircle 获取一些车圈的最新的Question
func GetSomeNewQuestionInCircle(num int, circleName string) *[]model.Question {
	var questions []model.Question
	util.Db.Where(&model.Question{QuestionCircleName: circleName}).Order("question_issue_time desc").Limit(num).Find(&questions)
	return &questions
}

// GetSomeHotQuestionInCircle 获取一些车圈的最热的Question
func GetSomeHotQuestionInCircle(num int, circleName string) *[]model.Question {
	var questions []model.Question
	util.Db.Where(&model.Question{QuestionCircleName: circleName}).Order("question_answer_num desc").Limit(num).Find(&questions)
	return &questions
}

// GetSomeNewQuestionForUser 获取一些给用户的最新的Question
func GetSomeNewQuestionForUser(num int, userName string) *[]model.Question {
	var questions []model.Question
	util.Db.Model(&model.Question{}).Order("question_issue_time").Limit(num).Find(&questions)
	// util.Db.Model(&model.Question{QuestionCircleName: circleName}).Order("question_issue_time").Limit(num).Find(&questions)
	return &questions
}

// GetSomeHotQuestionForUser 获取一些给用户的最热的Question
func GetSomeHotQuestionForUser(num int, userName string) *[]model.Question {
	var questions []model.Question
	circles := GetUserAllCircle(userName)
	var circlesData []string
	for _, circle := range *circles {
		circlesData = append(circlesData, circle.CircleName)
	}

	util.Db.Where("question_circle_name IN (?)", circlesData).Order("question_answer_num desc").Limit(num).Find(&questions)

	// 如果用户加入的车圈过少，就推荐一些其他车圈的最热问答
	if len(questions) <  num {
		temp := GetSomeHotQuestion(num-len(questions))
		for _, question := range *temp {
			questions = append(questions, question)
		}
	}

	return &questions
}

// GetQuestionAnswer 获取问答的回答
func GetQuestionAnswer(questionID string) *[]model.Answer {
	var answers []model.Answer
	util.Db.Where(&model.Answer{QuestionID: questionID}).Find(&answers)
	return &answers
}

// GetQuestionPhoto 获取问答的图片
func GetQuestionPhoto(questionID string) *[]string {
	var photos []model.QuestionPhoto
	util.Db.Where(&model.QuestionPhoto{QuestionID: questionID}).Find(&photos)

	var photoData []string
	for _, photo := range photos {
		photoData = append(photoData, photo.PhotoAddress)
	}
	return &photoData
}

// AddQuestionAnswerNum 问答回答数+1
func AddQuestionAnswerNum(questionID string) bool {
	question, empty := SearchQuestion(questionID)
	if empty { return false }
	util.Db.Model(&question).Updates(model.Question{QuestionAnswerNum: question.QuestionAnswerNum+1})
	return true
}

// MinusQuestionAnswerNum 问答回答数-1
func MinusQuestionAnswerNum(questionID string) bool {
	question, empty := SearchQuestion(questionID)
	if empty { return false }
	util.Db.Model(&question).Update("question_answer_num", question.QuestionAnswerNum-1)
	return true
}

// SetQuestionIsSolved 设置问答已解决
func SetQuestionIsSolved(answerID string, questionID string) bool {
	if !AcceptAnswer(answerID) { return false }
	question, empty := SearchQuestion(questionID)
	if empty { return false }
	util.Db.Model(&question).Updates(model.Question{QuestionSolved: true})
	return true
}

// SetQuestionIsUnSolved 设置问答未解决
func SetQuestionIsUnSolved(questionID string) {
	question, _ := SearchQuestion(questionID)
	util.Db.Model(&question).Update("question_solved", false)
}

// GetAllQuestionInDb 获取数据库中所有的问答
func GetAllQuestionInDb() *[]model.Question {
	var questions []model.Question
	util.Db.Find(&questions)
	return &questions
}

