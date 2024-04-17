package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func LogMiddleWare() gin.HandlerFunc {
	var (
		logFilePath = "./runtime/log" //File storage path
		logFileName = "system.log"
	)
	// log file
	fileName := path.Join(logFilePath, logFileName)
	// write to a file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failure to open/write file", err)
		return nil
	}
	// instantiated
	logger := logrus.New()
	// Log level
	logger.SetLevel(logrus.DebugLevel)
	// Setting the output
	logger.Out = file
	// Set rotate logs to split files.
	logWriter, err := rotateLogs.New(
		// Split file name
		fileName+".%Y%m%d.log",
		// Generate a soft link to the latest log file
		rotateLogs.WithLinkName(fileName),
		// Setting the maximum retention time (7 days)
		rotateLogs.WithMaxAge(7*24*time.Hour), //Integer in hours
		// Set log cutting interval (1 day)
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	// Setting up the hook mechanism
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	//Adding a hook to loggers
	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	return func(c *gin.Context) {
		c.Next()
		//request method
		method := c.Request.Method
		//Request Routing
		reqUrl := c.Request.RequestURI
		//status code
		statusCode := c.Writer.Status()
		//Request ip
		clientIP := c.ClientIP()
		// Print Log
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"client_ip":   clientIP,
			"req_method":  method,
			"req_uri":     reqUrl,
		}).Info()
	}
}
