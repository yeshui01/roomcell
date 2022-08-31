package gameroom

import (
	"math/rand"
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/timeutil"
)

// 数字炸弹
// 玩家游戏数据
type PlayerBombData struct {
	RoleID   int64  // 玩家ID
	Nickname string // 名称
	Icon     int32  // 图像
	Ready    int32  // 准备状态 0-未准备 1-准备了
	// 本次游戏猜测
	IsOnline   bool
	GuessNum   int32 // 本次猜测的数字
	PlayNumber int32 // 玩家编号
	IsTalked   bool
}

type RoomNumberBomb struct {
	*EmptyRoom
	RoomStep        int32
	StepTime        int64
	ToPlayPlayers   []iroom.IGamePlayer
	PlayerGameData  map[int64]*PlayerBombData
	TalkerCacheList []*PlayerBombData // 发言玩家缓存列表
	TalkRoleID      int64             // 发言玩家id
	TalkRoleNumber  int32             // 发言玩家的编号
	SysNumber       int32             // 系统生成的数字
	BombRoleID      int64             // 最后炸弹玩家
	MinNumber       int32             // 最小数字
	MaxNumber       int32             // 最大数字
	PlayNumber      int32
	MaxTurn         int32 // 最大轮数
	CurTurn         int32 // 当前轮数
}

func NewRoomNumberBomb(roomID int64, globalObj iroom.IRoomGlobal) *RoomNumberBomb {
	roomObj := &RoomNumberBomb{
		EmptyRoom:      NewEmptyRoom(roomID, globalObj),
		RoomStep:       sconst.EUndercoverStepReady,
		StepTime:       timeutil.NowTime(),
		PlayerGameData: make(map[int64]*PlayerBombData),
		TalkRoleID:     0,
		TalkRoleNumber: 0,
		SysNumber:      0,
		BombRoleID:     0,
		MinNumber:      1,
		MaxNumber:      99,
		PlayNumber:     0,
	}
	roomObj.RoomType = sconst.EGameRoomTypeNumberBomb
	return roomObj
}
func (roomObj *RoomNumberBomb) HoldPlayerData(roleID int64) *PlayerBombData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	playerData := &PlayerBombData{
		RoleID: roleID,
	}
	roomObj.PlayerGameData[roleID] = playerData
	return playerData
}
func (roomObj *RoomNumberBomb) GetPlayerData(roleID int64) *PlayerBombData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	return nil
}
func (roomObj *RoomNumberBomb) JoinPlayer(p iroom.IGamePlayer) {
	roomObj.EmptyRoom.JoinPlayer(p)
	roomObj.ToPlayPlayers = append(roomObj.ToPlayPlayers, p) // 按进入房间的顺序
	playerData := roomObj.HoldPlayerData(p.GetRoleID())
	playerData.Nickname = p.GetName()
	playerData.Icon = p.GetIcon()
	playerData.IsOnline = true
	if playerData.PlayNumber == 0 {
		// 分配玩家编号
		playerData.PlayNumber = roomObj.genNumberId()
	}
}
func (roomObj *RoomNumberBomb) PushRoomGameData() {
	pushList := make([]int64, 0)
	for k := range roomObj.PlayerList {
		pushList = append(pushList, k)
	}
	if len(pushList) > 0 {
		pushMsg := &pbclient.ECMsgGamePushNumberBombRoomDataNotify{
			RoomGameData: roomObj.ToGameDetailData(),
		}
		roomObj.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassGame, protocol.ECMsgGamePushUndercoverRoomData, pushMsg, pushList)
		loghlp.Debugf("PushUndercoverRoomGameData:%+v", pushMsg)
	}
}
func (roomObj *RoomNumberBomb) ToGameDetailData() *pbclient.RoomNumberBombDetail {
	gameData := &pbclient.RoomNumberBombDetail{
		RoomStep:       roomObj.RoomStep,
		StepTime:       roomObj.StepTime,
		MinNumber:      roomObj.MinNumber,
		MaxNumber:      roomObj.MaxNumber,
		TalkRoleID:     roomObj.TalkRoleID,
		TalkRoleNumber: roomObj.TalkRoleNumber,
		BombRoleID:     roomObj.BombRoleID,
		Turn:           roomObj.CurTurn,
		MaxTurn:        roomObj.MaxTurn,
	}
	// 玩家数据
	gameData.PlayersGameData = make(map[int64]*pbclient.NumberBombPlayerGameData)
	for k, v := range roomObj.PlayerGameData {
		gameData.PlayersGameData[k] = &pbclient.NumberBombPlayerGameData{
			RoleID:       k,
			Ready:        v.Ready,
			PlayerNumber: v.PlayNumber,
			Icon:         v.Icon,
			Nickname:     v.Nickname,
			IsOnline:     v.IsOnline,
			GuessNum:     v.GuessNum,
		}
	}
	gameData.MinNumber = roomObj.MinNumber
	gameData.MaxNumber = roomObj.MaxNumber
	return gameData
}
func (roomObj *RoomNumberBomb) ToRoomDetail() *pbclient.RoomData {
	roomData := roomObj.EmptyRoom.ToRoomDetail()
	// 游戏数据
	roomData.NumberbombRoomData = roomObj.ToGameDetailData()
	// 结果
	return roomData
}
func (roomObj *RoomNumberBomb) ChangeStep(step int32) {
	roomObj.RoomStep = step
	roomObj.StepTime = timeutil.NowTime()
	loghlp.Debugf("RoomNumberBomb room(%d) change_to_step(%d)", roomObj.RoomID, roomObj.RoomStep)
	// 广播房间数据
	roomObj.PushRoomGameData()
}
func (roomObj *RoomNumberBomb) IsAllReady() bool {
	if len(roomObj.PlayerList) < 4 {
		return false // 至少需要4个人
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

func (roomObj *RoomNumberBomb) Update(curTime int64) {
	roomObj.EmptyRoom.Update(curTime)
	switch roomObj.RoomStep {
	case sconst.EUndercoverStepReady:
		{
			if roomObj.IsAllReady() {
				roomObj.initPlayerNumber()
				roomObj.ChangeStep(sconst.ENumberBombStepGenNumber)
			}
			break
		}
	case sconst.ENumberBombStepGenNumber:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.genSysNumber()
				roomObj.nextTalker(false)
				roomObj.ChangeStep(sconst.ENumberBombStepGuessNumber)
			}
			break
		}
	case sconst.ENumberBombStepGuessNumber:
		{
			break
		}
	case sconst.ENumberBombStepTurnEnd:
		{
			if curTime-roomObj.StepTime >= 3 {
				if roomObj.CurTurn >= roomObj.MaxTurn {
					roomObj.ChangeStep(sconst.ENumberBombStepGameEnd)
				} else {
					roomObj.ChangeStep(sconst.ENumberBombStepGenNumber)
				}
			}
			break
		}
	case sconst.ENumberBombStepGameEnd:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.resetDataForGameEnd()
				roomObj.ChangeStep(sconst.ENumberBombStepReady)
			}
			break
		}
	}
}

func (roomObj *RoomNumberBomb) resetDataForGameEnd() {
	for _, p := range roomObj.PlayerGameData {
		if _, ok := roomObj.PlayerList[p.RoleID]; !ok {
			delete(roomObj.PlayerGameData, p.RoleID)
			continue
		}
		p.Ready = 0
		p.GuessNum = 0
		p.PlayNumber = 0
		p.IsTalked = false
	}
	roomObj.TalkRoleNumber = 0
	roomObj.TalkRoleID = 0
	roomObj.SysNumber = 0
	roomObj.MinNumber = sconst.NumberBombMinNumber
	roomObj.MaxNumber = sconst.NumberBombMaxNumber // 默认值
	roomObj.BombRoleID = 0
}

func (roomObj *RoomNumberBomb) initPlayerNumber() {
	roomObj.PlayNumber = 0
	roomObj.TalkerCacheList = make([]*PlayerBombData, 0)
	for _, p := range roomObj.ToPlayPlayers {
		playerData, ok := roomObj.PlayerGameData[p.GetRoleID()]
		if ok {
			playerData.PlayNumber = roomObj.genNumberId()
			roomObj.TalkerCacheList = append(roomObj.TalkerCacheList, playerData)
		}
	}
}

func (roomObj *RoomNumberBomb) genSysNumber() {
	roomObj.SysNumber = int32(rand.Intn(int(sconst.NumberBombMaxNumber))) + 1
	loghlp.Infof("numberbomb room(%d) genSysNumber:%d", roomObj.RoomID, roomObj.SysNumber)
}

func (roomObj *RoomNumberBomb) genNumberId() int32 {
	roomObj.PlayNumber++
	return roomObj.PlayNumber
}

func (roomObj *RoomNumberBomb) PlayerGuess(roleID int64, guessNumber int32) bool {
	playerData := roomObj.HoldPlayerData(roleID)
	if playerData == nil {
		return false
	}

	if guessNumber < sconst.NumberBombMinNumber || guessNumber > sconst.NumberBombMaxNumber {
		loghlp.Errorf("player(%d) guess number param error, guessNumber:%d", roleID, guessNumber)
		return false
	}
	if guessNumber == roomObj.SysNumber {
		//
		loghlp.Infof("player(%d) guess number bomb sysnumber(%d)", roleID, guessNumber)
		roomObj.BombRoleID = roleID
		roomObj.ChangeStep(sconst.ENumberBombStepTurnEnd)
		return true
	} else if guessNumber > roomObj.SysNumber {
		roomObj.MaxNumber = guessNumber - 1
		roomObj.PushRangeChange()
	} else {
		roomObj.MinNumber = guessNumber + 1
		roomObj.PushRangeChange()
	}
	playerData.IsTalked = true
	if !roomObj.nextTalker(true) {
		// 找不到下一位发言人,继续重新开始发言
		for _, p := range roomObj.PlayerGameData {
			p.IsTalked = false
			p.GuessNum = 0
		}
		roomObj.TalkRoleID = 0
		roomObj.TalkRoleNumber = 0
		if !roomObj.nextTalker(true) {
			// 结束
			loghlp.Infof("After PlayerGuess, not find next taler, turn will end")
			roomObj.ChangeStep(sconst.ENumberBombStepTurnEnd)
		}
	}
	return true
}

func (roomObj *RoomNumberBomb) PushRangeChange() {
	pushMsg := &pbclient.ECMsgGamePushNumberBombRangeChangeNotify{
		MinNumber: roomObj.MinNumber,
		MaxNumber: roomObj.MaxNumber,
	}
	roomObj.BroadCastRoomMsg(0,
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushNumberBombRangeChange,
		pushMsg)
}

// 变更发言人
func (roomObj *RoomNumberBomb) nextTalker(notify bool) bool {
	var findNewTalker bool = false
	for _, p := range roomObj.TalkerCacheList {
		if p.GuessNum != 0 {
			continue
		}
		if !p.IsOnline {
			// 不在线,直接跳过
			p.IsTalked = true
			p.GuessNum = 0
			continue
		}
		p.IsTalked = true
		findNewTalker = true
		pushMsg := &pbclient.ECMsgGamePushNumberBombGuesserChangeNotify{
			TalkRoleID: p.RoleID,
		}
		roomObj.TalkRoleNumber = p.PlayNumber
		roomObj.TalkRoleID = p.RoleID
		if notify {
			roomObj.BroadCastRoomMsg(0,
				protocol.ECMsgClassGame,
				protocol.ECMsgGamePushNumberBombGuesserChange,
				pushMsg)
		}
		break
	}
	if !findNewTalker {
		roomObj.TalkRoleNumber = 0
		roomObj.TalkRoleID = 0
	}
	return findNewTalker
}
