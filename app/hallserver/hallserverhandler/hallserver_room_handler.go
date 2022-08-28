package hallserverhandler

import (
	"roomcell/pkg/loghlp"
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
	hallGlobal := hallServe.GetHallGlobal()
	// 寻找玩家
	player := hallGlobal.FindPlayer(req.RoleID)
	if player == nil {
		// 直接返回
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	if player.RoomNode != nil {
		return protocol.ECodeRoomPlayerHasInRoom, nil, iframe.EHandleContent
	}
	req.RoleName = player.GetBaseData().GetRoleName()
	req.HallNode = &pbserver.NetNode{
		ZoneID:    trframe.GetFrameConfig().ZoneID,
		NodeType:  trnode.ETRNodeTypeHallServer,
		NodeIndex: trframe.GetCurNodeIndex(),
	}
	req.PlayerData = player.ToRoomPlayerInfo()
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("create room callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomCreateRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
		if okCode == protocol.ECodeSuccess {
			player.RoomNode = &trnode.TRNodeInfo{
				ZoneID:    cbRep.ZoneID,
				NodeType:  cbRep.NodeType,
				NodeIndex: cbRep.NodeIndex,
			}
			loghlp.Infof("hallplayer(%d) setting roomNode:%+v", req.RoleID, player.RoomNode)
		}
		trframe.SendReplyMessage(protocol.ECodeSuccess, cbRep, env)
		return
	}
	// 这里的session是frameSession
	cbEnv := trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		req,
		trnode.ETRNodeTypeHallRoomMgr,
		0,
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}

// 离开房间
func HandleHallRoomLeave(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomLeaveReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	hallGlobal := hallServe.GetHallGlobal()
	// 寻找玩家
	player := hallGlobal.FindPlayer(req.RoleID)
	if player == nil {
		// 直接返回
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	if player.RoomNode == nil {
		loghlp.Errorf("player not in room")
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("leave room callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgRoomLeaveRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
		if okCode == protocol.ECodeSuccess {
			player.RoomNode = nil
			loghlp.Infof("hall role(%d) leave room!!!", player.GetBaseData().GetRoleID())
		}
		trframe.SendReplyMessage(protocol.ECodeSuccess, cbRep, env)
		return
	}
	// 这里的session是frameSession
	cbEnv := trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomLeave,
		req,
		trnode.ETRNodeTypeHallRoom,
		player.RoomNode.NodeIndex,
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}

// 加入房间
func HandleRoomEnter(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomEnterReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	hallGlobal := hallServe.GetHallGlobal()
	// 寻找玩家
	player := hallGlobal.FindPlayer(req.RoleID)
	if player == nil {
		// 直接返回
		loghlp.Errorf("ECodeSysError not find hallplayer(%d)", req.RoleID)
		return protocol.ECodeSysError, nil, iframe.EHandleContent
	}
	if player.RoomNode != nil {
		loghlp.Errorf("ECodeRoomPlayerHasInRoom", req.RoleID)
		return protocol.ECodeRoomPlayerHasInRoom, nil, iframe.EHandleContent
	}
	req.RoleName = player.GetBaseData().GetRoleName()
	req.HallNode = &pbserver.NetNode{
		ZoneID:    trframe.GetFrameConfig().ZoneID,
		NodeType:  trnode.ETRNodeTypeHallServer,
		NodeIndex: trframe.GetCurNodeIndex(),
	}
	req.PlayerData = player.ToRoomPlayerInfo()

	// 获取房间节点
	findReq := &pbserver.ESMsgRoomFindReq{
		RoomID: req.RoomID,
	}
	cbFind := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("find room callback success,okCode:%d", okCode)
		cbFindRep := &pbserver.ESMsgRoomFindRep{}
		if !trframe.DecodePBMessage2(msgData, cbFindRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbFindRep)
		if okCode == protocol.ECodeSuccess {
			// 发送消息
			cbEnter := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("enter room callback success,okCode:%d", okCode)
				cbRep := &pbserver.ESMsgRoomEnterRep{}
				if !trframe.DecodePBMessage2(msgData, cbRep) {
					loghlp.Error("decode cbRep error")
					trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
					return
				}
				trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
				if okCode == protocol.ECodeSuccess {
					player.RoomNode = &trnode.TRNodeInfo{
						ZoneID:    cbRep.ZoneID,
						NodeType:  cbRep.NodeType,
						NodeIndex: cbRep.NodeIndex,
					}
					loghlp.Infof("hallplayer(%d) setting roomNode:%+v", req.RoleID, player.RoomNode)
				}
				trframe.SendReplyMessage(protocol.ECodeSuccess, cbRep, env)
				return
			}
			// 这里的session是frameSession
			//cbEnv := trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage)
			trframe.ForwardZoneMessage(
				protocol.ESMsgClassRoom,
				protocol.ESMsgRoomEnter,
				req,
				trnode.ETRNodeTypeHallRoom,
				cbFindRep.RoomBrief.NodeIndex,
				cbEnter,
				env, // 继续沿用
			)
		} else {
			trframe.SendReplyMessage(okCode, cbFindRep, env)
		}
	}
	// 这里的session是frameSession
	findEnv := trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomFind,
		findReq,
		trnode.ETRNodeTypeHallRoomMgr,
		0,
		cbFind,
		findEnv,
	)

	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}
