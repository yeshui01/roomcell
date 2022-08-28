package evhub

import (
	"net"

	"github.com/sirupsen/logrus"
)

func (hub *EventHub) handleNetEvent(netCmd *HubCommand) {
	switch netCmd.cmdType {
	case NetEventConnected:
		{
			sess := hub.AllocateSession(NetSessionTypeServer)
			if sess == nil {
				logrus.Error("sess is nil")
				return
			}
			sess.netConn = netCmd.cmdData.(net.Conn)
			sess.run(hub)
			if sess.sessionType == NetSessionTypeServer {
				hub.onSessionComeIn(sess)
			} else {
				hub.onClientComeIn(sess)
			}
		}
	case NetEventClose:
		{
			hub.onSessionClose(netCmd.cmdData.(*NetSession))
		}
	case NetEventExitWrite:
		{
			hub.onSessionExitWrite(netCmd.cmdData.(*NetSession))
		}
	case NetEventMessage:
		{
			sessionMsg := netCmd.cmdData.(*SessionMsg)
			hub.onSessionMessage(sessionMsg.session, sessionMsg.msg)
		}
	}
}

func (hub *EventHub) onSessionComeIn(session *NetSession) {
	if session != nil {
		session.status = SessionStatusActive
		hub.sessionMap[session.sessionID] = session
		logrus.Infof("session come in, sessionID:%d,addr:%s", session.sessionID, session.netConn.RemoteAddr())
		if hub.onSessionConnect != nil {
			hub.onSessionConnect(session)
		}
	}
}
func (hub *EventHub) onClientComeIn(session *NetSession) {
	if session != nil {
		session.status = SessionStatusActive
		hub.sessionMap[session.sessionID] = session
		delete(hub.reconnectSession, session.sessionID)
		logrus.Infof("client come in, sessionID:%d,addr:%s", session.sessionID, session.netConn.RemoteAddr())
	}
}
func (hub *EventHub) onSessionClose(session *NetSession) {
	if session != nil {
		logrus.Infof("session close,sessionID:%d", session.sessionID)
		if session.sessionType != NetSessionTypeClient {
			if hub.onNetClose != nil {
				hub.onNetClose(session)
			}
		} else {
			if hub.onClientClose != nil {
				hub.onClientClose(session)
			}
		}
		session.Close()
		delete(hub.sessionMap, session.sessionID)
		if !session.IsFullClosed() {
			hub.waitCloseSessions[session.sessionID] = session
		} else {
			delete(hub.waitCloseSessions, session.sessionID)
			if session.reConnect {
				hub.reconnectSession[session.sessionID] = session
			} else {
				hub.RecycleSession(session)
			}
		}
	}
}
func (hub *EventHub) onSessionMessage(session *NetSession, msg *NetMessage) {
	if session != nil {
		if hub.onMessage != nil {
			hub.onMessage(session, msg)
		}
	}
}

func (hub *EventHub) onSessionExitWrite(session *NetSession) {
	if session != nil {
		logrus.Infof("session exit write,sessionID:%d", session.sessionID)
		if !session.exitWrite {
			session.exitWrite = true
			session.writeCh = nil
		}
		delete(hub.sessionMap, session.sessionID)
		if !session.IsFullClosed() {
			hub.waitCloseSessions[session.sessionID] = session
		} else {
			delete(hub.waitCloseSessions, session.sessionID)
			if session.reConnect {
				hub.reconnectSession[session.sessionID] = session
			} else {
				hub.RecycleSession(session)
			}
		}
	}
}
