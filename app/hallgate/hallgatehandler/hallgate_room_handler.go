package hallgatehandler

import (
	"roomcell/app/hallgate/hallgatemain"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// 创建房间
func HandlePlayerCreateRoom(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgRoomCreateReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	gateUser := hallGateServe.GetUserManager().GetGateUser(hgsession.UserID)
	if gateUser == nil {
		loghlp.Errorf("not find gate user:%d", hgsession.UserID)
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	// 创建房间
	createReq := &pbserver.ESMsgRoomCreateReq{
		RoleID: hgsession.UserID,
		GateNode: &pbserver.NetNode{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeHallGate,
			NodeIndex: trframe.GetCurNodeIndex(),
		},
		RoomType: req.GameType,
	}
	hallServIndex := int32(0) // 暂时先默认用0
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("create room callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomCreateRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(env.SrcMessage, cbRep)
		cltRep := &pbclient.ECMsgRoomCreateRsp{}
		if okCode == protocol.ECodeSuccess {
			cltRep.RoomDetail = cbRep.RoomDetail
			roomNode := &trnode.TRNodeInfo{
				ZoneID:    cbRep.ZoneID,
				NodeType:  cbRep.NodeType,
				NodeIndex: cbRep.NodeIndex,
			}
			gateUser := hallGateServe.GetUserManager().GetGateUser(env.SrcMessage.SecondHead.ID)
			if gateUser != nil {
				gateUser.SetRoomNode(roomNode)
				loghlp.Infof("gateuser(%d) set room node:%+v!!!", env.SrcMessage.SecondHead.ID, roomNode)
			}
		} else {
			loghlp.Errorf("create room fail(%d)!!!", okCode)
		}
		hallGateServe.SendWSClientReplyMessage(okCode, cltRep, env)
		//trframe.SendReplyMessage(protocol.ECodeSuccess, cltRep, env)
	}

	cbEnv := trframe.MakeClientMsgEnv(hgsession.UserID,
		tmsgCtx.NetMessage, hgsession)

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		createReq,
		trnode.ETRNodeTypeHallServer,
		hallServIndex,
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}

// 离开房间
func HandlePlayerLeaveRoom(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgRoomLeaveReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	gateUser := hallGateServe.GetUserManager().GetGateUser(hgsession.UserID)
	if gateUser == nil {
		loghlp.Errorf("not find gate user:%d", hgsession.UserID)
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	// 创建房间
	createReq := &pbserver.ESMsgRoomLeaveReq{
		RoleID: hgsession.UserID,
	}
	hallServIndex := int32(0) // 暂时先默认用0
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("leave room callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomCreateRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(env.SrcMessage, cbRep)
		cltRep := &pbclient.ECMsgRoomLeaveRsp{}
		if okCode == protocol.ECodeSuccess {
			gateUser := hallGateServe.GetUserManager().GetGateUser(env.SrcMessage.SecondHead.ID)
			if gateUser != nil {
				gateUser.SetRoomNode(nil)
				loghlp.Infof("gateuser(%d) leave room succ!!!", env.SrcMessage.SecondHead.ID)
			}
		} else {
			loghlp.Errorf("leave room but fail!!!")
		}
		hallGateServe.SendWSClientReplyMessage(okCode, cltRep, env)
		//trframe.SendReplyMessage(protocol.ECodeSuccess, cltRep, env)
	}

	cbEnv := trframe.MakeClientMsgEnv(hgsession.UserID,
		tmsgCtx.NetMessage, hgsession)

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomLeave,
		createReq,
		trnode.ETRNodeTypeHallServer,
		hallServIndex,
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}

// 加入房间
func HandlePlayerJoinRoom(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgRoomEnterReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)

	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	gateUser := hallGateServe.GetUserManager().GetGateUser(hgsession.UserID)
	if gateUser == nil {
		loghlp.Errorf("not find gate user:%d", hgsession.UserID)
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	// 加入房间
	enterReq := &pbserver.ESMsgRoomEnterReq{
		RoleID: hgsession.UserID,
		GateNode: &pbserver.NetNode{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeHallGate,
			NodeIndex: trframe.GetCurNodeIndex(),
		},
		RoomID: req.RoomID,
	}
	hallServIndex := int32(0) // 暂时先默认用0
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("enter room callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomEnterRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(env.SrcMessage, cbRep)
		cltRep := &pbclient.ECMsgRoomEnterRsp{}
		if okCode == protocol.ECodeSuccess {
			cltRep.RoomDetail = cbRep.RoomDetail
			roomNode := &trnode.TRNodeInfo{
				ZoneID:    cbRep.ZoneID,
				NodeType:  cbRep.NodeType,
				NodeIndex: cbRep.NodeIndex,
			}
			gateUser := hallGateServe.GetUserManager().GetGateUser(env.SrcMessage.SecondHead.ID)
			if gateUser != nil {
				gateUser.SetRoomNode(roomNode)
				loghlp.Infof("gateuser(%d) set room node:%+v!!!", env.SrcMessage.SecondHead.ID, roomNode)
			}
		} else {
			loghlp.Errorf("enter room fail(%d)!!!", okCode)
		}
		hallGateServe.SendWSClientReplyMessage(okCode, cltRep, env)
		//trframe.SendReplyMessage(protocol.ECodeSuccess, cltRep, env)
	}

	cbEnv := trframe.MakeClientMsgEnv(hgsession.UserID,
		tmsgCtx.NetMessage, hgsession)

	trframe.ForwardZoneClientPBMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomEnter,
		enterReq,
		trnode.ETRNodeTypeHallServer,
		hallServIndex,
		cb,
		cbEnv,
		hgsession.UserID,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}
