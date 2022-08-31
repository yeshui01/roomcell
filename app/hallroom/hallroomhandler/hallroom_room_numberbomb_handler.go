package hallroomhandler

import (
	"roomcell/app/hallroom/hallroommain/gameroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
)

// 猜数字发言
func HandlePlayerGuessNumberTalk(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameNumberBombGuessReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	loghlp.Infof("HandlePlayerChatUndertalk")
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameNumberBombGuessRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomNumberBomb)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep != sconst.ENumberBombStepGuessNumber {
		loghlp.Errorf("roomObj.RoomStep != sconst.ENumberBombStepGuessNumber")
		return protocol.ECodeRoomNumberbombInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if playerPlayData.IsTalked {
		loghlp.Errorf("playerPlayData.IsOut ECodeRoomNumberbombHasTalked")
		return protocol.ECodeRoomNumberbombHasTalked, rep, iframe.EHandleContent
	}
	if req.GuessNumber < roomObj.MinNumber || req.GuessNumber > roomObj.MaxNumber {
		loghlp.Errorf("req.GuessNumber(%d) < roomObj.MinNumber(%d) || req.GuessNumber > roomObj.MaxNumber(%d)",
			req.GuessNumber,
			roomObj.MinNumber,
			roomObj.MaxNumber)
		return protocol.ECodeRoomNumberbombInvalideOption, rep, iframe.EHandleContent
	}
	// 推送发言
	pushMsg := &pbclient.ECMsgGamePushNumberBombGuessNotify{
		GuessNumber: req.GuessNumber,
		Talker: &pbclient.RoomTalker{
			RoleID:   roomPlayer.GetRoleID(),
			Nickname: roomPlayer.Nickname,
			Icon:     roomPlayer.Icon,
		},
	}
	roomObj.BroadCastRoomMsg(0,
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushNumberBombGuess,
		pushMsg)

	roomObj.PlayerGuess(roomPlayer.GetRoleID(), req.GuessNumber)

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 游戏设定
func HandlePlayerECMsgGameNumberBombSetting(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameNumberBombSettingReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	loghlp.Infof("HandlePlayerECMsgGameNumberBombSetting")
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameNumberBombSettingRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomNumberBomb)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep != sconst.ENumberBombStepReady {
		loghlp.Errorf("roomObj.RoomStep != sconst.ENumberBombStepReady")
		return protocol.ECodeRoomNumberbombInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	roomObj.MaxTurn = req.MaxTurn

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 选择惩罚
func HandlePlayerECMsgGameNumberBombChoosePunishment(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameNumberBombChoosePunishmentReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	loghlp.Infof("ECMsgGameNumberBombChoosePunishment")
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameNumberBombChoosePunishmentRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomNumberBomb)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep != sconst.ENumberBombStepReady {
		loghlp.Errorf("roomObj.RoomStep != sconst.ENumberBombStepReady")
		return protocol.ECodeRoomNumberbombInvalideOption, rep, iframe.EHandleContent
	}
	playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	//roomObj.MaxTurn = req.MaxTurn

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
