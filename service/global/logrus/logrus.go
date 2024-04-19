package logrus

import (
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func init() {

}

var (
	logFilePath = "./runtime/log" //File storagepath
	logFileName = "system.log"
)

func ReturnsInstance() *logrus.Logger {
	Logger := logrus.New()
	// log file
	fileName := path.Join(logFilePath, logFileName)
	// write file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open/write file", err)
	}
	//log level
Logger.SetLevel(logrus.DebugLevel)
//Set the output
Logger.Out = file
//Set rotate logs to achieve file splitting
logWriter, err := rotateLogs.New(
//Split file name
fileName+".%Y%m%d.log",
//Generate a soft link pointing to the latest log file
rotateLogs.WithLinkName(fileName),
//Set the maximum storage time (7 days)
rotateLogs.WithMaxAge(7*24*time.Hour), //Integer in hour
//Set the log cutting interval (1 day)
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	// Hook mechanism settings
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	//Add hooks to loggers
	Logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return Logger

}
