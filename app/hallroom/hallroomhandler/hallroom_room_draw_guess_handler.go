package hallroomhandler

import (
	"roomcell/app/hallroom/hallroommain/gameroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"strings"
)

// 玩家画画
func HandlePlayerDrawPaint(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameDrawPaintReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameDrawPaintRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	hallPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if hallPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if hallPlayer.RoomPtr == nil {
		loghlp.Errorf("room player(%d) roomptr is nil", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotInRoom, rep, iframe.EHandleContent
	}
	roomDrawGuess, ok := hallPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
	if !ok {
		loghlp.Errorf("roomDrawGuess convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	roomDrawGuess.DrawOpts = append(roomDrawGuess.DrawOpts, req.CurPain)
	// 广播给其他
	pushMsg := &pbclient.ECMsgGamePushDrawPaintNotify{
		RoleID:  hallPlayer.GetRoleID(),
		CurPain: req.CurPain,
	}
	roomDrawGuess.BroadCastRoomMsg(hallPlayer.GetRoleID(),
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushDrawPaint,
		pushMsg,
	)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 玩家猜词
func HandlePlayerGuessWords(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameDrawGuessWordsReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameDrawGuessWordsRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	hallPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if hallPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if hallPlayer.RoomPtr == nil {
		loghlp.Errorf("room player(%d) roomptr is nil", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotInRoom, rep, iframe.EHandleContent
	}
	roomDrawGuess, ok := hallPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
	if !ok {
		loghlp.Errorf("roomDrawGuess convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomDrawGuess.RoomStep != sconst.EDrawGuessStepDraw {
		loghlp.Errorf("roomDrawGuess.RoomStep != sconst.EDrawGuessStepDraw")
		return protocol.ECodeRoomDrawInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomDrawGuess.HoldPlayerData(hallPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if playerPlayData.CurWords != "" {
		return protocol.ECodeRoomDrawReapeatedGuess, rep, iframe.EHandleContent
	}
	playerPlayData.CurWords = req.Words

	pushMsg := &pbclient.ECMsgGamePushDrawGuessNotify{
		RoleID: hallPlayer.GetRoleID(),
		Words:  req.Words,
	}
	// 广播
	roomDrawGuess.BroadCastRoomMsg(hallPlayer.GetRoleID(),
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushDrawGuess,
		pushMsg,
	)
	if req.Words == roomDrawGuess.CurWords {
		// 猜对了
		roomDrawGuess.AnswerList = append(roomDrawGuess.AnswerList, hallPlayer.GetRoleID())
		playerPlayData.GuessCorrect++
		playerPlayData.Score = playerPlayData.Score + roomDrawGuess.RewardScore // 累计积分
		playerPlayData.CurScore = roomDrawGuess.RewardScore
		if roomDrawGuess.RewardScore > 1 {
			roomDrawGuess.RewardScore--
			// 每次获得的分数递减,最少1分
		}
		roomDrawGuess.CheckDrawEndTime()
		roomDrawGuess.CheckEndRightNow()
	} else {
		// 猜错了
		roomDrawGuess.CheckDrawEndTime()
		roomDrawGuess.CheckEndRightNow()
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 设定游戏规则
func HandleDrawGuessSetting(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameGrawSettingReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameGrawSettingRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	hallPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if hallPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if hallPlayer.RoomPtr == nil {
		loghlp.Errorf("room player(%d) roomptr is nil", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotInRoom, rep, iframe.EHandleContent
	}
	roomDrawGuess, ok := hallPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
	if !ok {
		loghlp.Errorf("roomDrawGuess convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomDrawGuess.MasterID != hallPlayer.GetRoleID() {
		return protocol.ECodeRoomDrawSettingAuthError, rep, iframe.EHandleContent
	}
	roomDrawGuess.MaxTurnNum = req.MaxTurnNum
	roomDrawGuess.DrawTime = req.DrawTime
	// 推送刷新
	pushSetting := &pbclient.ECMsgGamePushDrawSettingNotify{
		MaxTurnNum: req.MaxTurnNum,
		DrawTime:   req.DrawTime,
	}
	roomDrawGuess.BroadCastRoomMsg(hallPlayer.GetRoleID(),
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushDrawSetting,
		pushSetting,
	)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 选词
func HandleDrawSelectWords(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameDrawSelectWordsReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameDrawSelectWordsRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	hallPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if hallPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if hallPlayer.RoomPtr == nil {
		loghlp.Errorf("room player(%d) roomptr is nil", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotInRoom, rep, iframe.EHandleContent
	}
	roomDrawGuess, ok := hallPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
	if !ok {
		loghlp.Errorf("roomDrawGuess convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomDrawGuess.DrawerRoleID != hallPlayer.GetRoleID() {
		return protocol.ECodeRoomDrawSettingAuthError, rep, iframe.EHandleContent
	}
	if roomDrawGuess.RoomStep != sconst.EDrawGuessStepSelectWords {
		return protocol.ECodeRoomDrawInvalideOption, rep, iframe.EHandleContent
	}
	if req.Idx < 0 || req.Idx >= int32(len(roomDrawGuess.WordsToSelect)) {
		return protocol.ECodeParamError, rep, iframe.EHandleContent
	}
	roomDrawGuess.CurWords = roomDrawGuess.WordsToSelect[req.Idx].Words
	roomDrawGuess.CurWordType = roomDrawGuess.WordsToSelect[req.Idx].WordType
	roomDrawGuess.DrawEndTime = timeutil.NowTime() + int64(roomDrawGuess.DrawTime)
	roomDrawGuess.ChangeStep(sconst.EDrawGuessStepDraw)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 聊天猜词
func HandlePlayerChatGuessWords(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgRoomChatReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	loghlp.Infof("HandlePlayerChatGuessWords")
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgRoomChatRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	roomPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if roomPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if roomPlayer.RoomPtr == nil {
		loghlp.Errorf("room player(%d) roomptr is nil", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotInRoom, rep, iframe.EHandleContent
	}
	roomDrawGuess, ok := roomPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
	if !ok {
		loghlp.Errorf("roomDrawGuess convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomDrawGuess.RoomStep != sconst.EDrawGuessStepDraw {
		loghlp.Errorf("roomDrawGuess.RoomStep != sconst.EDrawGuessStepDraw")
		return protocol.ECodeRoomDrawInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomDrawGuess.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	// if playerPlayData.CurWords != "" {
	// 	return protocol.ECodeRoomDrawReapeatedGuess, rep, iframe.EHandleContent
	// }
	// playerPlayData.CurWords = req.TalkContent
	// 推送聊天
	pushMsg := &pbclient.ECMsgRoomPushChatNotify{
		TalkContent: req.TalkContent,
		Talker: &pbclient.RoomTalker{
			RoleID:   roomPlayer.GetRoleID(),
			Nickname: roomPlayer.Nickname,
			Icon:     roomPlayer.Icon,
		},
	}
	// // // 屏蔽字处理
	// wordsUtil := hallRoomGlobal.GetSensitiveWordsUtil()
	// if wordsUtil != nil {
	// 	req.TalkContent = wordsUtil.HandleWord(req.TalkContent, 'x')
	// }
	if strings.Index(req.TalkContent, roomDrawGuess.CurWords) != -1 {
		loghlp.Infof("player(%d) guess words(%s), talkContent:%s", roomPlayer.GetRoleID(), roomDrawGuess.CurWords, req.TalkContent)
		if len(playerPlayData.CurWords) < 1 {
			loghlp.Infof("player(%d) guess correct,talkContent:%s", roomPlayer.GetRoleID(), req.TalkContent)
			playerPlayData.CurWords = req.TalkContent
			// 猜对了
			roomDrawGuess.AnswerList = append(roomDrawGuess.AnswerList, roomPlayer.GetRoleID())
			playerPlayData.GuessCorrect++
			playerPlayData.Score = playerPlayData.Score + roomDrawGuess.RewardScore
			playerPlayData.CurScore = roomDrawGuess.RewardScore
			// 画画的人也加分
			drawerData := roomDrawGuess.HoldPlayerData(roomDrawGuess.DrawerRoleID)
			drawerData.Score = drawerData.Score + roomDrawGuess.RewardScore
			drawerData.CurScore = roomDrawGuess.RewardScore
			if roomDrawGuess.RewardScore > 1 {
				roomDrawGuess.RewardScore--
				// 每次获得的分数递减,最少1分
			}
		}
		// 答案替换为*
		nLen := len(roomDrawGuess.CurWords)
		hideStr := strings.Repeat("*", nLen)
		loghlp.Debugf("hideStr:%s,CurWords:%s", hideStr, roomDrawGuess.CurWords)
		req.TalkContent = strings.ReplaceAll(req.TalkContent, roomDrawGuess.CurWords, hideStr)
		if roomPlayer.RoomPtr != nil {
			pushMsg.TalkContent = req.TalkContent
			roomPlayer.RoomPtr.BroadCastRoomMsg(0, protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, pushMsg)
		}
		//roomDrawGuess.CheckDrawEndTime()
		// roomDrawGuess.CheckEndRightNow()
	} else {
		if roomPlayer.RoomPtr != nil {
			pushMsg.TalkContent = req.TalkContent
			roomPlayer.RoomPtr.BroadCastRoomMsg(0, protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, pushMsg)
		}
		// 猜错了
		//roomDrawGuess.CheckDrawEndTime()
		//roomDrawGuess.CheckEndRightNow()
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
