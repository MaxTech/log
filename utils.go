package log

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "reflect"
    "runtime"
    "strings"
    "sync"
    "time"
)

type Flag int

const (
    DEBUG Flag = 1
    INFO  Flag = 2
    WARN  Flag = 3
    ERROR Flag = 4
)

var flagTextMap = map[Flag]string{
    DEBUG: "Debug",
    INFO:  "Info",
    WARN:  "Warn",
    ERROR: "Error",
}

var flagCodeMap = map[Flag]int{
    DEBUG: 1,
    INFO:  2,
    WARN:  3,
    ERROR: 4,
}

func (f Flag) Text() string {
    return flagTextMap[f]
}

func (f Flag) Code() int {
    return flagCodeMap[f]
}

type AppLogger interface {
    Log(flag Flag, msg, logPosition string, v ...interface{})
    Debug(msg string, v ...interface{})
    HighQualityDebug(msg string, v ...interface{})
    Info(msg string, v ...interface{})
    HighQualityInfo(msg string, v ...interface{})
    Warn(msg string, v ...interface{})
    HighQualityWarn(msg string, v ...interface{})
    Error(msg string, v ...interface{})
    HighQualityError(msg string, v ...interface{})
}

type mxWriter interface {
    io.Writer
    SetFilePath(path string)
}

type logger struct {
    loggerName string
    filePath   string
    fileDate   string
    logLevel   int
    writer     mxWriter
    DEBUG      *log.Logger
    ERROR      *log.Logger
    INFO       *log.Logger
    WARN       *log.Logger
}

type mxLoggerWriter struct {
    mu       *sync.Mutex
    filePath string
}

func (mlw *mxLoggerWriter) SetFilePath(_filePath string) {
    mlw.filePath = _filePath
}

func (mlw *mxLoggerWriter) Write(_data []byte) (int, error) {
    mlw.mu.Lock()
    defer mlw.mu.Unlock()

    file, err := os.OpenFile(mlw.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
    if err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err)
        return 0, nil
    }
    defer file.Close()
    return file.Write(_data)
}

func newMxWriter() mxWriter {
    return &mxLoggerWriter{mu: new(sync.Mutex)}
}

func NewLogger(_loggerName string) AppLogger {
    logger := logger{
        loggerName: _loggerName,
        writer:     newMxWriter(),
    }
    logger.filePath, _ = filepath.Abs(fmt.Sprintf("./logs/%s", _loggerName))
    logger.updateFileDate()

    logger.DEBUG = log.New(nil, fmt.Sprintf("[%s]\t[%s]\t", _loggerName, DEBUG.Text()), DEBUG.Code())
    logger.ERROR = log.New(nil, fmt.Sprintf("[%s]\t[%s]\t", _loggerName, ERROR.Text()), ERROR.Code())
    logger.INFO = log.New(nil, fmt.Sprintf("[%s]\t[%s]\t", _loggerName, INFO.Text()), INFO.Code())
    logger.WARN = log.New(nil, fmt.Sprintf("[%s]\t[%s]\t", _loggerName, WARN.Text()), WARN.Code())

    logger.DEBUG.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.ERROR.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.INFO.SetFlags(log.LstdFlags | log.Lmicroseconds)
    logger.WARN.SetFlags(log.LstdFlags | log.Lmicroseconds)

    return &logger
}

func (l *logger) Log(_flag Flag, _msg, _logPosition string, _v ...interface{}) {
    logs := make([]string, 0)
    if _logPosition != "" {
        logs = append(logs, _logPosition)
    }
    logs = append(logs, _msg)
    
    if len(_v) == 1 {
        if values, ok := _v[0].([]string); ok {
            logs = append(logs, values...)
        }
    } else {
        for _, val := range _v {
            logs = append(logs, fmt.Sprintf("%v: %v", reflect.TypeOf(val), val))
        }
    }
    
    logStr := strings.Join(logs, "\t")

    l.updateFileDate()
    flagFilePath := l.getFlagFilePath(_flag)
    l.writer.SetFilePath(flagFilePath)

    switch _flag {
    case DEBUG:
        if l.logLevel <= _flag.Code() {
            l.DEBUG.SetOutput(l.writer)
            _ = l.DEBUG.Output(1, fmt.Sprintf("\t%s", logStr))
        }
    case ERROR:
        if l.logLevel <= _flag.Code() {
            l.ERROR.SetOutput(l.writer)
            _ = l.ERROR.Output(1, fmt.Sprintf("\t%s", logStr))
        }
    case INFO:
        if l.logLevel <= _flag.Code() {
            l.INFO.SetOutput(l.writer)
            _ = l.INFO.Output(1, fmt.Sprintf("\t%s", logStr))
        }
    case WARN:
        if l.logLevel <= _flag.Code() {
            l.WARN.SetOutput(l.writer)
            _ = l.WARN.Output(1, fmt.Sprintf("\t%s", logStr))
        }
    default:
        tempLogger := log.New(os.Stdout, fmt.Sprintf("[%s]\t[%s]\t", l.loggerName, "UNKNOWN"), ERROR.Code())
        tempLogger.SetFlags(log.LstdFlags | log.Lmicroseconds)
        tempLogger.SetOutput(l.writer)
        _ = tempLogger.Output(1, fmt.Sprintf("\t%s", logStr))
    }
}

func (l *logger) Debug(_msg string, _v ...interface{}) {
    l.Log(DEBUG, _msg, "", _v...)
}

//HighQuality
func (l *logger) HighQualityDebug(_msg string, _v ...interface{}) {
    funcPtr, file, line, ok := runtime.Caller(1)
    if ok {
        funcName := runtime.FuncForPC(funcPtr).Name()
        logPosition := fmt.Sprintf("[funcName: %s, file: %s:%d]", funcName, file, line)
        l.Log(DEBUG, _msg, logPosition, _v...)
        return
    }
    l.Log(DEBUG, _msg, "", _v...)
}

func (l *logger) Info(_msg string, _v ...interface{}) {
    l.Log(INFO, _msg, "", _v...)
}

//HighQuality
func (l *logger) HighQualityInfo(_msg string, _v ...interface{}) {
    funcPtr, file, line, ok := runtime.Caller(1)
    if ok {
        funcName := runtime.FuncForPC(funcPtr).Name()
        logPosition := fmt.Sprintf("[funcName: %s, file: %s:%d]", funcName, file, line)
        l.Log(INFO, _msg, logPosition, _v...)
        return
    }
    l.Log(INFO, _msg, "", _v...)
}

func (l *logger) Warn(_msg string, _v ...interface{}) {
    l.Log(WARN, _msg, "", _v...)
}

//HighQuality
func (l *logger) HighQualityWarn(_msg string, _v ...interface{}) {
    funcPtr, file, line, ok := runtime.Caller(1)
    if ok {
        funcName := runtime.FuncForPC(funcPtr).Name()
        logPosition := fmt.Sprintf("[funcName: %s, file: %s:%d]", funcName, file, line)
        l.Log(WARN, _msg, logPosition, _v...)
        return
    }
    l.Log(WARN, _msg, "", _v...)
}

func (l *logger) Error(_msg string, _v ...interface{}) {
    l.Log(ERROR, _msg, "", _v...)
}

//HighQuality
func (l *logger) HighQualityError(_msg string, _v ...interface{}) {
    funcPtr, file, line, ok := runtime.Caller(1)
    if ok {
        funcName := runtime.FuncForPC(funcPtr).Name()
        logPosition := fmt.Sprintf("[funcName: %s, file: %s:%d]", funcName, file, line)
        l.Log(ERROR, _msg, logPosition, _v...)
        return
    }
    l.Log(ERROR, _msg, "", _v...)
}

func (l *logger) updateFileDate() {
    dateString := time.Now().Format("20060102")
    if l.fileDate != dateString {
        l.fileDate = dateString
    }
}

func (l *logger) getFlagFilePath(_flag Flag) string {
    flagPath := fmt.Sprintf("%s/%s", l.filePath, strings.ToLower(_flag.Text()))
    _, err := os.Stat(flagPath)
    if err != nil && os.IsNotExist(err) {
        _ = os.MkdirAll(flagPath, os.ModePerm)
    }
    flagFilePath := fmt.Sprintf("%s/%s_%s_%s.log", flagPath, strings.ToLower(l.loggerName), strings.ToLower(_flag.Text()), l.fileDate)
    return flagFilePath
}
