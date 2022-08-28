package hallroomhandler

import (
	"roomcell/app/hallroom/hallroommain/gameroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"strings"
)

// 卧底发言
func HandlePlayerChatUndertalk(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomUndercover)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep != sconst.EUndercoverStepTalk {
		loghlp.Errorf("roomObj.RoomStep != sconst.EUndercoverStepTalk")
		return protocol.ECodeRoomUndercoverInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if playerPlayData.IsOut {
		loghlp.Errorf("playerPlayData.IsOut ECodeRoomUndercoverOutPlayerCantTalk")
		return protocol.ECodeRoomUndercoverOutPlayerCantTalk, rep, iframe.EHandleContent
	}
	// 推送聊天
	pushMsg := &pbclient.ECMsgRoomPushChatNotify{
		TalkContent: req.TalkContent,
		Talker: &pbclient.RoomTalker{
			RoleID:   roomPlayer.GetRoleID(),
			Nickname: roomPlayer.Nickname,
			Icon:     roomPlayer.Icon,
		},
	}
	if len(playerPlayData.CurTalk) == 0 {
		playerPlayData.CurTalk = req.TalkContent
		if strings.Index(req.TalkContent, playerPlayData.SelfWords) != -1 {
			loghlp.Warnf("player(%d) undercover talk,selfWords(%s), talkContent:%s", roomPlayer.GetRoleID(), playerPlayData.SelfWords, req.TalkContent)
			// 自己的词语替换为*防止作弊
			nLen := len(playerPlayData.SelfWords)
			hideStr := strings.Repeat("*", nLen)
			loghlp.Debugf("hideStr:%s,CurWords:%s", hideStr, playerPlayData.SelfWords)
			req.TalkContent = strings.ReplaceAll(req.TalkContent, playerPlayData.SelfWords, hideStr)
			if roomPlayer.RoomPtr != nil {
				pushMsg.TalkContent = req.TalkContent
				roomPlayer.RoomPtr.BroadCastRoomMsg(0, protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, pushMsg)
			}
		} else {
			loghlp.Infof("player(%d) undercover talk normal,selfWords(%s), talkContent:%s", roomPlayer.GetRoleID(), playerPlayData.SelfWords, req.TalkContent)
			if roomPlayer.RoomPtr != nil {
				pushMsg.TalkContent = req.TalkContent
				roomPlayer.RoomPtr.BroadCastRoomMsg(0, protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, pushMsg)
			}
		}
		roomObj.OnPlayerEndUnderTalk(roomPlayer.GetRoleID())
	} else {
		loghlp.Warnf("undercover player(%d) has talked, cant repeated talk, ignore talk!", roomPlayer.GetRoleID())
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 投票
func HandlePlayerUndercoverVote(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameUndercoverVoteReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	if req.TargetRoleID < 1 {
		return protocol.ECodeParamError, nil, iframe.EHandleContent
	}
	loghlp.Infof("HandlePlayerUndercoverVote")
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameUndercoverVoteRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	roomPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if roomPlayer == nil {
		loghlp.Errorf("not find room player:%d", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeRoomPlayerNotFound, rep, iframe.EHandleContent
	}
	if roomPlayer.RoomPtr != nil {
		roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomUndercover)
		if !ok {
			loghlp.Errorf("roomObj convert fail")
			return protocol.ECodeSysError, rep, iframe.EHandleContent
		}
		playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
		if playerPlayData.Voted {
			loghlp.Errorf("player(%d) opt error, ECodeRoomUndercoverVoted", roomPlayer.GetRoleID())
			return protocol.ECodeRoomUndercoverVoted, rep, iframe.EHandleContent
		}

		playerPlayData.Voted = true
		// 广播投票
		pushMsg := &pbclient.ECMsgGamePushUndercoverVoteNotify{
			RoleID:       roomPlayer.GetRoleID(),
			TargetRoleID: req.TargetRoleID,
		}
		roomPlayer.RoomPtr.BroadCastRoomMsg(roomPlayer.GetRoleID(), protocol.ECMsgClassRoom, protocol.ECMsgGamePushUndercoverVote, pushMsg)
		roomObj.OnPlayerVote(roomPlayer.GetRoleID(), req.TargetRoleID)
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
