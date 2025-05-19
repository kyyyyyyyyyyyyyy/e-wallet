package models

import (
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/golang-wallet?parseTime=true"))

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{})

	DB = db

}
