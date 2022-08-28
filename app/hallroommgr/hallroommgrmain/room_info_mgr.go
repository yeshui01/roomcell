package hallroommgrmain

import (
	"math/rand"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/trnode"
	"time"
)

// 房间信息
type RoomInfo struct {
	RoomID     int64
	CreateTime int64

	RoomNode *trnode.TRNodeInfo // 所在的room节点
}

type RoomInfoManager struct {
	roomList          map[int64]*RoomInfo
	lockPendingCreate map[int64]*RoomInfo
	seqID             int64
	reportTime        int64
}

func newRoomInfoMgr() *RoomInfoManager {
	mgr := &RoomInfoManager{
		roomList:          make(map[int64]*RoomInfo),
		seqID:             100000,
		lockPendingCreate: make(map[int64]*RoomInfo),
		reportTime:        timeutil.NowTime(),
	}
	rand.Seed(time.Now().UnixNano())
	return mgr
}

// 生成房间id
func (mgr *RoomInfoManager) GenRoomUid() int64 {
	// seqID : 10000-99998 随机一个
	seqID := rand.Intn(89999)
	roomUid := (seqID+10000)*10 + int(trframe.GetCurNodeIndex())
	loghlp.Debugf("GenRoomUid:%d", roomUid)
	var genCount int32 = 1
	for {
		_, ok := mgr.roomList[int64(roomUid)]
		pendingCreate, ok2 := mgr.lockPendingCreate[int64(roomUid)]
		if (!ok && !ok2) || (!ok && ok2 && (timeutil.NowTime()-pendingCreate.CreateTime) > 30) {
			break
		} else {
			seqID++
			if roomUid >= 89999 {
				seqID = 0
			}
			roomUid = (seqID+10000)*10 + int(trframe.GetCurNodeIndex())
			genCount = genCount + 1
			loghlp.Debugf("GenRoomUidNext:%d", roomUid)
		}
		if genCount >= 89999 {
			roomUid = 0
			loghlp.Errorf("gen roomid count(%d), exit", genCount)
			break
		}
	}
	return int64(roomUid)
}

// 获取房间数量
func (mgr *RoomInfoManager) GetRoomNum() int32 {
	return int32(len(mgr.roomList))
}

// 查找房间
func (mgr *RoomInfoManager) FindRoomInfoById(roomID int64) *RoomInfo {
	if r, ok := mgr.roomList[roomID]; ok {
		return r
	}
	return nil
}
func (mgr *RoomInfoManager) AddRoomInfo(roomInfo *RoomInfo) {
	mgr.roomList[roomInfo.RoomID] = roomInfo
	loghlp.Infof("addRoomInfo,nowRoomNum:%d", len(mgr.roomList))
}

// 缓存房间创建
func (mgr *RoomInfoManager) PendingRoomCreate(roomID int64) *RoomInfo {
	roomInfo := &RoomInfo{
		RoomID:     roomID,
		CreateTime: timeutil.NowTime(),
	}
	mgr.lockPendingCreate[roomID] = roomInfo
	return roomInfo
}
func (mgr *RoomInfoManager) FindPendingRoomCreate(roomID int64) *RoomInfo {
	if r, ok := mgr.lockPendingCreate[roomID]; ok {
		return r
	}
	return nil
}
func (mgr *RoomInfoManager) DeletePendingRoomCreate(roomID int64) {
	delete(mgr.lockPendingCreate, roomID)
}

func (mgr *RoomInfoManager) DeleteRoomInfo(roomID int64) {
	delete(mgr.roomList, roomID)
	loghlp.Infof("DeleteRoomInfo(%d),nowRoomNum:%d", roomID, len(mgr.roomList))
}

func (mgr *RoomInfoManager) SecUpdate(curTime int64) {
	var delList []int64
	for k, v := range mgr.lockPendingCreate {
		if curTime-v.CreateTime >= 30 {
			delList = append(delList, k)
		}
	}
	for _, v := range delList {
		loghlp.Errorf("delete lock pending room create:%d", v)
		delete(mgr.lockPendingCreate, v)
	}
	if curTime-mgr.reportTime >= 60 {
		mgr.ReportStatus()
		mgr.reportTime = curTime
	}
}
func (mgr *RoomInfoManager) ReportStatus() {
	loghlp.Infof("RoomInfoManager::ReportStatus,roomNum:%d,pendingCreateNum:%d", len(mgr.roomList), len(mgr.lockPendingCreate))
}
