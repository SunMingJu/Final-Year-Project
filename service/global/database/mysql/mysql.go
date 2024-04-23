package mysql

import (
	"simple-video-net/global/config"
	globalLog "simple-video-net/global/logrus"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type MyWriter struct {
	log *logrus.Logger
}

// Printf Implement gorm/logger.writer interface
func (m *MyWriter) Printf(format string, v ...interface{}) {
	m.log.Errorf(fmt.Sprintf(format, v...))
}

func NewMyWriter() *MyWriter {
	log := globalLog.ReturnsInstance()
	return &MyWriter{log: log}
}

func ReturnsInstance() *gorm.DB {
	var err error
	var mysqlConfig = config.Config.SqlConfig
	//sql logging
	myLogger := logger.New(
		//Set up logger
		NewMyWriter(),
		logger.Config{
			LogLevel: logger.Error,
		},
	)
	// create link
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", mysqlConfig.User, mysqlConfig.Password, mysqlConfig.IP, mysqlConfig.Port, mysqlConfig.Database)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger,
	})
	if err != nil {
		log.Fatalf("database connect err- %v \n", err)
	}
	if Db.Error != nil {
		log.Fatalf("database err- %v \n", Db.Error)
	}
	return Db
}