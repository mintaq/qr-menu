// @/config/database.go
package database

import (
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

func MysqlGormConnection() error {
	var err error

	mysqlConnURL, err := utils.ConnectionURLBuilder("mysql")
	if err != nil {
		panic(err)
	}
	Database, err = gorm.Open(mysql.Open(mysqlConnURL), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(&models.Book{})

	return nil
}
