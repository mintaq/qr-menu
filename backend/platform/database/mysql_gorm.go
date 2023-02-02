// @/config/database.go
package database

import (
	"log"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database instance
type DbInstance struct {
	Db *gorm.DB
}

var Database *gorm.DB

func MysqlGormConnection() {
	var err error

	mysqlConnURL, err := utils.ConnectionURLBuilder("mysql")
	if err != nil {
		log.Panic(err.Error())
	}
	Database, err = gorm.Open(mysql.Open(mysqlConnURL+"?parseTime=true"), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Panic(err.Error())
	}

	if err := Database.AutoMigrate(&models.User{}); err != nil {
		log.Panic(err.Error())
	}
}
