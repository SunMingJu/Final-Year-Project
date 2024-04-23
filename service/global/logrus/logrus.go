package logrus

import (
	"encoding/json"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type JsonInfo struct {
	Time     string `json:"time"`
	Level    string `json:"level"`
	Msg      string `json:"msg"`
	File     string `json:"file,omitempty"`
	Function string `json:"function,omitempty"`
}

//JsonFormatter Custom json parsing
type JsonFormatter struct {
	logrus.JSONFormatter
}

func (f *JsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Construct JSON data
	info := &JsonInfo{
		Time:  entry.Time.Format("2006.01.02 15:04:05"),
		Level: entry.Level.String(),
		Msg:   entry.Message,
	}
	//Only level matching will print caller information
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.WarnLevel || entry.Level == logrus.DebugLevel || entry.Level == logrus.PanicLevel {
		info.File = entry.Caller.File
		info.Function = entry.Caller.Function
	}
	lineBreak := "\n"
	jsonData, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	formattedMsg := string(jsonData) + lineBreak
	return []byte(formattedMsg), nil
}

var (
	logFilePath = "./runtime/log" //File storage path
)

func ReturnsInstance() *logrus.Logger {
	Logger := logrus.New()
	// Log level
	Logger.SetLevel(logrus.DebugLevel)
	//Print caller information
	Logger.SetReportCaller(true)
	//define to empty output
	Logger.SetOutput(ioutil.Discard)

	// Set rotate logs to achieve file splitting
	logInfoWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/info.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	logFataWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/fata.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	logDebugWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/debug.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	logWarnWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/warn.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	logErrorWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/error.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	logPanicWriter, _ := rotateLogs.New(
		logFilePath+"/%Y-%m-%d/panic.log",
		rotateLogs.WithMaxAge(7*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)

	// Hook mechanism settings
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logInfoWriter,
		logrus.FatalLevel: logFataWriter,
		logrus.DebugLevel: logDebugWriter,
		logrus.WarnLevel:  logWarnWriter,
		logrus.ErrorLevel: logErrorWriter,
		logrus.PanicLevel: logPanicWriter,
	}
	Logger.Formatter = &JsonFormatter{}
	//Add hooks to loggers
	Logger.AddHook(lfshook.NewHook(writerMap, &JsonFormatter{}))

	return Logger
}