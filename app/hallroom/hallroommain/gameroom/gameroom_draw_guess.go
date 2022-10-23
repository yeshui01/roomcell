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

// var WordsTypeDef = []string{
// 	"食物",
// 	"水果",
// 	"蔬菜",
// 	"调料",
// 	"饮料",
// 	"零食",
// 	"干果",
// 	"医学",
// 	"职业",
// 	"乐器",
// 	"财经",
// 	"汽车",
// 	"交通",
// 	"动物",
// 	"生活用品",
// 	"动漫",
// 	"自然现象",
// 	"运动",
// 	"身体部位",
// 	"植物",
// 	"武器",
// 	"成语"}

// 玩家游戏数据
type PlayerDrawData struct {
	RoleID       int64  // 玩家ID
	Nickname     string // 昵称
	Score        int32  // 积分
	GuessCorrect int32  // 猜对的次数
	Ready        int32  // 准备状态 0-未准备 1-准备了
	// 本次游戏猜测
	CurWords string
	CurScore int32
}
type SelectWords struct {
	WordType int32  //类型
	Words    string // 单词
}
type RoomDrawGuess struct {
	*EmptyRoom
	RoomStep      int32
	StepTime      int64
	ToDrawPlayers []iroom.IGamePlayer
	// 房主的游戏设定
	MaxTurnNum int32 // 游戏轮数
	DrawTime   int32 // 每次游戏绘画时间(秒)
	CurTurn    int32 // 当前第几轮
	// 当前房间游戏数据
	CurWords       string // 当前单词
	CurWordType    int32
	WordsToSelect  []*SelectWords // 待选词组列表
	AnswerList     []int64        // 回答正确的玩家列表,按顺序
	DrawerRoleID   int64          // 当前画图的玩家
	PlayerGameData map[int64]*PlayerDrawData
	DrawEndTime    int64 // 本次画画结束时间,会随着猜对的人数动态变化
	// 画图数据
	DrawOpts []*pbclient.DrawPainData
	// 当前画画的idx
	DrawerIndex int32
	// 当前分数
	RewardScore int32
}

func (roomObj *RoomDrawGuess) GetRoomType() int32 {
	return sconst.EGameRoomTypeDrawGuess
}

func NewRoomDrawGuess(roomID int64, globalObj iroom.IRoomGlobal) *RoomDrawGuess {
	r := &RoomDrawGuess{
		EmptyRoom:      NewEmptyRoom(roomID, globalObj),
		RoomStep:       sconst.EDrawGuessStepReady,
		StepTime:       timeutil.NowTime(),
		PlayerGameData: make(map[int64]*PlayerDrawData),
		DrawTime:       45,
		MaxTurnNum:     1,
	}
	r.RoomType = sconst.EGameRoomTypeDrawGuess
	return r
}
func (roomObj *RoomDrawGuess) ChangeStep(step int32) {
	roomObj.RoomStep = step
	roomObj.StepTime = timeutil.NowTime()
	loghlp.Debugf("room(%d) change to step(%d),drawRoleId:%d", roomObj.RoomID, roomObj.RoomStep, roomObj.DrawerRoleID)
	// 广播房间数据
	roomObj.PushRoomGameData()
}

func (roomObj *RoomDrawGuess) Update(curTime int64) {
	roomObj.EmptyRoom.Update(curTime)
	switch roomObj.RoomStep {
	case sconst.EDrawGuessStepReady:
		{
			// 是否都准备好了
			if roomObj.IsAllReady() {
				roomObj.ResetPlayerPlayData()
				roomObj.DrawerIndex = -1
				roomObj.CurTurn = 1
				roomObj.ChangeStep(sconst.EDrawGuessStepSelectDrawer)
			}
			break
		}
	case sconst.EDrawGuessStepSelectDrawer:
		{
			if curTime-roomObj.StepTime >= 3 {
				roomObj.ResetPlayerThisTimesData()
				roomObj.SelectCurDrawer()
				roomObj.GenWordsSelectList()
				roomObj.RewardScore = 10 // 测试分数
				roomObj.DrawEndTime = timeutil.NowTime() + int64(roomObj.DrawTime)
				roomObj.CurWords = ""
				roomObj.CurWordType = 0
				roomObj.AnswerList = make([]int64, 0)
				roomObj.ChangeStep(sconst.EDrawGuessStepSelectWords) // 放到选词之后
			}
			break
		}
	case sconst.EDrawGuessStepSelectWords:
		{
			// 等待玩家选词 DO NOThING
			break
		}
	case sconst.EDrawGuessStepDraw:
		{
			// 你画我猜
			if curTime >= roomObj.DrawEndTime {
				// 结束
				loghlp.Infof("room(%d) cur draw timeout, endgame", roomObj.RoomID)
				roomObj.ChangeStep(sconst.EDrawGuessStepEnd)
			}
			break
		}
	case sconst.EDrawGuessStepEnd:
		{
			// 本次游戏结束了
			// 下一位开始画
			if curTime-roomObj.StepTime >= 3 {
				roomObj.DrawOpts = make([]*pbclient.DrawPainData, 0)
				loghlp.Debugf("room(%d) step end, DrawerIndex:%d", roomObj.RoomID, roomObj.DrawerIndex)
				// 进入下一位
				if roomObj.DrawerIndex >= int32(len(roomObj.ToDrawPlayers)-1) {
					loghlp.Debugf("room(%d) this turn all player draw end, DrawerIndex:%d", roomObj.RoomID, roomObj.DrawerIndex)
					// 所有人都画了
					if roomObj.CurTurn >= roomObj.MaxTurnNum {
						// 最后一轮游戏结束
						// 回到准备状态
						//roomObj.ChangeStep(sconst.EDrawGuessStepReady)
						loghlp.Infof("room(%d)", roomObj.RoomID)
						loghlp.Debugf("room(%d) this turn end, MaxTurnNum:%d,CurTurn:%d", roomObj.RoomID, roomObj.MaxTurnNum, roomObj.CurTurn)
						roomObj.ChangeStep(sconst.EDrawGuessStepGameEnd)
					} else {
						// 下一轮
						roomObj.DrawerIndex = -1
						roomObj.CurTurn++
						loghlp.Debugf("room(%d) will enter next turn, MaxTurnNum:%d,NextCurTurn:%d, ", roomObj.RoomID, roomObj.MaxTurnNum, roomObj.CurTurn)
						//清理单词
						roomObj.CurWords = ""
						roomObj.CurWordType = 0
						roomObj.WordsToSelect = make([]*SelectWords, 0)
						roomObj.ChangeStep(sconst.EDrawGuessStepSelectDrawer)
					}
				} else {
					loghlp.Infof("room(%d) will start next draw", roomObj.RoomID)
					roomObj.ChangeStep(sconst.EDrawGuessStepSelectDrawer)
				}
			}
			break
		}
	case sconst.EDrawGuessStepGameEnd:
		{
			if curTime-roomObj.StepTime >= 5 {
				// 清理玩家游戏数据
				roomObj.ResetPlayerPlayData()
				roomObj.DrawerIndex = -1
				roomObj.CurTurn = 1
				roomObj.DrawerRoleID = 0
				roomObj.CurWords = ""
				roomObj.CurWordType = 0
				roomObj.WordsToSelect = make([]*SelectWords, 0)
				roomObj.ChangeStep(sconst.EDrawGuessStepReady)
			}
			break
		}
	}
}

func (roomObj *RoomDrawGuess) ToGameDetailData() *pbclient.RoomDrawDetail {
	gameData := &pbclient.RoomDrawDetail{
		MaxTurnNum: roomObj.MaxTurnNum,
		DrawTime:   roomObj.DrawTime,
		CurTurn:    roomObj.CurTurn,
		RoomStep:   roomObj.RoomStep,
		StepTime:   roomObj.StepTime,
		CurWords:   roomObj.CurWords,
		//WordsType:  roomObj.CurWordType,
		//WordsToSelect: roomObj.WordsToSelect,
		DrawerRoleID: roomObj.DrawerRoleID,
		DrawOpts:     roomObj.DrawOpts,
		DrawEndTime:  roomObj.DrawEndTime,
	}
	wordTypeCfg := configdata.Instance().GetWordTypeCfg(roomObj.CurWordType)
	if wordTypeCfg != nil {
		gameData.WordsType = wordTypeCfg.TypeName
	} else {
		gameData.WordsType = "unknown"
	}
	for _, v := range roomObj.WordsToSelect {
		gameData.WordsToSelect = append(gameData.WordsToSelect, v.Words)
	}
	// 玩家数据
	gameData.PlayersGameData = make(map[int64]*pbclient.DrawPlayerGameData)
	for k, v := range roomObj.PlayerGameData {
		gameData.PlayersGameData[k] = &pbclient.DrawPlayerGameData{
			RoleID:       k,
			TotalScore:   v.Score,
			GuessCorrect: v.GuessCorrect,
			Score:        v.CurScore,
			Nickname:     v.Nickname,
		}
	}
	return gameData
}

func (roomObj *RoomDrawGuess) ToRoomDetail() *pbclient.RoomData {
	roomData := roomObj.EmptyRoom.ToRoomDetail()
	// 游戏数据
	roomData.DrawRoomData = roomObj.ToGameDetailData()
	return roomData
}

func (roomObj *RoomDrawGuess) JoinPlayer(p iroom.IGamePlayer) {
	roomObj.EmptyRoom.JoinPlayer(p)
	playerData := roomObj.HoldPlayerData(p.GetRoleID())
	playerData.Nickname = p.GetName()
	roomObj.ToDrawPlayers = append(roomObj.ToDrawPlayers, p) // 按进入房间的顺序
}

func (roomObj *RoomDrawGuess) LeavePlayer(roleID int64) {
	// 排队列表删除
	for i, p := range roomObj.ToDrawPlayers {
		if p.GetRoleID() == roleID {
			// 移除待绘画的人
			for j := i; j+1 < len(roomObj.ToDrawPlayers); j++ {
				roomObj.ToDrawPlayers[j] = roomObj.ToDrawPlayers[j+1]
			}
			roomObj.ToDrawPlayers = roomObj.ToDrawPlayers[0:(len(roomObj.ToDrawPlayers) - 1)]
			break
		}
	}
	roomObj.EmptyRoom.LeavePlayer(roleID)
}
func (roomObj *RoomDrawGuess) PushRoomGameData() {
	pushList := make([]int64, 0)
	for k := range roomObj.PlayerList {
		pushList = append(pushList, k)
	}
	if len(pushList) > 0 {
		pushMsg := &pbclient.ECMsgGamePushDrawRoomDataNotify{
			RoomGameData: roomObj.ToGameDetailData(),
		}
		roomObj.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassGame, protocol.ECMsgGamePushDrawRoomData, pushMsg, pushList)
		loghlp.Debugf("PushRoomGameData:%+v", pushMsg)
	}
}

// 生成单词选择列表
func (roomObj *RoomDrawGuess) GenWordsSelectList() {
	typeWordsList := configdata.Instance().GetWordsTypeList()
	typeLen := len(typeWordsList)
	if typeLen < 1 {
		loghlp.Errorf("no word types")
		return
	}
	// 随机4种
	selectTypeList := make([]*configdata.DrawTypeWordsCfg, 0)
	if typeLen > 4 {
		for {
			if len(selectTypeList) >= 4 {
				break
			}
			idx := rand.Intn(typeLen)
			selectTypeList = append(selectTypeList, typeWordsList[idx])
			tmp := typeWordsList[typeLen-1]
			typeWordsList[typeLen-1] = typeWordsList[idx]
			typeWordsList[idx] = tmp
			typeWordsList = typeWordsList[0 : typeLen-1]
			typeLen--
		}
	} else {
		copy(selectTypeList, typeWordsList)
	}
	// 每个类型随机一个单词
	roomObj.WordsToSelect = make([]*SelectWords, 0)
	for _, selectType := range selectTypeList {
		if len(selectType.WordsList) < 1 {
			continue
		}
		idx := rand.Intn(len(selectType.WordsList))
		if idx >= 0 && idx < len(selectType.WordsList) {
			roomObj.WordsToSelect = append(roomObj.WordsToSelect, &SelectWords{
				Words:    selectType.WordsList[idx].Word,
				WordType: selectType.WordType,
			})
		}
	}
	loghlp.Infof("GenWordsSelectList,room(%d), to select words:%+v", roomObj.RoomID, roomObj.WordsToSelect)
}

// 选择当前画画的玩家
func (roomObj *RoomDrawGuess) SelectCurDrawer() {
	roomObj.DrawerIndex++
	loghlp.Debugf("SelectCurDrawer,toDrawPlayerLen:%d,idx:%d", int32(len(roomObj.ToDrawPlayers)), roomObj.DrawerIndex)
	if roomObj.DrawerIndex >= int32(len(roomObj.ToDrawPlayers)) {
		// 所有人都画了,这里不应该走到
		return
	}
	roomObj.DrawerRoleID = roomObj.ToDrawPlayers[roomObj.DrawerIndex].GetRoleID()
	playData := roomObj.HoldPlayerData(roomObj.DrawerRoleID)
	playData.CurWords = ""
	playData.CurScore = 0
}
func (roomObj *RoomDrawGuess) IsAllReady() bool {
	if len(roomObj.PlayerList) < 2 {
		return false // 至少需要两个人
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
func (roomObj *RoomDrawGuess) HoldPlayerData(roleID int64) *PlayerDrawData {
	if playerData, ok := roomObj.PlayerGameData[roleID]; ok {
		return playerData
	}
	playerData := &PlayerDrawData{
		RoleID:   roleID,
		Score:    0,
		CurWords: "",
		CurScore: 0,
	}
	roomObj.PlayerGameData[roleID] = playerData
	return playerData
}

func (roomObj *RoomDrawGuess) CheckDrawEndTime() {
	// roomObj.DrawEndTime
	// 猜对人数,超过1/3时,如果剩余时间大于10秒,那么剩余时间变为10秒
	n := len(roomObj.AnswerList)
	if n >= len(roomObj.PlayerGameData)/3 {
		curTime := timeutil.NowTime()
		if roomObj.DrawEndTime-curTime > 10 {
			roomObj.DrawEndTime = curTime + 10
			loghlp.Infof("CheckDrawEndTime, dynamic change draw end time:%d", roomObj.DrawEndTime)
		}
	}
}
func (roomObj *RoomDrawGuess) CheckEndRightNow() {
	// 全部猜对或者猜错,直接结束
	correctNum := 0
	errNum := 0
	for _, p := range roomObj.PlayerGameData {
		if len(p.CurWords) < 1 {
			continue
		}
		if p.CurWords == roomObj.CurWords {
			correctNum++
		} else {
			errNum++
		}
	}
	playerNum := len(roomObj.PlayerGameData) - 1 // 减去1个画图的人
	if errNum == playerNum || correctNum == playerNum {
		loghlp.Infof("draw room(%d) end Right now,correctNum:%d, errorNum:%d", roomObj.RoomID, correctNum, errNum)
		roomObj.ChangeStep(sconst.EDrawGuessStepEnd)
		return
	}
	if errNum+correctNum >= playerNum {
		// 都回答完了,结束
		loghlp.Infof("draw room(%d) end Right now for all player answered,correctNum:%d, errorNum:%d", roomObj.RoomID, correctNum, errNum)
		playerPlayData := roomObj.HoldPlayerData(roomObj.DrawerRoleID)
		playerPlayData.CurScore = 10 // 每次10分
		playerPlayData.Score = playerPlayData.Score + 10
		roomObj.ChangeStep(sconst.EDrawGuessStepEnd)
	}
}

// 清理玩家本次游戏数据
func (roomObj *RoomDrawGuess) ResetPlayerThisTimesData() {
	for _, p := range roomObj.PlayerGameData {
		p.CurScore = 0
		p.CurWords = ""
		roomObj.CurWordType = 0
	}
}

func (roomObj *RoomDrawGuess) IsCanJoin() bool {
	return roomObj.RoomStep == sconst.EDrawGuessStepReady
}

// 清理玩家本次游戏数据
func (roomObj *RoomDrawGuess) ResetPlayerPlayData() {
	for _, p := range roomObj.PlayerGameData {
		p.CurScore = 0
		p.CurWords = ""
		roomObj.CurWordType = 0
		p.Score = 0
		p.GuessCorrect = 0
		p.Ready = 0
	}
	for _, p := range roomObj.EmptyRoom.PlayerList {
		p.SetReady(0)
	}
}
