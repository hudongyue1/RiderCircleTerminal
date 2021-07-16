package util

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func InitDB() {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "test01"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	Db, err = gorm.Open(driverName, args)

	// 建表
	//Db.AutoMigrate(&model.Circle{})
	//Db.AutoMigrate(&model.CirclePhoto{})
	//Db.AutoMigrate(&model.User{})
	//Db.AutoMigrate(&model.UserRelation{})
	//Db.AutoMigrate(&model.Post{})
	//Db.AutoMigrate(&model.PostPhoto{})
	//Db.AutoMigrate(&model.Commentary{})
	//Db.AutoMigrate(&model.Reply{})
	//Db.AutoMigrate(&model.PostUpRelation{})
	//Db.AutoMigrate(&model.Question{})
	//Db.AutoMigrate(&model.QuestionPhoto{})
	//Db.AutoMigrate(&model.Answer{})
	//Db.AutoMigrate(&model.Draft{})
	//Db.AutoMigrate(&model.DraftPhoto{})
	//Db.AutoMigrate(&model.Administrator{})


	if err != nil {
		panic("falied to connect database, err:" + err.Error())
	}
}



//func init() {
//	db, err = gorm.Open("mysql", "root:Hdy15608156313@(localhost)/test01?charset=utf8mb4&parseTime=True&loc=Local")
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//}
