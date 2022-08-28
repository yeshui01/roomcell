package tserver

import (
	"roomcell/pkg/appconfig"
	"roomcell/pkg/evhub"
)

// 客户端连接
func (serv *TServer) OnClientConnect(servSession *TServerSession, userData interface{}) {
	if userData != nil {
		netUserData, err := userData.(*UserData)
		if err == false {
			if netUserData.dataType == EUserDataTypeNetInfo {
				servSession.nodeType = netUserData.NodeType
				NetNodeMgr().AddServerNode(appconfig.Instance().ServerID, netUserData.NodeType, netUserData.NodeIndex, servSession)
			}
		}
	}
	serv.netSessionList[servSession.GetSessionID()] = servSession
}

// 服务器连接
func (serv *TServer) OnSessionConnect(servSession *TServerSession) {
	serv.netSessionList[servSession.GetSessionID()] = servSession
}

// 收到网络消息
func (serv *TServer) OnNetMessage(servSession *TServerSession, msg *evhub.NetMessage) {
	// TODO
}

// 连接关闭
func (serv *TServer) OnSessionDisconnect(servSession *TServerSession) {
	ssID := servSession.GetSessionID()
	if _, ok := serv.netSessionList[ssID]; ok {
		delete(serv.netSessionList, ssID)
	}
}
func (serv *TServer) OnClientDisconnect(servSession *TServerSession) {
	ssID := servSession.GetSessionID()
	if _, ok := serv.netSessionList[ssID]; ok {
		delete(serv.netSessionList, ssID)
	}
}

// 包装一下hubcommand
type TServCommand struct {
	userCmd *evhub.HubCommand
	hub     *evhub.EventHub
}

// 用户命令
func (serv *TServer) OnUserCommand(servCmd *TServCommand) {
	// TODO
}

func (serv *TServer) RegCallback() {
	serv.evHub.OnClientConnection(func(netSession *evhub.NetSession, userData interface{}) {
		servSess := NewServerSession(netSession)
		serv.OnClientConnect(servSess, userData)
	})
	serv.evHub.OnSessionConnection(func(netSession *evhub.NetSession) {
		servSess := NewServerSession(netSession)
		serv.OnSessionConnect(servSess)
	})
	serv.evHub.OnMessage(func(netSession *evhub.NetSession, netMsg *evhub.NetMessage) {
		if s, ok := serv.netSessionList[netSession.GetSessionID()]; ok {
			serv.OnNetMessage(s, netMsg)
		}
	})
	serv.evHub.OnClientDisconnect(func(netSession *evhub.NetSession) {
		if s, ok := serv.netSessionList[netSession.GetSessionID()]; ok {
			serv.OnClientDisconnect(s)
		}
	})
	serv.evHub.OnNetDisconnect(func(netSession *evhub.NetSession) {
		if s, ok := serv.netSessionList[netSession.GetSessionID()]; ok {
			serv.OnSessionDisconnect(s)
		}
	})
}
