package hallserverhandler

import (
	"roomcell/app/hallserver/hallservermain"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// 登录大厅
func HandlePlayerLoginHall(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerLoginHallReq{}
	rep := &pbserver.ESMsgPlayerLoginHallRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	hallGlobal := hallServe.GetHallGlobal()
	// 寻找玩家
	player := hallGlobal.FindPlayer(req.UserID)
	if player != nil {
		// 如果玩家在线
		if player.IsOnline {
			if trframe.GetFrameSysNowTime()-player.GetHeartTime() < 5 {
				// 玩家已经在线
				loghlp.Warnf("player(%d) has online, cant repeated login", req.UserID)
				return protocol.ECodeRoleHasOnline, rep, iframe.EHandleContent
			} else {
				loghlp.Warnf("player(%d) has online,kickout old player", req.UserID)
				// 踢掉之前的玩家
				hallGlobal.HandlePlayerOffline(player)
				// 通知room
				if player.RoomNode != nil {
					loghlp.Infof("player(%d) offline from roomnode:%+v", player.GetBaseData().GetRoleID(), player.RoomNode)
					offlineReq := &pbserver.ESMsgPlayerDisconnectReq{
						RoleID: player.GetBaseData().GetRoleID(),
						Reason: sconst.EPlayerOfflineReasonReplaceLogin,
					}
					cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
						loghlp.Infof("player(%d) offline notify room success,okCode:%d", player.GetBaseData().GetRoleID(), okCode)
					}
					trframe.ForwardZoneMessage(
						protocol.ESMsgClassPlayer,
						protocol.ESMsgPlayerDisconnect,
						offlineReq,
						trnode.ETRNodeTypeHallRoom,
						player.RoomNode.NodeIndex,
						cb,
						nil,
					)
					player.RoomNode = nil
				}
				// 通知gate
				if player.GateNode != nil {
					loghlp.Infof("player(%d) replace login kickout offline from gatenode:%+v", player.GetBaseData().GetRoleID(), player.GateNode)
					offlineReq := &pbserver.ESMsgPlayerKickOutReq{
						RoleID: player.GetBaseData().GetRoleID(),
						Reason: sconst.EPlayerOfflineReasonReplaceLogin,
					}
					cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
						loghlp.Infof("player(%d) replace login kickout offline notify gate success,okCode:%d", player.GetBaseData().GetRoleID(), okCode)
					}
					trframe.ForwardZoneMessage(
						protocol.ESMsgClassPlayer,
						protocol.ESMsgPlayerKickOut,
						offlineReq,
						trnode.ETRNodeTypeHallGate,
						player.GateNode.NodeIndex,
						cb,
						nil,
					)
					player.GateNode = &trnode.TRNodeInfo{
						ZoneID:    req.GateInfo.ZoneID,
						NodeType:  req.GateInfo.NodeType,
						NodeIndex: req.GateInfo.NodeIndex,
					}
				}
			}
		}

		// 直接返回
		rep.RoleData = player.ToClientRoleInfo()
		// 这里后续处理放到下一帧处理
		player.IsOnline = true
		player.UpdateHeartTime(timeutil.NowTime())
		// hallGlobal.PostJob(func() {
		// 	hallGlobal.HandlePlayerOnline(player)
		// })
		trframe.AfterMsgJob(func() {
			hallGlobal.HandlePlayerOnline(player)
		})
		return protocol.ECodeSuccess, rep, iframe.EHandleContent
	}

	reqDB := &pbserver.ESMsgPlayerLoadRoleReq{
		RoleID:     req.UserID,
		WithCreate: true,
		CreateParam: &pbserver.CreateRoleInfo{
			UserID:  req.UserID,
			Account: req.Account,
			Level:   1,
			CchId:   req.CchId,
		},
	}
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("load roledata callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgPlayerLoadRoleRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
		//hallServe
		player = hallservermain.NewHallPlayer()
		hallGlobal := hallServe.GetHallGlobal()
		player.LoadData(cbRep.RoleDetailData)
		hallGlobal.AddPlayer(player)
		player.GateNode = &trnode.TRNodeInfo{
			ZoneID:    req.GateInfo.ZoneID,
			NodeType:  req.GateInfo.NodeType,
			NodeIndex: req.GateInfo.NodeIndex,
		}
		// 这里后续处理放到下一帧处理
		player.IsOnline = true
		player.UpdateHeartTime(timeutil.NowTime())
		// hallGlobal.PostJob(func() {
		// 	player.UpdateHeartTime(timeutil.NowTime())
		// 	hallGlobal.HandlePlayerOnline(player)
		// })
		trframe.AfterMsgJob(func() {
			player.UpdateHeartTime(timeutil.NowTime())
			hallGlobal.HandlePlayerOnline(player)
		})
		cltRep := &pbserver.ESMsgPlayerLoginHallRep{
			RoleData: player.ToClientRoleInfo(),
		}
		// hallGateServe.SendWSClientReplyMessage(protocol.ECodeSuccess, cltRep, env)
		trframe.SendReplyMessage(protocol.ECodeSuccess, cltRep, env)
		return
	}
	// 这里的session是frameSession
	cbEnv := trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoadRole,
		reqDB,
		trnode.ETRNodeTypeHallData,
		0,
		cb,
		cbEnv,
	)
	return protocol.ECodeSuccess, rep, iframe.EHandlePending
}

func HandlePlayerDisconnect(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerDisconnectReq{}
	rep := &pbserver.ESMsgPlayerDisconnectRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	loghlp.Infof("player(%d) offline, reason(%d)", req.RoleID, req.Reason)
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	hallGlobal := hallServe.GetHallGlobal()
	// 寻找玩家
	player := hallGlobal.FindPlayer(req.RoleID)
	if player != nil {
		hallGlobal.HandlePlayerOffline(player)
	}

	// 通知room
	if player != nil && player.RoomNode != nil {
		loghlp.Infof("player(%d) offline from roomnode:%+v", req.RoleID, player.RoomNode)
		offlineReq := &pbserver.ESMsgPlayerDisconnectReq{
			RoleID: player.GetBaseData().GetRoleID(),
			Reason: req.Reason,
		}
		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
			loghlp.Infof("player(%d) offline notify room success,okCode:%d", req.RoleID, okCode)
		}
		trframe.ForwardZoneMessage(
			protocol.ESMsgClassPlayer,
			protocol.ESMsgPlayerDisconnect,
			offlineReq,
			trnode.ETRNodeTypeHallRoom,
			player.RoomNode.NodeIndex,
			cb,
			nil,
		)
		player.RoomNode = nil
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
func HandlePlayerHeartTime(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerKeepHeartReq{}
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	roleID := tmsgCtx.NetMessage.SecondHead.ID
	loghlp.Infof("hallplayer(%d) keep heart", roleID)
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	hallGlobal := hallServe.GetHallGlobal()
	player := hallGlobal.FindPlayer(roleID)
	if player != nil {
		player.UpdateHeartTime(timeutil.NowTime())
	} else {
		loghlp.Errorf("not find heart player(%d)", roleID)
	}
	// 这里只是更新,不需要回调
	return protocol.ECodeSuccess, rep, iframe.EHandleNone
}
