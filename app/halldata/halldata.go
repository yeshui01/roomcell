package halldata

import (
	"time"

	"roomcell/app/halldata/halldatahandler"
	"roomcell/app/halldata/halldatamain"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type HallData struct {
	hallDataGlobal *halldatamain.HallDataGlobal
}

func NewHallData() *HallData {
	s := &HallData{
		hallDataGlobal: halldatamain.NewHallDataGlobal(),
	}
	s.RegisterMsgHandler()
	halldatahandler.InitHallDataObj(s)
	return s
}

// 打开本地数据库
func (hdata *HallData) OpenLocalDB(dbFile string) bool {
	dbLogger := logger.New(
		logrus.New(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		panic("open dbfile faile")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)
	hdata.hallDataGlobal.InitGameDB(db)
	return true
}

// 打开Mysql数据库
func (hdata *HallData) OpenMysqlDB(connStr string) bool {
	dbLogger := logger.New(
		logrus.New(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		panic("open dbfile faile")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)
	hdata.hallDataGlobal.InitGameDB(db)
	return true
}
func (hdata *HallData) RegisterMsgHandler() {
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoadRole,
		halldatahandler.HandlePlayerLoadRoleData)
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerSaveRole,
		halldatahandler.HandlePlayerSaveRoleData)
}

func (hdata *HallData) GetGameDB() *gorm.DB {
	return hdata.hallDataGlobal.GetGameDB()
}
func (hdata *HallData) HallDataGlobal() *halldatamain.HallDataGlobal {
	return hdata.hallDataGlobal
}
func (hdata *HallData) FrameRun(curTimeMs int64) {
	hdata.hallDataGlobal.Update(curTimeMs)
}
