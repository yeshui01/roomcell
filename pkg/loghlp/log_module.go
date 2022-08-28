/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-01-12 18:28:05
 * @LastEditTime: 2022-01-17 13:32:47
 * @FilePath: roomcell\pkg\loghlp\log_module.go
 */

package loghlp

import (
	"os"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
)

type LogHelper struct {
	consoleLog     *logrus.Logger
	openConsoleLog bool
	fileLog        *logrus.Logger
	userLog        *logrus.Logger
}

var loginst *LogHelper
var mu sync.Mutex

func instance() *LogHelper {
	if loginst == nil {
		mu.Lock()
		defer mu.Unlock()
		if loginst == nil {
			loginst = &LogHelper{
				consoleLog:     nil,
				fileLog:        nil,
				userLog:        nil,
				openConsoleLog: false,
			}
		}
	}
	return loginst
}

type colorHook struct {
}

func (c *colorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel}
}

func (c *colorHook) Fire(entry *logrus.Entry) error {
	// fmt.Printf("\033[1;34;40m%s\033[0m\n", "这是一个error级别的信息")
	// fmt.Printf("\033[0;33;40m%s : %s\033[0m\n", "这是一个error级别的信息", entry.Message)
	return nil
}

/**
 * @description: 激活文件日志
 * @param {string} filePath: 日志文件路径
 * @param {string} servName: 服务名字标识
 * @return {*}
 */
func ActiveFileLog(filePath string, servName string) {
	if instance().fileLog == nil {
		instance().fileLog = logrus.New()
		// instance().fileLog.ReportCaller = reportCaller
		writer, _ := rotatelogs.New(
			//filePath+"servName_%Y%m%d_%H%M.log",
			filePath+"//"+servName+"_%Y%m%d.log",
			// rotatelogs.WithLinkName(filePath),
			rotatelogs.WithMaxAge(time.Duration(3600*24*7)*time.Second),   // 保留最近一周
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second), // 每天回滚
			rotatelogs.WithRotationSize(1024*1024*60),                     // 一个日志文件大小设为60M
			rotatelogs.ForceNewFile(),
		)
		instance().fileLog.SetOutput(writer)
		instance().fileLog.SetLevel(logrus.DebugLevel)
	}
}
func ActiveFileLogReportCaller(reportCaller bool) {
	if instance().fileLog != nil {
		instance().fileLog.ReportCaller = reportCaller
	}
}

/**
 * @description: 激活用户数据日志
 * @param {string} filePath: 日志文件路径
 * @param {string} servName: 服务名字标识
 * @return {*}
 */
func ActiveUserLog(filePath string, servName string) {
	instance().userLog = logrus.New()
	writer, _ := rotatelogs.New(
		filePath+"userlog_servName%Y%m%d_%H%M.log",
		// rotatelogs.WithLinkName(filePath),
		rotatelogs.WithMaxAge(time.Duration(3600*24*7)*time.Second), // 保留最近一周
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		rotatelogs.ForceNewFile(),
	)
	instance().userLog.SetOutput(writer)
	instance().userLog.SetLevel(logrus.DebugLevel)
}

/**
 * @description: 激活控制台日志
 * @param {string} filePath:文件路径
 * @return {*}
 */
func ActiveConsoleLog() {
	if instance().consoleLog == nil {
		instance().openConsoleLog = true
		instance().consoleLog = logrus.New()
		instance().consoleLog.SetLevel(logrus.DebugLevel)
		instance().consoleLog.SetOutput(os.Stdout)
		// instance().consoleLog.SetFormatter(&logrus.JSONFormatter{})
		// instance().consoleLog.AddHook(&colorHook{})
		instance().consoleLog.SetFormatter(&logrus.TextFormatter{
			ForceColors:     false,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

/**
 * @description: 设置控制台日志等级
 * @param {logrus.Level} logLv:日志等级,logrus.DebugLevel|logrus.InfoLevel,logrus.WarnLevel|logrus.ErrorLevel
 * @return {*}
 */
func SetConsoleLogLevel(logLv logrus.Level) {
	if instance().consoleLog != nil {
		instance().consoleLog.SetLevel(logLv)
	}
}

/**
 * @description: 设置文件日志等级
 * @param {logrus.Level} logLv:日志等级,logrus.DebugLevel|logrus.InfoLevel,logrus.WarnLevel|logrus.ErrorLevel
 * @return {*}
 */
func SetFileLogLevel(logLv logrus.Level) {
	if instance().fileLog != nil {
		instance().fileLog.SetLevel(logLv)
	}
}

/**
 * @description: 设置用户日志等级
 * @param {logrus.Level} logLv:日志等级,logrus.DebugLevel|logrus.InfoLevel,logrus.WarnLevel|logrus.ErrorLevel
 * @return {*}
 */
func SetUserLogLevel(logLv logrus.Level) {
	if instance().userLog != nil {
		instance().userLog.SetLevel(logLv)
	}
}

// debug
func Debug(args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Debug(args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Debug(args...)
	}
}
func Debugf(format string, args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Debugf(format, args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Debugf(format, args...)
	}
}

// Info
func Info(args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Info(args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Info(args...)
	}
}
func Infof(format string, args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Infof(format, args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Infof(format, args...)
	}
}

// Warn
func Warn(args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Warn(args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Warn(args...)
	}
}
func Warnf(format string, args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Warnf(format, args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Warnf(format, args...)
	}
}

// Error
func Error(args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Error(args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Error(args...)
	}
}
func Errorf(format string, args ...interface{}) {
	if instance().consoleLog != nil && instance().openConsoleLog {
		instance().consoleLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Errorf(format, args...)
	}
	if instance().fileLog != nil {
		instance().fileLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Errorf(format, args...)
	}
}

func UserDebugf(format string, args ...interface{}) {
	if instance().userLog != nil {
		instance().userLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Debugf(format, args...)
	}
}
func UserInfof(format string, args ...interface{}) {
	if instance().userLog != nil {
		instance().userLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Infof(format, args...)
	}
}
func UserWarnf(format string, args ...interface{}) {
	if instance().userLog != nil {
		instance().userLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Warnf(format, args...)
	}
}
func UserErrorf(format string, args ...interface{}) {
	if instance().userLog != nil {
		instance().userLog.WithField("msec", time.Now().UnixNano()/1e6%1000).Errorf(format, args...)
	}
}
