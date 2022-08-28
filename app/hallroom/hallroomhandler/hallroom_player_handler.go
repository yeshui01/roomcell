package hallroomhandler

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
)

// 玩家心跳
func HandleRoomPlayerKeepHeart(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerKeepHeartReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}
	hallRoomGlobal := roomServe.GetGlobalData()
	roomPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if roomPlayer != nil {
		loghlp.Debugf("roomplayer(%d) keep heart", roomPlayer.GetRoleID())
		roomPlayer.HeartTime = trframe.GetFrameSysNowTime()
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleNone
}
