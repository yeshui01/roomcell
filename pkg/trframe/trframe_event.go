package trframe

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// 客户端连接
func (tf *TRFrame) onClientConnect(frameSession *FrameSession, userData interface{}) {
	if userData != nil {
		netUserData, ok := userData.(*iframe.SessionUserData)
		if ok {
			if netUserData.DataType == iframe.ESessionDataTypeNetInfo {
				frameSession.nodeType = netUserData.NodeType
				frameSession.isServerClient = netUserData.IsServerClient
				// 关联节点信息
				frameSession.nodeInfo = &trnode.TRNodeInfo{
					ZoneID:    tf.frameConfig.ZoneID,
					NodeType:  netUserData.NodeType,
					NodeIndex: netUserData.NodeIndex,
					DesInfo:   netUserData.DesInfo,
				}
				tf.frameNodeMgr.AddNode(frameSession)
				loghlp.Infof("onClientConnect,addNode:%+v", netUserData)
			}
		}
	}
	tf.netSessionList[frameSession.GetSessionID()] = frameSession
}

// 服务器连接
func (tf *TRFrame) onSessionConnect(frameSession *FrameSession) {
	tf.netSessionList[frameSession.GetSessionID()] = frameSession
}

// 收到网络消息
func (tf *TRFrame) onNetMessage(frameSession *FrameSession, msg *evhub.NetMessage, customData interface{}) {
	if msg.Head.HasSecond > 0 {
		loghlp.Debugf("TRFrame::onNetMessage2(%d_%d),hd:%+v,hd2:%+v", msg.Head.MsgClass, msg.Head.MsgType, *(msg.Head), *(msg.SecondHead))
	} else {
		loghlp.Debugf("TRFrame::onNetMessage1(%d_%d),hd:%+v", msg.Head.MsgClass, msg.Head.MsgType, *(msg.Head))
	}

	// 检测是否为回调消息
	if tf.remoteMsgMgr.checkHandleCallbackMsg(msg) {
		return
	}
	// 如果是心跳,默认处理,直接返回
	if msg.Head.MsgClass == uint16(protocol.EMsgClassFrame) && msg.Head.MsgType == uint16(protocol.EFrameMsgKeepNodeHeart) {
		msg.SecondHead.RepID = msg.SecondHead.ReqID
		frameSession.SendMsg(msg)
		if frameSession.nodeInfo != nil {
			loghlp.Debugf("recv node heart_time msg,nodeInfo:%s", frameSession.nodeInfo.DesInfo)
		} else {
			loghlp.Debugf("recv session heart_time msg")
		}
		return
	}
	// 其他消息,分发处理
	tf.msgDispatcher.Dispatch(frameSession, msg, customData)
	if len(tf.msgDoneList) > 0 {
		for _, doJob := range tf.msgDoneList {
			doJob()
		}
		tf.msgDoneList = nil
	}
}

// 连接关闭
func (tf *TRFrame) onSessionDisconnect(frameSession *FrameSession) {
	ssID := frameSession.GetSessionID()
	if _, ok := tf.netSessionList[ssID]; ok {
		// 如果是节点,移除
		if frameSession.nodeInfo != nil {
			tf.frameNodeMgr.RemoveNode2(ssID)
		}
		delete(tf.netSessionList, ssID)
	}
}
func (tf *TRFrame) onClientDisconnect(frameSession *FrameSession) {
	loghlp.Debugf("frame.onClientDisconnect:%d", frameSession.GetSessionID())
	ssID := frameSession.GetSessionID()
	if _, ok := tf.netSessionList[ssID]; ok {
		// 如果是节点,移除
		if frameSession.nodeInfo != nil {
			tf.frameNodeMgr.RemoveNode2(ssID)
		}
		delete(tf.netSessionList, ssID)
	}
}

// 用户命令
func (tf *TRFrame) onUserCommand(servCmd *TRFrameCommand) {
	// TODO
	if tf.userCmdHandle != nil {
		tf.userCmdHandle(servCmd)
	}
}

func (tf *TRFrame) regCallback() {
	tf.evHub.OnClientConnection(func(evSession *evhub.NetSession, userData interface{}) {
		servSess := &FrameSession{
			netSession: evSession,
		}
		tf.onClientConnect(servSess, userData)
	})
	tf.evHub.OnSessionConnection(func(evSession *evhub.NetSession) {
		servSess := &FrameSession{
			netSession: evSession,
		}
		tf.onSessionConnect(servSess)
	})
	tf.evHub.OnMessage(func(evSession *evhub.NetSession, netMsg *evhub.NetMessage) {
		if s, ok := tf.netSessionList[evSession.GetSessionID()]; ok {
			tf.onNetMessage(s, netMsg, nil)
		}
	})
	tf.evHub.OnClientDisconnect(func(evSession *evhub.NetSession) {
		loghlp.Debugf("OnClientDisconnect:%d", evSession.GetSessionID())
		if s, ok := tf.netSessionList[evSession.GetSessionID()]; ok {
			tf.onClientDisconnect(s)
		}
	})
	tf.evHub.OnNetDisconnect(func(evSession *evhub.NetSession) {
		if s, ok := tf.netSessionList[evSession.GetSessionID()]; ok {
			tf.onSessionDisconnect(s)
		}
	})
	tf.evHub.OnUserHubCommand(func(hubCmd *evhub.HubCommand) {
		trcmd := TRFrameCommand{
			UserCmd: hubCmd,
			Hub:     tf.evHub,
		}
		tf.onUserCommand(&trcmd)
	})
}
