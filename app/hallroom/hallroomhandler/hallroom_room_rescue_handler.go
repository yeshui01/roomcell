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

// 游戏设定
func HandlePlayerECMsgGameRescueSetting(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameRescueSettingReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	if req.MaxHp < 1 || req.MaxTime < 1 {
		return protocol.ECodeParamError, nil, iframe.EHandleContent
	}

	rep := &pbclient.ECMsgGameRescueSettingRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomRescue)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep != sconst.ERescueStepReady {
		loghlp.Errorf("roomObj.RoomStep != sconst.ERescueStepReady")
		return protocol.ECodeInvalideOperation, rep, iframe.EHandleContent
	}
	if roomPlayer.GetRoleID() != roomObj.MasterID {
		loghlp.Errorf("roomPlayer.GetRoleID() != roomObj.MasterID")
		return protocol.ECodeInvalideOperation, rep, iframe.EHandleContent
	}
	playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
	if playerPlayData == nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	roomObj.MaxTime = req.MaxTime
	roomObj.MaxHp = req.MaxHp
	pushSetting := &pbclient.ECMsgGamePushRescueSettingNotify{
		MaxHp:   req.MaxHp,
		MaxTime: req.MaxTime,
	}
	roomObj.BroadCastRoomMsg(roomPlayer.GetRoleID(),
		protocol.ECMsgClassGame,
		protocol.ECMsgGamePushRescueSetting,
		pushSetting,
	)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 收到礼物
func HandlePlayerECMsgGameRescueRecvGift(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameRescueRecvGiftReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	if req.ItemID < 1 {
		return protocol.ECodeParamError, nil, iframe.EHandleContent
	}

	rep := &pbclient.ECMsgGameRescueRecvGiftRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomRescue)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	// if roomObj.RoomStep != sconst.ERescueStepRunning {
	// 	loghlp.Errorf("roomObj.RoomStep != sconst.ERescueStepRunning")
	// 	return protocol.ECodeInvalideOperation, rep, iframe.EHandleContent
	// }
	if roomObj.RoomStep == sconst.ERescueStepRunning {
		playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
		if playerPlayData == nil {
			return protocol.ECodeSysError, rep, iframe.EHandleContent
		}

		pushMsg := &pbclient.ECMsgGamePushRescueRecvGiftNotify{
			ItemID: req.ItemID,
			RoleID: roomPlayer.GetRoleID(),
		}
		roomObj.BroadCastRoomMsg(roomPlayer.GetRoleID(),
			protocol.ECMsgClassGame,
			protocol.ECMsgGamePushRescueRecvGift,
			pushMsg,
		)
		// 道具效果判断TODO

		// 检查游戏是否结束
		roomObj.CheckGameEnd()
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
