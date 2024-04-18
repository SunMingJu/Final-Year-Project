package mysql

import (
	"fmt"
	"simple-video-net/global/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ReturnsInstance() *gorm.DB {
	var err error

	var mysqlConfig = config.Config.SqlConfig
	// Create a link
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", mysqlConfig.User, mysqlConfig.Password, mysqlConfig.IP, mysqlConfig.Port, mysqlConfig.Database)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), Print all sql
	})
	if err != nil {
		fmt.Printf("Database link error - %v \n", err)
	}
	if Db.Error != nil {
		fmt.Printf("Database error - %v \n", Db.Error)
	}
	return Db
}
