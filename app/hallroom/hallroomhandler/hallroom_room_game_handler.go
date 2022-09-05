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

// 准备
func HandleDrawReadyOpt(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgGameReadyOptReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgGameReadyOptRsp{}
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
	// 就不用接口了,直接实例判断
	switch hallPlayer.RoomPtr.GetRoomType() {
	case sconst.EGameRoomTypeDrawGuess:
		{
			roomDrawGuess, ok := hallPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
			if !ok {
				loghlp.Errorf("roomDrawGuess convert fail")
				return protocol.ECodeSysError, rep, iframe.EHandleContent
			}
			if roomDrawGuess.RoomStep != sconst.EDrawGuessStepReady {
				return protocol.ECodeRoomDrawInvalideOption, rep, iframe.EHandleContent
			}
			playData := roomDrawGuess.HoldPlayerData(hallPlayer.GetRoleID())
			playData.Ready = req.Ready
			break
		}
	case sconst.EGameRoomTypeUndercover:
		{
			roomObj, ok := hallPlayer.RoomPtr.(*gameroom.RoomUndercover)
			if !ok {
				loghlp.Errorf("RoomUndercover convert fail")
				return protocol.ECodeSysError, rep, iframe.EHandleContent
			}
			if roomObj.RoomStep != sconst.EUndercoverStepReady {
				return protocol.ECodeRoomUndercoverInvalideOption, rep, iframe.EHandleContent
			}
			playData := roomObj.HoldPlayerData(hallPlayer.GetRoleID())
			playData.Ready = req.Ready
			break
		}
	case sconst.EGameRoomTypeNumberBomb:
		{
			roomObj, ok := hallPlayer.RoomPtr.(*gameroom.RoomNumberBomb)
			if !ok {
				loghlp.Errorf("RoomNumberBomb convert fail")
				return protocol.ECodeSysError, rep, iframe.EHandleContent
			}
			if roomObj.RoomStep != sconst.EUndercoverStepReady {
				return protocol.ECodeRoomNumberbombInvalideOption, rep, iframe.EHandleContent
			}
			playData := roomObj.HoldPlayerData(hallPlayer.GetRoleID())
			playData.Ready = req.Ready
			break
		}
	case sconst.EGameRoomTypeRescuePlayer:
		{
			roomObj, ok := hallPlayer.RoomPtr.(*gameroom.RoomRescue)
			if !ok {
				loghlp.Errorf("RoomRescue convert fail")
				return protocol.ECodeSysError, rep, iframe.EHandleContent
			}
			if roomObj.RoomStep != sconst.ERescueStepReady {
				return protocol.ECodeInvalideOperation, rep, iframe.EHandleContent
			}
			playData := roomObj.HoldPlayerData(hallPlayer.GetRoleID())
			playData.Ready = req.Ready
			break
		}
	default:
		{
			loghlp.Errorf("no handled player ready opt,room type:%d", hallPlayer.RoomPtr.GetRoomType())
		}
	}

	pushReady := &pbclient.ECMsgGamePushPlayerReadyStatusNotify{
		RoleID: hallPlayer.GetRoleID(),
		Ready:  req.Ready,
	}
	hallPlayer.RoomPtr.BroadCastRoomMsg(hallPlayer.GetRoleID(), protocol.ECMsgClassGame, protocol.ECMsgGamePushPlayerReadyStatus, pushReady)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
