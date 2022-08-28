package hallgatehandler

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
)

func HandleServerNodeRegister(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	// TODO:
	return protocol.ECodeSuccess, nil, iframe.EHandleContent
}

// 推送消息给客户端
func HandleFramePushMsgToClient(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.EFrameMsgPushMsgToClientReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleNone
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	gateUser := hallGateServe.GetUserManager().GetGateUser(req.RoleID)
	if gateUser != nil {
		pushMessage := evhub.MakeMessage(req.MsgClass, req.MsgType, req.MsgData)
		gateUser.SendMessageToSelf(pushMessage)
		loghlp.Debugf("push player(%d) msg(%d_%d)", req.RoleID, req.MsgClass, req.MsgType)
	}
	return protocol.ECodeSuccess, nil, iframe.EHandleNone
}

// 广播消息给客户端
func HandleFrameBroadcastMsgToClient(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.EFrameMsgBroadcastMsgToClientReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleNone
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	for _, RoleID := range req.RoleList {
		gateUser := hallGateServe.GetUserManager().GetGateUser(RoleID)
		if gateUser != nil {
			pushMessage := evhub.MakeMessage(req.MsgClass, req.MsgType, req.MsgData)
			gateUser.SendMessageToSelf(pushMessage)
			loghlp.Debugf("broadcast push player(%d) msg(%d_%d)", RoleID, req.MsgClass, req.MsgType)
		}
	}
	return protocol.ECodeSuccess, nil, iframe.EHandleNone
}
