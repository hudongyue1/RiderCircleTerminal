package model

import (
	"time"
)

// BlogType BlogType `gorm:"ForeignKey:BlogTypeID;AssociationForeignKey:ID"`

// Circle -----------------------------------------------
type Circle struct {
	CircleName string `gorm:"primary_key; not null"`
	CircleMasterName string `gorm:"not null"`
	UserNum int	`gorm:"default:0; not null"`
	CircleDescription string `gorm:"default:'暂无描述！'; not null"`

	CreatedAt    time.Time
}

type CirclePhoto struct {
	CirclePhotoID string `gorm:"primary_key"`

	Circle Circle `gorm:"ForeignKey:CircleName;AssociationForeignKey:CircleName"`
	CircleName string `gorm:"not null"`

	PhotoAddress string `gorm:"not null"`
}

// User -----------------------------------------------
type User struct {
	UserName string	`gorm:"primary_key; not null"`
	Password string `gorm:"not null"`
	UserPhoto string	`gorm:"default:'./public/photo/default/userPhoto.png'; not null"`
	UserDescription string `gorm:"default:'这家伙很懒，没有个人描述'; not null"`

	ActivatedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRelation struct {
	UserRelationID string `gorm:"primary_key"`

	Circle Circle `gorm:"ForeignKey:CircleName;AssociationForeignKey:CircleName"`
	CircleName string `gorm:"not null"`

	User User `gorm:"ForeignKey:UserName;AssociationForeignKey:UserName"`
	UserName string `gorm:"not null"`
}

// Post -----------------------------------------------
type Post struct {
	PostID string `gorm:"primary_key"`

	User User `gorm:"ForeignKey:PostIssuerName;AssociationForeignKey:UserName"`
	PostIssuerName string `gorm:"not null"`

	PostIssueTime time.Time `gorm:"not null"`
	PostDescription string `gorm:"not null"`

	Circle Circle `gorm:"ForeignKey:PostCircleName;AssociationForeignKey:CircleName"`
	PostCircleName string `gorm:"not null"`

	PostUpNum int `gorm:"default:0; not null"`
	PostCommentaryNum int `gorm:"default:0; not null"`
}

type PostPhoto struct {
	PostPhotoID string	`gorm:"primary_key"`

	Post Post `gorm:"ForeignKey:PostID;AssociationForeignKey:PostID"`
	PostID string `gorm:"not null"`

	PhotoAddress string `gorm:"not null"`
}

type Commentary struct {
	CommentaryID string	`gorm:"primary_key"`

	Post Post `gorm:"ForeignKey:PostID;AssociationForeignKey:PostID"`
	PostID string `gorm:"not null"`

	User User `gorm:"ForeignKey:CommenterName;AssociationForeignKey:UserName"`
	CommenterName string `gorm:"not null"`

	CommentaryDescription string `gorm:"not null"`
	CommentaryTime time.Time `gorm:"not null"`
	ReplyNum int `gorm:"default:0; not null"`
}

type Reply struct {
	ReplyID string	`gorm:"primary_key"`

	Commentary Commentary `gorm:"ForeignKey:CommentaryID;AssociationForeignKey:CommentaryID"`
	CommentaryID string `gorm:"not null"`

	User User `gorm:"ForeignKey:ReplyerName;AssociationForeignKey:UserName"`
	ReplyerName string `gorm:"not null"`

	ReplyDescription string `gorm:"not null"`
	ReplyTime time.Time `gorm:"not null"`
}

// PostUpRelation -----------------------------------------------
type PostUpRelation struct {
	PostUpRelationID string	`gorm:"primary_key; not null"`

	Post Post `gorm:"ForeignKey:PostID;AssociationForeignKey:PostID"`
	PostID string `gorm:"not null"`

	User User `gorm:"ForeignKey:UserName;AssociationForeignKey:UserName"`
	UserName string `gorm:"not null"`
}

// Question -----------------------------------------------
type Question struct {
	QuestionID string `gorm:"primary_key"`
	QuestionIssuerName string `gorm:"not null"`
	QuestionIssueTime time.Time `gorm:"not null"`
	QuestionDescription string `gorm:"not null"`
	QuestionCircleName string `gorm:"not null"`
	QuestionSolved bool `gorm:"default:false; not null"`
	QuestionAnswerNum int `gorm:"default:0; not null"`
}

type QuestionPhoto struct {
	QuestionPhotoID string	`gorm:"primary_key"`

	Question Question `gorm:"ForeignKey:QuestionID;AssociationForeignKey:QuestionID"`
	QuestionID string `gorm:"not null"`

	PhotoAddress string `gorm:"not null"`
}

type Answer struct {
	AnswerID string	`gorm:"primary_key"`

	Question Question `gorm:"ForeignKey:QuestionID;AssociationForeignKey:QuestionID"`
	QuestionID string `gorm:"not null"`

	User User `gorm:"ForeignKey:AnswererName;AssociationForeignKey:UserName"`
	AnswererName string `gorm:"not null"`

	AnswerTime time.Time `gorm:"not null"`
	AnswerDescription string `gorm:"not null"`
	AnswerAcceptance bool `gorm:"default:false; not null"`
}

// Draft -----------------------------------------------
type Draft struct {
	DraftID string `gorm:"primary_key"`

	User User `gorm:"ForeignKey:UserName;AssociationForeignKey:UserName"`
	UserName string `gorm:"not null"`

	Choose int `gorm:"default:0; not null"`
	Description string
	CircleName string
}

type DraftPhoto struct {
	DraftPhotoID string `gorm:"primary_key"`

	Draft   Draft  `gorm:"ForeignKey:DraftID;AssociationForeignKey:DraftID"`
	DraftID string

	PhotoAddress string
}

// Administrator -----------------------------------------------
type Administrator struct {
	AdminName string `gorm:"primary_key; not null"`
	Password string `gorm:"not null"`
}