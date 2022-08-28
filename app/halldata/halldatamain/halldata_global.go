package halldatamain

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/timeutil"
	"time"

	"gorm.io/gorm"
)

type HallDataDBJob struct {
	DoJob func()
}

type HallDataGlobal struct {
	gameDB         *gorm.DB
	dataPlayers    map[int64]*HallDataPlayer
	lastUpdateTime int64
	dbJobCh        chan *HallDataDBJob
	dbJobStop      bool
}

func NewHallDataGlobal() *HallDataGlobal {
	return &HallDataGlobal{
		gameDB:      nil,
		dataPlayers: make(map[int64]*HallDataPlayer),
		dbJobCh:     make(chan *HallDataDBJob, 1024),
		dbJobStop:   false,
	}
}

func (servGlobal *HallDataGlobal) GetGameDB() *gorm.DB {
	return servGlobal.gameDB
}
func (servGlobal *HallDataGlobal) InitGameDB(gameDB *gorm.DB) {
	servGlobal.gameDB = gameDB
	servGlobal.StartDBRun()
}

func (servGlobal *HallDataGlobal) FindDataPlayer(roleID int64) *HallDataPlayer {
	if p, ok := servGlobal.dataPlayers[roleID]; ok {
		p.VisitTime = timeutil.NowTime()
		return p
	}
	return nil
}
func (servGlobal *HallDataGlobal) AddDataPlayer(roleID int64, player *HallDataPlayer) {
	if _, ok := servGlobal.dataPlayers[roleID]; ok {
		return
	}
	player.VisitTime = timeutil.NowTime()
	servGlobal.dataPlayers[roleID] = player
}
func (servGlobal *HallDataGlobal) Update(curTimeMs int64) {
	servGlobal.SecUpdate(curTimeMs / 1000)
}
func (servGlobal *HallDataGlobal) SecUpdate(curTime int64) {
	if curTime <= servGlobal.lastUpdateTime {
		return
	}
	var delList []int64
	for _, v := range servGlobal.dataPlayers {
		if v.VisitTime == 0 {
			v.VisitTime = curTime
		}
		if curTime-v.VisitTime >= 1200 {
			delList = append(delList, v.RoleID)
		}
	}
	for _, v := range delList {
		delete(servGlobal.dataPlayers, v)
	}
	servGlobal.lastUpdateTime = curTime
}

// DB更新线程
func (servGlobal *HallDataGlobal) StartDBRun() {
	go func() {
		for {
			select {
			case dbOpt, ok := <-servGlobal.dbJobCh:
				{
					if ok {
						loghlp.Debugf("do db job")
						dbOpt.DoJob()
					} else {
						servGlobal.dbJobStop = true
					}
					break
				}
			}
			if servGlobal.dbJobStop {
				break
			}
		}
		loghlp.Info("exit db run")
	}()
}
func (servGlobal *HallDataGlobal) StopDBRun() {
	close(servGlobal.dbJobCh)
	// 等待结束
	for {
		if servGlobal.dbJobStop {
			time.Sleep(time.Second)
			break
		} else {
			time.Sleep(time.Second)
		}
	}
}
func (servGlobal *HallDataGlobal) PostDBJob(dbJob *HallDataDBJob) {
	servGlobal.dbJobCh <- dbJob
}
