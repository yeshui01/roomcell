package hallgate

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbcmd"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"

	"github.com/gorilla/websocket"
)

func (hg *HallGate) HandleCommand(frameCmd *trframe.TRFrameCommand) {
	hKey := GetCmdHandlerKey(frameCmd.UserCmd.GetCmdClass(), frameCmd.UserCmd.GetCmdType())
	if h, ok := hg.CmdHandlerMap[hKey]; ok {
		h(frameCmd)
	} else {
		loghlp.Errorf("not find cmd handler(%d_%d)", frameCmd.UserCmd.GetCmdClass(), frameCmd.UserCmd.GetCmdType())
	}
}

func (hg *HallGate) HandleCmdWSConnect(frameCmd *trframe.TRFrameCommand) {
	wsConn, ok := frameCmd.UserCmd.GetCmdData().(*websocket.Conn)
	if !ok {
		loghlp.Errorf("HandleCmdWSConnect data error")
	}
	hg.ConnMgr.AddConnection(wsConn)
}

func (hg *HallGate) HandleCmdWSMessage(frameCmd *trframe.TRFrameCommand) {
	wsMessage, ok := frameCmd.UserCmd.GetCmdData().(*pbcmd.CmdTypeWebsocketMessageData)
	if !ok {
		loghlp.Errorf("HandleCmdWSMessage data error")
	}
	switch wsMessage.WsMsgType {
	case websocket.TextMessage:
		{
			loghlp.Infof("recv wsTextMessage:%s", string(wsMessage.MsgData))
			sendMessage := evhub.MakeMessage(0, 1, wsMessage.MsgData) // 0_1 回复文本消息
			hgc := hg.ConnMgr.GetConnection(wsMessage.WsConn)
			if hgc != nil {
				hgc.SendMsg(sendMessage) // 测试
			}
		}
	case websocket.BinaryMessage:
		{
			loghlp.Infof("recv wsBinaryMessage,msgLen:%d", len(wsMessage.MsgData))
			// if wsMessage.HubMsg.Head.MsgClass == protocol.ECMsgClassPlayer && wsMessage.HubMsg.Head.MsgType == protocol.ECMsgPlayerLoginHall {
			// 	pbReq := &pbclient.ECMsgPlayerLoginHallReq{}
			// 	pbRep := &pbclient.ECMsgPlayerLoginHallRsp{
			// 		RoleData: &pbclient.RoleInfo{
			// 			RoleID: 1000001,
			// 			Level:  1,
			// 			Name:   "tttt",
			// 		},
			// 	}
			// 	errParse := proto.Unmarshal(wsMessage.HubMsg.Data, pbReq)
			// 	if errParse == nil {
			// 		loghlp.Debugf("recv ECMsgPlayerLoginHallRsp:%+v", pbReq)
			// 		// 回复数据
			// 		repData, _ := proto.Marshal(pbRep)
			// 		repMsg := evhub.MakeMessage(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall, repData)
			// 		repMsg.SecondHead = &evhub.NetMsgSecondHead{
			// 			ReqID: wsMessage.HubMsg.SecondHead.ReqID,
			// 		}
			// 		hgc := hg.ConnMgr.GetConnection(wsMessage.WsConn)
			// 		if hgc != nil {
			// 			hgc.SendMsg(repMsg)
			// 		}
			// 	} else {
			// 		loghlp.Errorf("decode ECMsgPlayerLoginHallRsp error:%s", errParse.Error())
			// 	}
			// }

			// 处理用户消息
			hgc := hg.ConnMgr.GetConnection(wsMessage.WsConn)
			if hgc != nil {
				if hg.IsForwardToRoom(int32(wsMessage.HubMsg.Head.MsgClass), int32(wsMessage.HubMsg.Head.MsgType)) {
					cbEnv := trframe.MakeMsgEnv(0,
						wsMessage.HubMsg)
					cbEnv.UserData = hgc
					cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
						loghlp.Infof("trans client msg callback succ,okCode:%d", okCode)
						hg.SendWSClientReplyMessage2(okCode, msgData, env)
					}
					gateUser := hg.UserMgr.GetGateUser(hgc.UserID)
					if gateUser != nil {
						if gateUser.GetRoomNode() != nil {
							trframe.ForwardZoneClientMessage(
								int32(wsMessage.HubMsg.Head.MsgClass),
								int32(wsMessage.HubMsg.Head.MsgType),
								wsMessage.HubMsg.Data,
								trnode.ETRNodeTypeHallRoom,
								gateUser.GetRoomNode().NodeIndex,
								cb,
								cbEnv,
								hgc.UserID,
							)
						} else {
							loghlp.Errorf("not find player(%d) room node", hgc.UserID)
							trframe.ForwardZoneClientMessage(
								int32(wsMessage.HubMsg.Head.MsgClass),
								int32(wsMessage.HubMsg.Head.MsgType),
								wsMessage.HubMsg.Data,
								trnode.ETRNodeTypeHallRoom,
								0,
								cb,
								cbEnv,
								hgc.UserID,
							)
						}
					} else {
						loghlp.Errorf("not find gate user:%d", hgc.UserID)
					}
				} else if !trframe.DispatchMsg(hgc, wsMessage.HubMsg, nil) {
					// 网关不处理的消息,转发到大厅处理
					cbEnv := trframe.MakeMsgEnv(0,
						wsMessage.HubMsg)
					cbEnv.UserData = hgc
					cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
						loghlp.Infof("trans client msg callback succ,okCode:%d", okCode)
						hg.SendWSClientReplyMessage2(okCode, msgData, env)
					}
					gateUser := hg.UserMgr.GetGateUser(hgc.UserID)
					if gateUser != nil {
						trframe.ForwardZoneClientMessage(
							int32(wsMessage.HubMsg.Head.MsgClass),
							int32(wsMessage.HubMsg.Head.MsgType),
							wsMessage.HubMsg.Data,
							trnode.ETRNodeTypeHallServer,
							gateUser.GetHallNode().NodeIndex,
							cb,
							cbEnv,
							hgc.UserID,
						)
					} else {
						trframe.ForwardZoneClientMessage(
							int32(wsMessage.HubMsg.Head.MsgClass),
							int32(wsMessage.HubMsg.Head.MsgType),
							wsMessage.HubMsg.Data,
							trnode.ETRNodeTypeHallServer,
							0,
							cb,
							cbEnv,
							hgc.UserID,
						)
					}
				}
			} else {
				loghlp.Error("not find hgc object!!!!")
			}

			break
			// hub message
		}
	case websocket.CloseMessage:
		{
			loghlp.Infof("recv wsCloseMessage:%s", string(wsMessage.MsgData))
		}
	case websocket.PingMessage:
		{
			loghlp.Infof("recv wsPingMessage:%s", string(wsMessage.MsgData))
		}
	case websocket.PongMessage:
		{
			loghlp.Infof("recv wsPongMessage:%s", string(wsMessage.MsgData))
		}
	default:
		{
			loghlp.Warnf("unhandled ws messageType:%d", wsMessage.WsMsgType)
		}
	}
}
func (hg *HallGate) HandleCmdWSClose(frameCmd *trframe.TRFrameCommand) {
	wsConn, ok := frameCmd.UserCmd.GetCmdData().(*websocket.Conn)
	if !ok {
		loghlp.Errorf("HandleCmdWSConnect data error")
	}
	hgc := hg.ConnMgr.GetConnection(wsConn)
	if hgc != nil {
		hgc.Stop()
	}
	loghlp.Info("HandleCmdWSClose")

	// 通知hallserver玩家断线
	if hgc.UserID > 0 {
		userHallIndex := int32(0)
		gateUser := hg.UserMgr.GetGateUser(hgc.UserID)
		if gateUser != nil {
			hallNode := gateUser.GetHallNode()
			if hallNode != nil {
				userHallIndex = hallNode.NodeIndex
			}
		}
		offlineReq := &pbserver.ESMsgPlayerDisconnectReq{
			RoleID: hgc.UserID,
			Reason: sconst.EPlayerOfflineReasonNormal,
		}
		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
			loghlp.Infof("offline notify success,okCode:%d", okCode)
		}
		trframe.ForwardZoneMessage(
			protocol.ESMsgClassPlayer,
			protocol.ESMsgPlayerDisconnect,
			offlineReq,
			trnode.ETRNodeTypeHallServer,
			userHallIndex,
			cb,
			nil,
		)
		// 删除用户
		hg.UserMgr.DelGateUser(hgc.UserID)
	}

	hg.ConnMgr.RemoveConnection(wsConn)
}

const (
	MsgFactor = 1000
)

// 是否直接转发到room
func (hg *HallGate) IsForwardToRoom(msgClass int32, msgType int32) bool {
	msgKey := msgClass*MsgFactor + msgType
	switch msgKey {
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameReadyOpt:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameDrawPaint:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameDrawGuessWords:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameGrawSetting:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameDrawSelectWords:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameUndercoverTalk:
		return true
	case protocol.ECMsgClassRoom*MsgFactor + protocol.ECMsgRoomChat:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameUndercoverVote:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameNumberBombGuess:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameNumberBombChoosePunishment:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameNumberBombSetting:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRescueSetting:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRescueRecvGift:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRescueChangeHp:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRunningReachEnd:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRunningSendItem:
		return true
	case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameRunningSetting:
		return true
	}
	return false
}
