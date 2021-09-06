package bblog

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

//Log levels
const (
	DEBUG = iota
	INFO
	ERROR
	FATAL
)

//BANNER printed after initialized log.
// const BANNER = `
// [START] =================== User: %s =================== [START]
// `

//DefaultLog type
type DefaultLog struct {
	logger *log.Logger //logger will write serialized when simultaneously accessed by multiple goroutines
	level  int
	// user   string
}

//Local log object
var localLog = new(DefaultLog)

//Init initializes the localLog for use according to specified levels
func Init(filePath string) error {
	LogFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	LogLvl := "DEBUG"
	localLog.logger = log.New(LogFile, "["+"DEBUG"+"] ", log.Ldate|log.Ltime)
	switch LogLvl {
	case "DEBUG":
		localLog.level = DEBUG
	case "INFO":
		localLog.level = INFO
	case "ERROR":
		localLog.level = ERROR
	case "FATAL":
		localLog.level = FATAL
	}
	return nil
}

func SetLogLevel(logLevel int) {
	localLog.level = logLevel
}

//Write to log if at correct log lvl. Logging at ERROR level will be fatal as a call to os.Exit will follow.
func Write(logLevel int, data interface{}) {
	if localLog.logger == nil { //avoids being called if the localLog has not been initialized by Init
		return
	}
	msg := fmt.Sprintf("%v", data)
	if logLevel >= localLog.level {
		_, file, line, _ := runtime.Caller(1) //1 as argument signifies retrieving info of the function call that called this function
		switch logLevel {
		case DEBUG:
			localLog.logger.Printf("[DEBUG] [%s:%d] %s", file, line, msg)
		case INFO:
			localLog.logger.Printf("[INFO] [%s:%d] %s", file, line, msg)
		case ERROR:
			localLog.logger.Fatalf("[ERROR] [%s:%d] %s", file, line, msg)
		case FATAL:
			localLog.logger.Fatalf("[FATAL] [%s:%d] %s", file, line, msg)
		}
	}

}

//Debug logs an entry at debug level.
func Debug(data interface{}) {
	Write(DEBUG, data)
}

//Info logs an entry at info level.
func Info(data interface{}) {
	Write(INFO, data)
}

//Error logs an entry at fatal level.
func Error(data interface{}) {
	Write(ERROR, data)
}

//Fatal logs an entry at fatal level and follows with a call to os.Exit(1).
func Fatal(data interface{}) {
	Write(FATAL, data)
}
