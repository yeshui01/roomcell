package hallroommgrhandler

import (
	//"crypto/rand"
	"math/rand"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// 创建房间
func HandleRoomCreate(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomCreateReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	session := tmsgCtx.Session.(*trframe.FrameSession)
	roomID := roomMgrServe.GetGlobalData().RoomInfoMgr.GenRoomUid()
	if roomID == 0 {
		loghlp.Errorf("gen room id fail")
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	roomMgrServe.GetGlobalData().RoomInfoMgr.PendingRoomCreate(roomID)
	roomNodeList := trframe.GetNodeListByType(trnode.ETRNodeTypeHallRoom)
	if len(roomNodeList) < 1 {
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	req.RoomID = roomID
	roomNodeIdx := rand.Intn(len(roomNodeList))
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("create room success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomCreateRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, cbRep, env)
			return
		}
		trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
		if okCode == protocol.ECodeSuccess {
			// 创建房间成功
			mgrGlobal := roomMgrServe.GetGlobalData()
			roomInfo := mgrGlobal.RoomInfoMgr.FindPendingRoomCreate(roomID)
			if roomInfo == nil {
				loghlp.Errorf("not find pending room create:%d", roomID)
				trframe.SendReplyMessage(protocol.ECodeSysError, cbRep, env)
				return
			}
			mgrGlobal.RoomInfoMgr.DeletePendingRoomCreate(roomID)
			roomInfo.RoomNode = &trnode.TRNodeInfo{
				ZoneID:    cbRep.ZoneID,
				NodeType:  cbRep.NodeType,
				NodeIndex: cbRep.NodeIndex,
			}
			mgrGlobal.RoomInfoMgr.AddRoomInfo(roomInfo)
		} else {
			loghlp.Errorf("create room fail!!!")
		}
		//hallGateServe.SendWSClientReplyMessage(okCode, cltRep, env)
		trframe.SendReplyMessage(okCode, cbRep, env)
	}

	cbEnv := trframe.MakeMsgEnv(session.GetSessionID(),
		tmsgCtx.NetMessage)

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		req,
		trnode.ETRNodeTypeHallRoom,
		int32(roomNodeIdx),
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}
func HandleRoomAutoDelete(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomAutoDeleteReq{}
	rep := &pbserver.ESMsgRoomAutoDeleteRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	roomMgrServe.GetGlobalData().RoomInfoMgr.DeleteRoomInfo(req.RoomID)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 获取房间信息
func HandleRoomFindBrief(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomFindReq{}
	rep := &pbserver.ESMsgRoomFindRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	roomInfo := roomMgrServe.GetGlobalData().RoomInfoMgr.FindRoomInfoById(req.RoomID)
	if roomInfo == nil {
		loghlp.Errorf("not find room(%d) info", req.RoomID)
		return protocol.ECodeRoomNotExisted, rep, iframe.EHandleContent
	}
	rep.RoomBrief = &pbclient.RoomSimple{
		RoomID:    req.RoomID,
		ZoneID:    roomInfo.RoomNode.ZoneID,
		NodeIndex: roomInfo.RoomNode.NodeIndex,
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
