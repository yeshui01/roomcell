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
)

// 发送道具
func HandlePlayerECMsgGameRunningSendItem(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameRunningSendItemReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameRunningSendItemRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomRunning)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep == sconst.ERunningStepRunning {
		playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
		if playerPlayData == nil {
			return protocol.ECodeSysError, rep, iframe.EHandleContent
		}

		pushMsg := &pbclient.ECMsgGamePushRunningSendItemNotify{
			RoleID: roomPlayer.GetRoleID(),
			ItemID: req.ItemID,
		}
		roomObj.BroadCastRoomMsg(roomPlayer.GetRoleID(),
			protocol.ECMsgClassGame,
			protocol.ECMsgGamePushRunningSendItem,
			pushMsg,
		)
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 到达终点
func HandlePlayerECMsgGameRunningReachEnd(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameRunningReachEndReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameRunningReachEndRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomRunning)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep == sconst.ERunningStepRunning {
		playerPlayData := roomObj.HoldPlayerData(roomPlayer.GetRoleID())
		if playerPlayData == nil {
			return protocol.ECodeSysError, rep, iframe.EHandleContent
		}
		playerPlayData.ReachTime = timeutil.NowTime()
		playerPlayData.ReachType = gameroom.ERunningReachForNormal
		pushMsg := &pbclient.ECMsgGamePushRunningReachEndNotify{
			RoleID: roomPlayer.GetRoleID(),
		}
		roomObj.BroadCastRoomMsg(roomPlayer.GetRoleID(),
			protocol.ECMsgClassGame,
			protocol.ECMsgGamePushRunningReachEnd,
			pushMsg,
		)
		roomObj.CheckGameEnd()
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 游戏设定
func HandlePlayerECMsgGameRunningSetting(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameRunningSettingReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	rep := &pbclient.ECMsgGameRunningSettingRsp{}
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
	roomObj, ok := roomPlayer.RoomPtr.(*gameroom.RoomRunning)
	if !ok {
		loghlp.Errorf("roomObj convert fail")
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	if roomObj.RoomStep == sconst.ERunningStepReady {
		roomObj.GameTime = req.GameTime
		roomObj.Distance = req.Distance

		pushMsg := &pbclient.ECMsgGamePushRunningSettingNotify{
			Distance: req.Distance,
			GameTime: req.GameTime,
		}
		roomObj.BroadCastRoomMsg(0,
			protocol.ECMsgClassGame,
			protocol.ECMsgGamePushRunningSetting,
			pushMsg,
		)
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
