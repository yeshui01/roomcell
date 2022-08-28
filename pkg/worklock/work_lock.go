package worklock

import (
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

var lockReport bool = true // 默认开启
var debugMode bool = false // 调试模式
// 设置是否开启加锁报告
func SetLockReport(r bool) {
	lockReport = r
}

// 获取枷锁调用
func getLockCaller() (string, bool) {
	functime, _, _, ok := runtime.Caller(2)
	if ok {
		funName := runtime.FuncForPC(functime).Name()
		return funName, true
	}
	return "", false
}

// 互斥锁
type WorkLock struct {
	mu         sync.Mutex
	moduleName string
}

func NewWorkLock(wkName string) *WorkLock {
	return &WorkLock{
		moduleName: wkName,
	}
}

func (l *WorkLock) Lock() {
	l.mu.Lock()
	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkLogicLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkLogicLock")
		}
	}
}
func (l *WorkLock) UnLock() {
	l.mu.Unlock()
	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkLogicUnLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkLogicUnLock")
		}
	}
}

// 读写锁
type WorkRWLock struct {
	mu         sync.RWMutex
	moduleName string
}

func NewWorkRWLock(wkModule string) *WorkRWLock {
	return &WorkRWLock{
		moduleName: wkModule,
	}
}
func (l *WorkRWLock) Lock() {
	l.mu.Lock()
	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkRWLogicLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkRWLogicLock")
		}
	}
}
func (l *WorkRWLock) Unlock() {
	l.mu.Unlock()
	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkRWLogicUnLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkRWLogicUnLock")
		}
	}
}
func (l *WorkRWLock) RLock() {
	if debugMode {
		l.mu.Lock() // 测试,故意加互斥锁,便于排除死锁
	} else {
		l.mu.RLock()
	}

	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkRWLogicRLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkRWLogicRLock")
		}
	}
}
func (l *WorkRWLock) RUnlock() {
	if debugMode {
		l.mu.Unlock() // 测试,故意加互斥锁,便于排除死锁
	} else {
		l.mu.RUnlock()
	}

	if lockReport {
		funcName, ok := getLockCaller()
		if ok {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debugf("WorkRWLogicRUnLock,funcName:%s", funcName)
		} else {
			logrus.WithFields(logrus.Fields{"lockModuleName": l.moduleName}).Debug("WorkRWLogicRUnLock")
		}
	}
}
