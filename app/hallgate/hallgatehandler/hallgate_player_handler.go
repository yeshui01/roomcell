package hallgatehandler

import (
	"roomcell/app/hallgate/hallgatemain"
	"roomcell/pkg/crossdef"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
)

// 登录大厅
func HandlePlayerLoginHall(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerLoginHallReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	//rep := &pbserver.ESMsgPlayerLoadRoleRep{}

	// 验证touken
	ok, tokenRes := crossdef.TokenAuthClaims(req.Token, crossdef.SignKey)
	if !ok {
		logrus.Error("token parse fail")
		return protocol.ECodeTokenExpire, nil, iframe.EHandleContent
	} else {
		loghlp.Infof("parse player token success:%+v", *tokenRes)
	}
	// hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	// if hgsession.UserID > 0 {
	// 	//
	// 	gateUser := hallGateServe.GetUserManager().GetGateUser(hgsession.UserID)
	// 	if gateUser != nil {
	// 		loghlp.Debugf("recv gateuser(%d) heart", hgsession.UserID)
	// 		hallServIndex := gateUser.GetHallNode().NodeIndex
	// 		trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
	// 			protocol.ECMsgPlayerKeepHeart,
	// 			req,
	// 			trnode.ETRNodeTypeHallServer,
	// 			hallServIndex,
	// 			hgsession.UserID,
	// 		)
	// 		if gateUser.GetRoomNode() != nil {
	// 			trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
	// 				protocol.ECMsgPlayerKeepHeart,
	// 				req,
	// 				trnode.ETRNodeTypeHallRoom,
	// 				gateUser.GetRoomNode().NodeIndex,
	// 				hgsession.UserID,
	// 			)
	// 		}
	// 	}
	// }

	loginReq := &pbserver.ESMsgPlayerLoginHallReq{
		UserID:   tokenRes.UserID,
		Account:  tokenRes.Nickname,
		DataZone: tokenRes.DataZone,
		GateInfo: &pbserver.ServerNodeInfo{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeHallGate,
			NodeIndex: trframe.GetCurNodeIndex(),
		},
	}
	// // 测试
	// loginReq := &pbserver.ESMsgPlayerLoginHallReq{
	// 	UserID:   1,
	// 	Account:  "testAcc1",
	// 	DataZone: 1,
	// 	GateInfo: &pbserver.ServerNodeInfo{
	// 		ZoneID:    trframe.GetFrameConfig().ZoneID,
	// 		NodeType:  trnode.ETRNodeTypeHallGate,
	// 		NodeIndex: trframe.GetCurNodeIndex(),
	// 	},
	// }
	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	hallServIndex := int32(0) // 暂时先默认用0
	// 发送消息
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("loginHall callback success,okCode:%d", okCode)
		cbRep := &pbserver.ESMsgPlayerLoginHallRep{}
		if !trframe.DecodePBMessage2(msgData, cbRep) {
			loghlp.Error("decode cbRep error")
			trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
			return
		}
		cltRep := &pbclient.ECMsgPlayerLoginHallRsp{
			RoleData: cbRep.RoleData,
		}
		if okCode == protocol.ECodeSuccess {
			// 关联玩家数据
			hgsession.UserID = cbRep.RoleData.RoleID
			loghlp.Infof("user_login_hall_succ, roleid:%d, ipAddr:%s",
				hgsession.UserID,
				hgsession.WSConn.RemoteAddr(),
			)
			// gate user数据初始化
			gateUser := hallgatemain.NewGateUser(cbRep.RoleData.RoleID)
			gateUser.SetHallNode(&trnode.TRNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeHallServer,
				NodeIndex: hallServIndex,
			})
			gateUser.SetGateConnect(hgsession)
			hallGateServe.GetUserManager().AddGateUser(cbRep.RoleData.RoleID, gateUser)
		} else {
			loghlp.Errorf("login hall fail!!!")
		}
		hallGateServe.SendWSClientReplyMessage(okCode, cltRep, env)
		//trframe.SendReplyMessage(protocol.ECodeSuccess, cltRep, env)
	}

	cbEnv := trframe.MakeMsgEnv(0,
		tmsgCtx.NetMessage)
	cbEnv.UserData = hgsession

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginHall,
		loginReq,
		trnode.ETRNodeTypeHallServer,
		hallServIndex,
		cb,
		cbEnv,
	)
	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}

// 心跳
func HandlePlayerHeart(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerKeepHeartReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}

	// 发送到hall更新heartTime
	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*hallgatemain.HGateConnction)
	gateUser := hallGateServe.GetUserManager().GetGateUser(hgsession.UserID)
	if gateUser != nil {
		loghlp.Debugf("recv gateuser(%d) heart", hgsession.UserID)
		hallServIndex := gateUser.GetHallNode().NodeIndex
		trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
			protocol.ECMsgPlayerKeepHeart,
			req,
			trnode.ETRNodeTypeHallServer,
			hallServIndex,
			hgsession.UserID,
		)
		if gateUser.GetRoomNode() != nil {
			trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
				protocol.ECMsgPlayerKeepHeart,
				req,
				trnode.ETRNodeTypeHallRoom,
				gateUser.GetRoomNode().NodeIndex,
				hgsession.UserID,
			)
		}
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 踢人
func HandlePlayerKickout(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerKickOutReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	loghlp.Warnf("HandlePlayerKickout:%+v", req)
	rep := &pbserver.ESMsgPlayerKickOutRep{}
	gateUser := hallGateServe.GetUserManager().GetGateUser(req.RoleID)
	if gateUser != nil {
		loghlp.Warnf("kickout player(%d), reason:%d succ!!!", gateUser.GetGateConnect().UserID, req.Reason)

		gateConnect := gateUser.GetGateConnect()
		if req.Reason == sconst.EPlayerOfflineReasonReplaceLogin {
			loghlp.Warnf("role(%d) be kickout for replace login", req.RoleID)
			// 发送一个提示消息
			tipMsg := evhub.MakeMessage(protocol.ECMsgClassPlayer,
				protocol.ECMsgPlayerPushLoginKick,
				make([]byte, 0),
			)
			gateConnect.SendMsg(tipMsg)
		}
		// 网关断掉
		emptyMsg := evhub.MakeEmptyMessage()
		emptyMsg.Head.HasSecond = 1
		emptyMsg.SecondHead = &evhub.NetMsgSecondHead{
			ID: req.RoleID,
		}
		gateConnect.SendMsg(emptyMsg)

		gateConnect.UserID = 0 // 解除关联
		// 删除用户
		hallGateServe.GetUserManager().DelGateUser(req.RoleID)
		hallGateServe.GetGateConnMgr().RemoveConnection(gateConnect.WSConn)
	} else {
		loghlp.Warnf("kickout player(%d), reason:%d, but not find gate user", req.RoleID, req.Reason)
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
