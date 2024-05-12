package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnDatabase() {
	DB, err = gorm.Open(sqlite.Open("jkbd.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = DB.DB()

	fmt.Println("数据库连接成功！")
}
