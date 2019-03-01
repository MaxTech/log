package max_log

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    DEBUG = "Debug"
    INFO  = "Info"
    WARN  = "Warn"
    ERROR = "Error"
)

var logFlagMap = map[LogFlag]int{
    DEBUG: 1,
    INFO:  2,
    WARN:  3,
    ERROR: 4,
}

type LogFlag string

func (lf LogFlag) Code() int {
    return logFlagMap[lf]
}

type AppLogger interface {
    Log(flag string, v ...interface{})
    Debug(v ...interface{})
    Info(v ...interface{})
    Warn(v ...interface{})
    Error(v ...interface{})
}

type logger struct {
    loggerName string
    filePath   string
    fileDate   string
    logLevel   int
    DEBUG      log.Logger
    ERROR      log.Logger
    INFO       log.Logger
    WARN       log.Logger
}

func NewLogger(loggerName, basePath string) AppLogger {
    logger := logger{
        loggerName: loggerName,
    }
    logger.filePath, _ = filepath.Abs(fmt.Sprintf("%s/%s", basePath, loggerName))
    logger.updateFileDate()

    logger.DEBUG = log.Logger{}
    logger.ERROR = log.Logger{}
    logger.INFO = log.Logger{}
    logger.WARN = log.Logger{}

    logger.DEBUG.SetPrefix(fmt.Sprintf("[%s]\t[DEBUG]\t", loggerName))
    logger.ERROR.SetPrefix(fmt.Sprintf("[%s]\t[ERROR]\t", loggerName))
    logger.INFO.SetPrefix(fmt.Sprintf("[%s]\t[INFO]\t", loggerName))
    logger.WARN.SetPrefix(fmt.Sprintf("[%s]\t[WARN]\t", loggerName))

    logger.DEBUG.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.ERROR.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.INFO.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.WARN.SetFlags(log.LstdFlags | log.Lmicroseconds)

    return &logger
}

func (l *logger) Log(flag string, v ...interface{}) {
    logs := make([]string, 0)
    for _, val := range v {
        logs = append(logs, fmt.Sprint(val))
    }
    logStr := strings.Join(logs, "\t")

    l.updateFileDate()
    flagFilePath := l.getFlagFilePath(flag)
    flagFile, _ := os.OpenFile(flagFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
    defer flagFile.Close()

    flagEnum := LogFlag(flag)

    switch flag {
    case DEBUG:
        if l.logLevel <= flagEnum.Code() {
            l.doOutPut(l.DEBUG, flagFile, fmt.Sprintf("\t%s", logStr))
        }
    case ERROR:
        if l.logLevel <= flagEnum.Code() {
            l.doOutPut(l.ERROR, flagFile, fmt.Sprintf("\t%s", logStr))
        }
    case INFO:
        if l.logLevel <= flagEnum.Code() {
            l.doOutPut(l.INFO, flagFile, fmt.Sprintf("\t%s", logStr))
        }
    case WARN:
        if l.logLevel <= flagEnum.Code() {
            l.doOutPut(l.WARN, flagFile, fmt.Sprintf("\t%s", logStr))
        }
    default:
        tempLogger := log.Logger{}
        tempLogger.SetPrefix(fmt.Sprintf("[%s]\t%s\t", l.loggerName, flag))
        tempLogger.SetFlags(log.LstdFlags | log.Lmicroseconds)
        l.doOutPut(tempLogger, flagFile, fmt.Sprintf("\t%s", logStr))
    }
}

func (l *logger) Debug(v ...interface{}) {
    l.Log(DEBUG, v...)
}

func (l *logger) Info(v ...interface{}) {
    l.Log(INFO, v...)
}

func (l *logger) Warn(v ...interface{}) {
    l.Log(WARN, v...)
}

func (l *logger) Error(v ...interface{}) {
    l.Log(ERROR, v...)
}

func (l *logger) updateFileDate() {
    dateString := time.Now().Format("20060102")
    if l.fileDate != dateString {
        l.fileDate = dateString
    }
}

func checkDir(dirPath string) {
    _, err := os.Stat(dirPath)
    if err != nil && os.IsNotExist(err) {
        _ = os.MkdirAll(dirPath, os.ModePerm)
    }
}

func (l *logger) getFlagFilePath(flag string) string {
    flagPath := fmt.Sprintf("%s/%s", l.filePath, strings.ToLower(flag))
    checkDir(flagPath)
    flagFilePath := fmt.Sprintf("%s/%s.log", flagPath, l.fileDate)
    return flagFilePath
}

func (l *logger) doOutPut(logger log.Logger, writer io.Writer, logString string)  {
    logger.SetOutput(writer)
    _ = logger.Output(2, fmt.Sprintf("\t%s", logString))
}
