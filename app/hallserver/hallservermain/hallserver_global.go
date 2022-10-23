package hallservermain

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/ormdef"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
	"time"
)

type HallServerGlobal struct {
	playerList     map[int64]*HallPlayer
	lastUpdateTime int64    // 上次更新的时间戳(秒)
	onceJobList    []func() // 单次job列表
	reportTime     int64
}

func NewHallServerGlobal() *HallServerGlobal {
	return &HallServerGlobal{
		playerList: make(map[int64]*HallPlayer),
		reportTime: timeutil.NowTime(),
	}
}

func (hg *HallServerGlobal) FindPlayer(roleID int64) *HallPlayer {
	if p, ok := hg.playerList[roleID]; ok {
		return p
	}

	return nil
}
func (hg *HallServerGlobal) AddPlayer(p *HallPlayer) {
	if _, ok := hg.playerList[p.RoleID]; ok {
		loghlp.Errorf("player is existed:%d", p.RoleID)
		return
	}
	hg.playerList[p.RoleID] = p
}

func (hg *HallServerGlobal) Update(curTimeMs int64) {
	hg.SecUpdate(curTimeMs / 1000)
	hg.updateJob(curTimeMs)
}

func (hg *HallServerGlobal) SecUpdate(curTime int64) {
	if curTime <= hg.lastUpdateTime {
		return
	}
	hg.UpdateSavePlayer(curTime)
	for _, v := range hg.playerList {
		v.SecUpdate(curTime)
	}
	hg.UpdateRemoveIdlePlayer(curTime)
	if curTime-hg.reportTime >= 60 {
		hg.ReportStatus()
		hg.reportTime = curTime
	}
	hg.lastUpdateTime = curTime
}
func (hg *HallServerGlobal) UpdateSavePlayer(curTime int64) {
	var saveCount int = 0
	for _, player := range hg.playerList {
		if curTime-player.LastSaveTime >= 30 {
			// 先固定30秒保存一次
			hg.SavePlayer(player)
			saveCount++
			if saveCount >= 200 {
				// 每次最多保存200个
				break
			}
		}
	}
}
func (hg *HallServerGlobal) SavePlayer(player *HallPlayer) {
	reqDB := &pbserver.ESMsgPlayerSaveRoleReq{
		RoleID: player.RoleID,
	}
	roleBase := &pbserver.DbTableData{
		TableID: ormdef.ETableRoleBase,
	}
	roleBase.Data, _ = player.tbRoleBase.ToBytes()
	reqDB.RoleTables = append(reqDB.RoleTables, roleBase)

	// 发送到db保存
	// 发送消息
	roleID := player.RoleID
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("save player(%d) callback success,okCode:%d", roleID, okCode)
		return
	}
	// 这里的session是frameSession
	cbEnv := trframe.MakeMsgEnv(0, nil)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerSaveRole,
		reqDB,
		trnode.ETRNodeTypeHallData,
		0,
		cb,
		cbEnv,
	)
	player.LastSaveTime = time.Now().Unix()
}

func (hg *HallServerGlobal) HandlePlayerOnline(player *HallPlayer) {
	player.Online()
}
func (hg *HallServerGlobal) HandlePlayerOffline(player *HallPlayer) {
	loghlp.Infof("HandlePlayerOffline:%d", player.RoleID)
	player.Offline()
}
func (hg *HallServerGlobal) UpdateRemoveIdlePlayer(curTime int64) {
	var idlePlayers []int64
	for k, v := range hg.playerList {
		if !v.IsOnline && curTime-v.OfflineTime >= 300 {
			loghlp.Debugf("will remove idleplayer(%d) for offline timeout, offlineTime:%d", k, v.OfflineTime)
			idlePlayers = append(idlePlayers, k)
			if len(idlePlayers) >= 200 {
				// 一次移除200个
				break
			}
		} else if v.IsOnline {
			if v.GetHeartTime() == 0 {
				v.UpdateHeartTime(curTime)
			} else if curTime-v.GetHeartTime() >= 300 {
				loghlp.Debugf("will remove idleplayer(%d) for hearttime timeout, hearttime:%d", k, v.GetHeartTime())
				idlePlayers = append(idlePlayers, k)
				if len(idlePlayers) >= 200 {
					// 一次移除200个
					break
				}
			}
		}
	}
	for _, v := range idlePlayers {
		loghlp.Infof("remove idle player(%d)", v)
		delete(hg.playerList, v)
	}
}
func (hg *HallServerGlobal) PostJob(doJob func()) {
	hg.onceJobList = append(hg.onceJobList, doJob)
}

func (hg *HallServerGlobal) updateJob(curTimeMs int64) {
	if len(hg.onceJobList) > 0 {
		for _, doJob := range hg.onceJobList {
			doJob()
		}
		hg.onceJobList = nil
	}
}

func (hg *HallServerGlobal) ReportStatus() {
	loghlp.Infof("ReportStatus,playerNum:%d,onceJobList:%d", len(hg.playerList), len(hg.onceJobList))
}
