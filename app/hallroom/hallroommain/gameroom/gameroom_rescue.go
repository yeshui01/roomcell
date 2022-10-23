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

// 拯救玩家
// 玩家游戏数据
type RescuePlayerData struct {
	RoleID   int64  // 玩家ID
	Nickname string // 名称
	Icon     int32  // 图像
	Ready    int32  // 准备状态 0-未准备 1-准备了
	// 本次数据
	IsOnline bool
	Hp       int32
	HpTime   int64 // hp时间
}
type RoomRescue struct {
	*EmptyRoom
	RoomStep       int32
	StepTime       int64
	ToPlayPlayers  []iroom.IGamePlayer
	PlayerGameData map[int64]*RescuePlayerData
	MaxHp          int32
	MaxTime        int32
	HpRank         []int64
}

func NewRoomRescue(roomID int64, globalObj iroom.IRoomGlobal) *RoomRescue {
	roomObj := &RoomRescue{
		EmptyRoom:      NewEmptyRoom(roomID, globalObj),
		RoomStep:       sconst.EUndercoverStepReady,
		StepTime:       timeutil.NowTime(),
		PlayerGameData: make(map[int64]*RescuePlayerData),
		MaxHp:          10,
		MaxTime:        120,
		HpRank:         make([]int64, 0),
	}
	roomObj.RoomType = sconst.EGameRoomTypeRescuePlayer
	return roomObj
}
func (roomObj *RoomRescue) HoldPlayerData(roleID int64) *RescuePlayerData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	playerData := &RescuePlayerData{
		RoleID: roleID,
	}
	roomObj.PlayerGameData[roleID] = playerData
	return playerData
}
func (roomObj *RoomRescue) GetPlayerData(roleID int64) *RescuePlayerData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	return nil
}
func (roomObj *RoomRescue) JoinPlayer(p iroom.IGamePlayer) {
	roomObj.EmptyRoom.JoinPlayer(p)
	roomObj.ToPlayPlayers = append(roomObj.ToPlayPlayers, p) // 按进入房间的顺序
	playerData := roomObj.HoldPlayerData(p.GetRoleID())
	playerData.Nickname = p.GetName()
	playerData.Icon = p.GetIcon()
	playerData.IsOnline = true
}

func (roomObj *RoomRescue) PushRoomGameData() {
	pushList := make([]int64, 0)
	for k := range roomObj.PlayerList {
		pushList = append(pushList, k)
	}
	if len(pushList) > 0 {
		pushMsg := &pbclient.ECMsgGamePushRescueRoomDataNotify{
			RoomGameData: roomObj.ToGameDetailData(),
		}
		roomObj.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassGame, protocol.ECMsgGamePushRescueRoomData, pushMsg, pushList)
		loghlp.Debugf("PushRescueRoomGameData:%+v", pushMsg)
	}
}
func (roomObj *RoomRescue) ToGameDetailData() *pbclient.RoomRescueDetail {
	gameData := &pbclient.RoomRescueDetail{
		RoomStep: roomObj.RoomStep,
		StepTime: roomObj.StepTime,
		MaxTime:  roomObj.MaxTime,
	}
	// 玩家数据
	gameData.PlayersGameData = make(map[int64]*pbclient.RescuePlayerGameData)
	for k, v := range roomObj.PlayerGameData {
		gameData.PlayersGameData[k] = &pbclient.RescuePlayerGameData{
			RoleID:   k,
			Ready:    v.Ready,
			Icon:     v.Icon,
			Nickname: v.Nickname,
			IsOnline: v.IsOnline,
			Hp:       v.Hp,
		}
	}
	gameData.MaxHp = roomObj.MaxHp
	gameData.MaxTime = roomObj.MaxTime
	if gameData.RoomStep == sconst.ERescueStepEnd {
		// 排名数据
		gameData.RankList = make([]int64, len(roomObj.HpRank))
		copy(gameData.RankList, roomObj.HpRank)
	}
	return gameData
}
func (roomObj *RoomRescue) ChangeStep(step int32) {
	roomObj.RoomStep = step
	roomObj.StepTime = timeutil.NowTime()
	loghlp.Debugf("RoomRescue room(%d) change_to_step(%d)", roomObj.RoomID, roomObj.RoomStep)
	// 广播房间数据
	roomObj.PushRoomGameData()
}
func (roomObj *RoomRescue) IsAllReady() bool {
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
func (roomObj *RoomRescue) Update(curTime int64) {
	roomObj.EmptyRoom.Update(curTime)
	switch roomObj.RoomStep {
	case sconst.ERescueStepReady:
		{
			if roomObj.IsAllReady() {
				roomObj.initDataForReady()
				roomObj.ChangeStep(sconst.ERescueStepRunning)
			}
			break
		}
	case sconst.ERescueStepRunning:
		{
			if curTime-roomObj.StepTime >= int64(roomObj.MaxTime) {
				loghlp.Infof("rescue room(%d) game end for time reach", roomObj.RoomID)
				roomObj.calcRank()
				roomObj.ChangeStep(sconst.ERescueStepEnd)
			}
			break
		}
	case sconst.ERescueStepEnd:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.resetDataForGameEnd()
				roomObj.ChangeStep(sconst.ERescueStepReady)
			}
			break
		}
	}
}
func (roomObj *RoomRescue) resetDataForGameEnd() {
	for _, p := range roomObj.PlayerGameData {
		if _, ok := roomObj.PlayerList[p.RoleID]; !ok {
			delete(roomObj.PlayerGameData, p.RoleID)
			continue
		}
		p.Ready = 0
		p.Hp = 0
		p.HpTime = 0
	}
	for _, p := range roomObj.EmptyRoom.PlayerList {
		p.SetReady(0)
	}
}

func (roomObj *RoomRescue) initDataForReady() {
	for _, p := range roomObj.PlayerGameData {
		p.Hp = roomObj.MaxHp
	}
}
func (roomObj *RoomRescue) IsCanJoin() bool {
	return roomObj.RoomStep == sconst.ERescueStepReady
}
func (roomObj *RoomRescue) LeavePlayer(roleID int64) {
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
	playerData.Hp = 0
	playerData.IsOnline = false
	roomObj.EmptyRoom.LeavePlayer(roleID)
}

func (roomObj *RoomRescue) calcRank() {
	rankList := make([]*RescuePlayerData, 0)
	for _, p := range roomObj.PlayerGameData {
		rankList = append(rankList, p)
	}
	sort.SliceStable(rankList, func(i, j int) bool {
		if rankList[i].Hp != rankList[j].Hp {
			return rankList[i].Hp > rankList[j].Hp
		}
		if rankList[i].HpTime != rankList[j].HpTime {
			return rankList[i].HpTime > rankList[j].HpTime
		}
		return rankList[i].RoleID < rankList[j].RoleID
	})
	roomObj.HpRank = make([]int64, len(rankList))
	for i, v := range rankList {
		roomObj.HpRank[i] = v.RoleID
	}
}

func (roomObj *RoomRescue) CheckGameEnd() {
	// 如果只剩下一个人了,游戏结束
	aliveNum := int32(0)
	for _, p := range roomObj.PlayerGameData {
		if p.Hp < 1 {
			continue
		}
		aliveNum++
	}
	if aliveNum <= 1 {
		loghlp.Infof("alive num <= 1, game end")
		roomObj.calcRank()
		roomObj.ChangeStep(sconst.ERescueStepEnd)
	}
}
func (roomObj *RoomRescue) ToRoomDetail() *pbclient.RoomData {
	roomData := roomObj.EmptyRoom.ToRoomDetail()
	// 游戏数据
	roomData.RescueRoomData = roomObj.ToGameDetailData()
	// 结果
	return roomData
}
