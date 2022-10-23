package hallroommain

import (
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/sensitivewds"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
)

type HallRoomGlobal struct {
	playerList        map[int64]*RoomPlayer
	lastSecUpdateTime int64
	roomList          map[int64]iroom.IGameRoom
	onceJobList       []func() // 单次job列表
	reportTime        int64
	sensitiveWords    *sensitivewds.DFAUtil
}

func NewHallRoomGlobal() *HallRoomGlobal {
	g := &HallRoomGlobal{
		playerList:        make(map[int64]*RoomPlayer),
		lastSecUpdateTime: timeutil.NowTime(),
		roomList:          make(map[int64]iroom.IGameRoom),
		reportTime:        timeutil.NowTime(),
		sensitiveWords:    nil,
	}

	return g
}

func (mgr *HallRoomGlobal) Update(curTimeMs int64) {
	curTime := curTimeMs / 1000
	if curTime > mgr.lastSecUpdateTime {
		mgr.SecUpdate(curTime)
		mgr.lastSecUpdateTime = curTime
	}
	mgr.updateJob(curTimeMs)
}
func (mgr *HallRoomGlobal) SecUpdate(curTime int64) {
	for _, r := range mgr.roomList {
		r.Update(curTime)
	}
	mgr.updateIdlePlayer(curTime)
	mgr.updateEmptyRoom(curTime)
	if curTime-mgr.reportTime >= 60 {
		mgr.ReportStatus()
		mgr.reportTime = curTime
	}
}

func (mgr *HallRoomGlobal) FindRoomPlayer(roleID int64) *RoomPlayer {
	if p, ok := mgr.playerList[roleID]; ok {
		return p
	}

	return nil
}

func (mgr *HallRoomGlobal) AddRoomPlayer(p *RoomPlayer) {
	if _, ok := mgr.playerList[p.RoleID]; !ok {
		mgr.playerList[p.RoleID] = p
	}
}

func (mgr *HallRoomGlobal) DelRoomPlayer(p *RoomPlayer) {
	delete(mgr.playerList, p.RoleID)
}

func (s *HallRoomGlobal) BroadcastMsgToClient(msgClass int32, msgType int32, pbMsg proto.Message, roleList []int64) {
	msgData, err := proto.Marshal(pbMsg)
	if err != nil {
		return
	}
	pushMsg := &pbframe.EFrameMsgBroadcastMsgToClientReq{
		MsgClass: msgClass,
		MsgType:  msgType,
		MsgData:  msgData,
	}
	pushMsg.RoleList = make([]int64, len(roleList))
	copy(pushMsg.RoleList, roleList)
	trframe.BroadcastMessage(protocol.EMsgClassFrame, protocol.EFrameMsgBroadcastMsgToClient, pushMsg, trnode.ETRNodeTypeHallGate)
}

func (s *HallRoomGlobal) PostJob(doJob func()) {
	s.onceJobList = append(s.onceJobList, doJob)
}

func (s *HallRoomGlobal) updateJob(curTimeMs int64) {
	if len(s.onceJobList) > 0 {
		for _, doJob := range s.onceJobList {
			doJob()
		}
		s.onceJobList = nil
	}
}
func (mgr *HallRoomGlobal) updateIdlePlayer(curTime int64) {
	var delList []*RoomPlayer
	for _, r := range mgr.playerList {
		if r.HeartTime == 0 {
			r.HeartTime = curTime
		}
		if curTime-r.HeartTime >= sconst.PlayerHeartTime {
			loghlp.Warnf("updateIdleRoomPlayer(%d) heart timeout", r.GetRoleID())
			delList = append(delList, r)
		}
	}
	for _, v := range delList {
		if v.RoomPtr != nil {
			v.RoomPtr.LeavePlayer(v.GetRoleID())
		}
		delete(mgr.playerList, v.GetRoleID())
	}
}
func (mgr *HallRoomGlobal) ReportStatus() {
	loghlp.Infof("ReportStatus,playerNum:%d,roomNum:%d", len(mgr.playerList), len(mgr.roomList))
}

func (mgr *HallRoomGlobal) InitSensitiveWords(wordsFile string) {
	ssWords, err := sensitivewds.InitSensitiveWords(wordsFile)
	if err != nil {
		loghlp.Errorf("InitSensitiveWords(%s) Fail:%s", wordsFile, err.Error())
		return
	}
	loghlp.Infof("InitSensitiveWords(%s) success", wordsFile)
	mgr.sensitiveWords = ssWords
}

func (mgr *HallRoomGlobal) GetSensitiveWordsUtil() *sensitivewds.DFAUtil {
	return mgr.sensitiveWords
}
