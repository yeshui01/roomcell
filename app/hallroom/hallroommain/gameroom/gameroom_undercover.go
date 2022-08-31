/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-08-22 15:28:05
 * @game: 谁是卧底
 */
package gameroom

import (
	"math/rand"
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/configdata"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/timeutil"
)

// 玩家卧底游戏数据
type PlayerUndercoverData struct {
	RoleID   int64 // 角色id
	Nickname string
	Icon     int32

	IsOut      bool   // 本轮是否出局
	SelfWords  string // 玩家拿到的词语
	CurTalk    string // 当前发言
	PlayNumber int32  // 玩家游戏编号(1开始)
	VoteNum    int32  // 本次得票数
	Ready      int32
	Voted      bool // 本轮是否已经投票
}

func NewPlayerUndercoverData(roleID int64) *PlayerUndercoverData {
	p := &PlayerUndercoverData{
		RoleID:     roleID,
		IsOut:      false,
		SelfWords:  "",
		CurTalk:    "",
		PlayNumber: 0,
		Ready:      0,
		Voted:      false,
	}

	return p
}

type RoomUndercover struct {
	*EmptyRoom
	RoomStep         int32
	StepTime         int64
	ToPlayPlayers    []iroom.IGamePlayer
	PlayerGameData   map[int64]*PlayerUndercoverData
	PlayNumber       int32   // number生成器
	TalkRoleID       int64   // 发言玩家id
	TalkRoleNumber   int32   // 发言玩家的编号
	UnderCoverWords  string  // 卧底词语
	OtherWords       string  // 其他玩家词语
	UndercoverRoleID int64   // 卧底玩家
	TalkerCacheList  []int64 // 发言玩家缓存列表
	IsUndercoverSucc bool    // 卧底是否胜利
}

func NewRoomUndercover(roomID int64, globalObj iroom.IRoomGlobal) *RoomUndercover {
	roomObj := &RoomUndercover{
		EmptyRoom:        NewEmptyRoom(roomID, globalObj),
		RoomStep:         sconst.EUndercoverStepReady,
		StepTime:         timeutil.NowTime(),
		PlayerGameData:   make(map[int64]*PlayerUndercoverData),
		PlayNumber:       1,
		TalkRoleID:       0,
		TalkRoleNumber:   0,
		UndercoverRoleID: 0,
		IsUndercoverSucc: false,
	}
	roomObj.RoomType = sconst.EGameRoomTypeUndercover
	return roomObj
}
func (roomObj *RoomUndercover) GetRoomType() int32 {
	return sconst.EGameRoomTypeUndercover
}
func (roomObj *RoomUndercover) JoinPlayer(p iroom.IGamePlayer) {
	roomObj.EmptyRoom.JoinPlayer(p)
	roomObj.ToPlayPlayers = append(roomObj.ToPlayPlayers, p) // 按进入房间的顺序
	playerData := roomObj.HoldPlayerData(p.GetRoleID())
	playerData.Nickname = p.GetName()
	playerData.Icon = p.GetIcon()
	if playerData.PlayNumber == 0 {
		// 分配玩家编号
		playerData.PlayNumber = roomObj.PlayNumber
		roomObj.PlayNumber = roomObj.PlayNumber + 1
	}
}

func (roomObj *RoomUndercover) LeavePlayer(roleID int64) {
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
	playerData.IsOut = true

	roomObj.EmptyRoom.LeavePlayer(roleID)
}
func (roomObj *RoomUndercover) HoldPlayerData(roleID int64) *PlayerUndercoverData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	playerData := NewPlayerUndercoverData(roleID)
	roomObj.PlayerGameData[roleID] = playerData
	return playerData
}
func (roomObj *RoomUndercover) GetPlayerData(roleID int64) *PlayerUndercoverData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	return nil
}
func (roomObj *RoomUndercover) Update(curTime int64) {
	roomObj.EmptyRoom.Update(curTime)
	switch roomObj.RoomStep {
	case sconst.EUndercoverStepReady:
		{
			if roomObj.IsAllReady() {
				roomObj.initPlayerNumber()

				roomObj.ChangeStep(sconst.EUndercoverStepGenWords)
			}
			break
		}
	case sconst.EUndercoverStepGenWords:
		{
			if len(roomObj.UnderCoverWords) < 1 {
				roomObj.genWords()
				// TODO 推送玩家词汇
			}
			// 3秒后变为发言
			if curTime-roomObj.StepTime >= 3 {
				roomObj.TalkRoleNumber = 1 // 初始化发言人为第一个
				for _, p := range roomObj.PlayerGameData {
					if p.PlayNumber == roomObj.TalkRoleNumber {
						roomObj.TalkRoleID = p.RoleID
						break
					}
				}
				roomObj.ChangeStep(sconst.EUndercoverStepTalk)
			}
			break
		}
	case sconst.EUndercoverStepTalk:
		{
			// DO NOTHING 等待玩家游戏
			break
		}
	case sconst.EUndercoverStepVote:
		{
			// DO NOTHING 等待玩家投票
			break
		}
	case sconst.EUndercoverStepVoteSummary:
		{
			// 停留5秒
			if curTime-roomObj.StepTime >= 5 {
				if roomObj.VoteSummary() {
					roomObj.ChangeStep(sconst.EUndercoverStepEnd)
				} else {
					// 继续发言投票
					roomObj.nextTalkTurn()
				}
			}
			break
		}
	case sconst.EUndercoverStepEnd:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.ChangeStep(sconst.EUndercoverStepReady)
			}
			break
		}
	}
}
func (roomObj *RoomUndercover) ChangeStep(step int32) {
	roomObj.RoomStep = step
	roomObj.StepTime = timeutil.NowTime()
	loghlp.Debugf("undercover room(%d) change_to_step(%d)", roomObj.RoomID, roomObj.RoomStep)
	// 广播房间数据
	roomObj.PushRoomGameData()
}

// 分配玩家编号
func (roomObj *RoomUndercover) initPlayerNumber() {
	roomObj.PlayNumber = 1
	roomObj.TalkerCacheList = make([]int64, len(roomObj.ToPlayPlayers))
	for i, p := range roomObj.ToPlayPlayers {
		playerData := roomObj.HoldPlayerData(p.GetRoleID())
		playerData.PlayNumber = roomObj.PlayNumber
		playerData.CurTalk = ""
		playerData.IsOut = false
		playerData.SelfWords = ""
		playerData.VoteNum = 0
		playerData.Voted = false
		playerData.Ready = 0
		roomObj.TalkerCacheList[i] = playerData.RoleID
		roomObj.PlayNumber = roomObj.PlayNumber + 1
	}

	roomObj.TalkRoleID = 0
	roomObj.TalkRoleNumber = 0
	roomObj.IsUndercoverSucc = false
	roomObj.UnderCoverWords = ""
	roomObj.OtherWords = ""
}

// 生成词语
func (roomObj *RoomUndercover) genWords() {
	roomObj.UnderCoverWords = "iphone"
	roomObj.OtherWords = "ipad"
	// 读取配置
	typeWordsList := configdata.Instance().GetUndercoverCfgList()
	wordsLen := len(typeWordsList)
	if wordsLen > 0 {
		idx := rand.Intn(wordsLen)
		roomObj.UnderCoverWords = typeWordsList[idx].Undercover
		roomObj.OtherWords = typeWordsList[idx].Other
	}
	// 随机分配玩家词语
	if len(roomObj.TalkerCacheList) < 3 {
		return
	}
	randRoleID := roomObj.TalkerCacheList[rand.Intn(len(roomObj.TalkerCacheList))]
	for _, playerData := range roomObj.PlayerGameData {
		if playerData.RoleID == randRoleID {
			// 卧底玩家
			playerData.SelfWords = roomObj.UnderCoverWords
			roomObj.UndercoverRoleID = playerData.RoleID
		} else {
			playerData.SelfWords = roomObj.OtherWords
		}
	}
	loghlp.Infof("undercover genWords finish, undercover(%s),other(%s),UndercoverRoleID(%d)", roomObj.UnderCoverWords, roomObj.OtherWords, roomObj.UndercoverRoleID)
}
func (roomObj *RoomUndercover) IsCanJoin() bool {
	return roomObj.RoomStep == sconst.EUndercoverStepReady
}

func (roomObj *RoomUndercover) IsAllReady() bool {
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
func (roomObj *RoomUndercover) PushRoomGameData() {
	pushList := make([]int64, 0)
	for k := range roomObj.PlayerList {
		pushList = append(pushList, k)
	}
	if len(pushList) > 0 {
		pushMsg := &pbclient.ECMsgGamePushUndercoverRoomDataNotify{
			RoomGameData: roomObj.ToGameDetailData(),
		}
		roomObj.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassGame, protocol.ECMsgGamePushUndercoverRoomData, pushMsg, pushList)
		loghlp.Debugf("PushUndercoverRoomGameData:%+v", pushMsg)
	}
}
func (roomObj *RoomUndercover) ToGameDetailData() *pbclient.RoomUndercoverDetail {
	gameData := &pbclient.RoomUndercoverDetail{
		RoomStep: roomObj.RoomStep,
		StepTime: roomObj.StepTime,
	}
	// 玩家数据
	gameData.PlayersGameData = make(map[int64]*pbclient.UndercoverPlayerGameData)
	for k, v := range roomObj.PlayerGameData {
		gameData.PlayersGameData[k] = &pbclient.UndercoverPlayerGameData{
			RoleID:       k,
			SelfWords:    v.SelfWords,
			Ready:        v.Ready,
			IsOut:        v.IsOut,
			PlayerNumber: v.PlayNumber,
			VoteNum:      v.VoteNum,
		}
	}
	gameData.UnderWin = roomObj.IsUndercoverSucc
	return gameData
}
func (roomObj *RoomUndercover) ToRoomDetail() *pbclient.RoomData {
	roomData := roomObj.EmptyRoom.ToRoomDetail()
	// 游戏数据
	roomData.UndercoverRoomData = roomObj.ToGameDetailData()
	// 结果
	return roomData
}

// 玩家结束卧底发言
func (roomObj *RoomUndercover) OnPlayerEndUnderTalk(roleID int64) {
	// 是否都发言完了
	var allTalked bool = true
	for k := range roomObj.EmptyRoom.PlayerList {
		playerData := roomObj.HoldPlayerData(k)
		if !playerData.IsOut {
			if len(playerData.CurTalk) < 1 {
				allTalked = false
				break
			}
		}
	}
	if allTalked {
		// 都发言完了
		loghlp.Infof("room(%d) all player has under talked", roomObj.RoomID)
		roomObj.ChangeStep(sconst.EUndercoverStepVote)
	} else {
		// 下一位发言人
		if roomObj.changeUnderTalker() {
			// 广播消息
			pushMsg := &pbclient.ECMsgGamePushUndercoverTalkerChangeNotify{
				TalkRoleID: roomObj.TalkRoleID,
			}
			roomObj.BroadCastRoomMsg(0,
				protocol.ECMsgClassGame,
				protocol.ECMsgGamePushUndercoverTalkerChange,
				pushMsg)
			loghlp.Infof("changeUnderTalker succ, room(%d) new underTalker(%d), push:%+v", roomObj.RoomID, roomObj.UndercoverRoleID, pushMsg)
		} else {
			// 结束,进入投票
			loghlp.Infof("change next UnderTalker finish, enter vote step")
			roomObj.ChangeStep(sconst.EUndercoverStepVote)
		}
	}
}
func (roomObj *RoomUndercover) changeUnderTalker() bool {
	var ret bool = false
	for idx := roomObj.TalkRoleNumber - 1; idx < int32(len(roomObj.TalkerCacheList)); idx++ {
		playerData := roomObj.GetPlayerData(roomObj.TalkerCacheList[idx])
		if playerData == nil {
			continue
		}
		if playerData.IsOut {
			continue
		}
		if len(playerData.CurTalk) > 0 {
			continue
		}
		// 变更
		roomObj.TalkRoleNumber = idx + 1
		roomObj.TalkRoleID = playerData.RoleID
		ret = true
		break
	}
	return ret
}

// 玩家投票
func (roomObj *RoomUndercover) OnPlayerVote(roleID int64, targetRoleID int64) {
	playerData := roomObj.GetPlayerData(roleID)
	if playerData == nil {
		loghlp.Errorf("OnPlayerVote, not find player_data(%d)", roleID)
		return
	}
	targetPlayerData := roomObj.GetPlayerData(targetRoleID)
	if targetPlayerData == nil {
		loghlp.Errorf("OnPlayerVote, not find target player_data(%d)", targetRoleID)
		return
	}
	playerData.Voted = true
	targetPlayerData.VoteNum = targetPlayerData.VoteNum + 1

	// 是否都投票了
	var allVoted bool = true
	for k := range roomObj.EmptyRoom.PlayerList {
		playerData := roomObj.HoldPlayerData(k)
		if !playerData.IsOut {
			if !playerData.Voted {
				allVoted = false
				break
			}
		}
	}
	if allVoted {
		// 计算得票数
		loghlp.Debugf("room(%d) on voted end", roomObj.RoomID)
		roomObj.ChangeStep(sconst.EUndercoverStepVoteSummary)
	}
}

// 投票结束
func (roomObj *RoomUndercover) VoteSummary() bool {
	loghlp.Debugf("room(%d) on vote summary", roomObj.RoomID)
	// 当有玩家被投票大于50%时，第一窗口“Zhuangtai”文字变更为“已出局”（播放相应动画），此时所有
	// 玩家“Zhuangtai/Toupiao”关闭，出局玩家两窗口锁死，不可再操作。
	// 如无玩家被投票达到50%，1号玩家继续发言。
	curPlayerNum := int32(len(roomObj.TalkerCacheList))
	outVoteNum := curPlayerNum / 2
	var undercoverPlayerOut bool = false
	for _, playerData := range roomObj.PlayerGameData {
		if playerData.IsOut {
			continue
		}
		if playerData.VoteNum >= outVoteNum {
			// 出局
			pushMsg := &pbclient.ECMsgGamePushUndercoverOutNotify{
				RoleID: playerData.RoleID,
			}
			playerData.IsOut = true
			roomObj.BroadCastRoomMsg(0,
				protocol.ECMsgClassGame,
				protocol.ECMsgGamePushUndercoverOut,
				pushMsg)

			if playerData.RoleID == roomObj.UndercoverRoleID {
				undercoverPlayerOut = true
			}

			loghlp.Debugf("push undercover player(%d) out for voted num >= 50%,push:%+v", playerData.RoleID, pushMsg)
		}
	}
	if undercoverPlayerOut {
		// 卧底出局了
		roomObj.IsUndercoverSucc = false
		return true
	}
	// 卧底没有出局
	if len(roomObj.PlayerList) <= 3 {
		roomObj.IsUndercoverSucc = true
		return true
	}
	return false
}

func (roomObj *RoomUndercover) nextTalkTurn() {
	// 继续发言投票
	loghlp.Debugf("undercover room(%d) nextTalkTurn", roomObj.RoomID)
	for _, playerData := range roomObj.PlayerGameData {
		playerData.CurTalk = ""
		playerData.VoteNum = 0
		playerData.Voted = false
	}
	// 下一位发言人
	roomObj.TalkRoleNumber = 1
	if roomObj.changeUnderTalker() {
		// 广播消息
		pushMsg := &pbclient.ECMsgGamePushUndercoverTalkerChangeNotify{
			TalkRoleID: roomObj.TalkRoleID,
		}
		roomObj.BroadCastRoomMsg(0,
			protocol.ECMsgClassGame,
			protocol.ECMsgGamePushUndercoverTalkerChange,
			pushMsg)
		loghlp.Infof("changeUnderTalker succ, room(%d) new underTalker(%d), push:%+v", roomObj.RoomID, roomObj.UndercoverRoleID, pushMsg)
		roomObj.ChangeStep(sconst.EUndercoverStepTalk)
	} else {
		// 结束,进入投票
		loghlp.Infof("change next UnderTalker finish2, enter vote step")
		roomObj.ChangeStep(sconst.EUndercoverStepVote)
	}
}
