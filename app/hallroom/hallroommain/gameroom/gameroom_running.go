package gameroom

import (
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/timeutil"
	"sort"
)

// 热血奔跑
const (
	ERunningReachForNormal  = 0
	ERunningReachForOffline = 1
)

// 玩家游戏数据
type RunningPlayerData struct {
	RoleID   int64  // 玩家ID
	Nickname string // 名称
	Icon     int32  // 图像
	Ready    int32  // 准备状态 0-未准备 1-准备了
	// 本次数据
	IsOnline  bool
	ReachTime int64 // 到达终点的时间
	ReachType int32 // 0-normal 1-offline
}
type RoomRunning struct {
	*EmptyRoom
	RoomStep       int32
	StepTime       int64
	ToPlayPlayers  []iroom.IGamePlayer
	PlayerGameData map[int64]*RunningPlayerData
	RunningRank    []int64
	GameTime       int32
	Distance       int32
}

func NewRoomRunning(roomID int64, globalObj iroom.IRoomGlobal) *RoomRunning {
	roomObj := &RoomRunning{
		EmptyRoom:      NewEmptyRoom(roomID, globalObj),
		RoomStep:       sconst.EUndercoverStepReady,
		StepTime:       timeutil.NowTime(),
		PlayerGameData: make(map[int64]*RunningPlayerData),
		RunningRank:    make([]int64, 0),
	}
	roomObj.RoomType = sconst.EGameRoomTypeRunning
	return roomObj
}
func (roomObj *RoomRunning) HoldPlayerData(roleID int64) *RunningPlayerData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	playerData := &RunningPlayerData{
		RoleID: roleID,
	}
	roomObj.PlayerGameData[roleID] = playerData
	return playerData
}
func (roomObj *RoomRunning) GetPlayerData(roleID int64) *RunningPlayerData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	return nil
}
func (roomObj *RoomRunning) JoinPlayer(p iroom.IGamePlayer) {
	roomObj.EmptyRoom.JoinPlayer(p)
	roomObj.ToPlayPlayers = append(roomObj.ToPlayPlayers, p) // 按进入房间的顺序
	playerData := roomObj.HoldPlayerData(p.GetRoleID())
	playerData.Nickname = p.GetName()
	playerData.Icon = p.GetIcon()
	playerData.IsOnline = true
}
func (roomObj *RoomRunning) PushRoomGameData() {
	pushList := make([]int64, 0)
	for k := range roomObj.PlayerList {
		pushList = append(pushList, k)
	}
	if len(pushList) > 0 {
		pushMsg := &pbclient.ECMsgGamePushRunningRoomDataNotify{
			RoomGameData: roomObj.ToGameDetailData(),
		}
		roomObj.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassGame, protocol.ECMsgGamePushRunningRoomData, pushMsg, pushList)
		loghlp.Debugf("PushRunningRoomGameData:%+v", pushMsg)
	}
}
func (roomObj *RoomRunning) ToGameDetailData() *pbclient.RoomRunningDetail {
	gameData := &pbclient.RoomRunningDetail{
		RoomStep: roomObj.RoomStep,
		StepTime: roomObj.StepTime,
	}
	// 玩家数据
	gameData.PlayersGameData = make(map[int64]*pbclient.RunningPlayerGameData)
	for k, v := range roomObj.PlayerGameData {
		gameData.PlayersGameData[k] = &pbclient.RunningPlayerGameData{
			RoleID:   k,
			Ready:    v.Ready,
			Icon:     v.Icon,
			Nickname: v.Nickname,
			IsOnline: v.IsOnline,
		}
	}

	if gameData.RoomStep == sconst.ERunningStepEnd {
		// 排名数据
		gameData.RankList = make([]int64, len(roomObj.RunningRank))
		copy(gameData.RankList, roomObj.RunningRank)
	}
	return gameData
}

func (roomObj *RoomRunning) ChangeStep(step int32) {
	roomObj.RoomStep = step
	roomObj.StepTime = timeutil.NowTime()
	loghlp.Debugf("RoomRunning room(%d) change_to_step(%d)", roomObj.RoomID, roomObj.RoomStep)
	// 广播房间数据
	roomObj.PushRoomGameData()
}
func (roomObj *RoomRunning) IsAllReady() bool {
	if len(roomObj.PlayerList) < 2 {
		return false
	}
	rReady := true
	for _, p := range roomObj.PlayerList {
		if playerData, ok := roomObj.PlayerGameData[p.GetRoleID()]; ok {
			if playerData.Ready == 0 {
				rReady = false
				break
			}
		} else {
			rReady = false
			break
		}
	}
	return rReady
}

func (roomObj *RoomRunning) Update(curTime int64) {
	roomObj.EmptyRoom.Update(curTime)
	switch roomObj.RoomStep {
	case sconst.ERunningStepReady:
		{
			if roomObj.IsAllReady() {
				roomObj.initDataForReady()
				roomObj.ChangeStep(sconst.ERunningStepRunning)
			}
			break
		}
	case sconst.ERunningStepRunning:
		{
			break
		}
	case sconst.ERunningStepEnd:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.resetDataForGameEnd()
				roomObj.ChangeStep(sconst.ERunningStepReady)
			}
			break
		}
	}
}
func (roomObj *RoomRunning) resetDataForGameEnd() {
	for _, p := range roomObj.PlayerGameData {
		if _, ok := roomObj.PlayerList[p.RoleID]; !ok {
			delete(roomObj.PlayerGameData, p.RoleID)
			continue
		}
		p.Ready = 0
		p.ReachTime = 0
		p.ReachType = 0
	}
	for _, p := range roomObj.EmptyRoom.PlayerList {
		p.SetReady(0)
	}
}

func (roomObj *RoomRunning) initDataForReady() {
	for _, p := range roomObj.PlayerGameData {
		p.ReachTime = 0
		p.ReachType = 0
	}
}
func (roomObj *RoomRunning) IsCanJoin() bool {
	return roomObj.RoomStep == sconst.ERunningStepReady
}

func (roomObj *RoomRunning) LeavePlayer(roleID int64) {
	// 排队列表删除
	for i, p := range roomObj.ToPlayPlayers {
		if p.GetRoleID() == roleID {
			for j := i; j+1 < len(roomObj.ToPlayPlayers); j++ {
				roomObj.ToPlayPlayers[j] = roomObj.ToPlayPlayers[j+1]
			}
			roomObj.ToPlayPlayers = roomObj.ToPlayPlayers[0:(len(roomObj.ToPlayPlayers) - 1)]
			break
		}
	}

	// 自动视为出局
	playerData := roomObj.HoldPlayerData(roleID)
	playerData.IsOnline = false
	if playerData.ReachTime == 0 {
		playerData.ReachTime = timeutil.NowTime()
		playerData.ReachType = ERunningReachForOffline
	}

	roomObj.EmptyRoom.LeavePlayer(roleID)
}

func (roomObj *RoomRunning) calcRank() {
	rankList := make([]*RunningPlayerData, 0)
	for _, p := range roomObj.PlayerGameData {
		rankList = append(rankList, p)
	}
	sort.SliceStable(rankList, func(i, j int) bool {
		if rankList[i].ReachType != rankList[j].ReachType {
			return rankList[i].ReachType < rankList[j].ReachType
		}
		if rankList[i].ReachTime != rankList[j].ReachTime {
			if rankList[i].ReachType == 0 {
				return rankList[i].ReachTime < rankList[j].ReachTime
			} else {
				return rankList[i].ReachTime > rankList[j].ReachTime
			}
		}
		return rankList[i].RoleID < rankList[j].RoleID
	})
	roomObj.RunningRank = make([]int64, len(rankList))
	for i, v := range rankList {
		roomObj.RunningRank[i] = v.RoleID
	}
}

func (roomObj *RoomRunning) CheckGameEnd() {
	// 如果只剩下一个人了,游戏结束
	reachNum := int32(0)
	for _, p := range roomObj.PlayerGameData {
		if p.ReachTime > 0 {
			reachNum++
		}
	}
	if reachNum >= int32(len(roomObj.PlayerGameData)) {
		loghlp.Infof("reachNum >= int32(len(roomObj.PlayerGameData)), game end")
		roomObj.calcRank()
		roomObj.ChangeStep(sconst.ERunningStepEnd)
	}
}
func (roomObj *RoomRunning) ToRoomDetail() *pbclient.RoomData {
	roomData := roomObj.EmptyRoom.ToRoomDetail()
	// 游戏数据
	roomData.RunningRoomData = roomObj.ToGameDetailData()
	// 结果
	return roomData
}
