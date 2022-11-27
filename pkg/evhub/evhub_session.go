package evhub

import (
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	NetSessionTypeClient uint8 = 0 // 客户端session
	NetSessionTypeServer uint8 = 1 // 服务器session
)

const (
	SessionStatusIdle   uint8 = 0 // 空闲
	SessionStatusActive uint8 = 1 // 正常
	SessionStatusClose  uint8 = 2 // 关闭
)

type INetSession interface {
	SendMsg(netMsg INetMessage) bool
	IsClosed() bool
	GetID() int32
	GetType() uint8
}

type NetSession struct {
	sessionID         int32
	sessionType       uint8
	netConn           net.Conn
	lastAliveTime     int64
	status            uint8
	writeCh           chan *NetMessage
	headerSlice       []byte
	secondHeaderBytes []byte
	userData          interface{}
	exitWrite         bool
	reConnect         bool     // 是否需要保持重连
	connMode          int32    // 连接模式 tcp or unix
	connAddr          string   // 连接地址
	fromReconn        bool     // 来自重连的连接
	connTODO          []func() // 连接后的行为
}

type SessionMsg struct {
	session *NetSession
	msg     *NetMessage
}

func NewNetSession(t uint8, sessID int32) *NetSession {
	return &NetSession{
		sessionType:       t,
		netConn:           nil,
		lastAliveTime:     time.Now().Unix(),
		sessionID:         sessID,
		status:            SessionStatusIdle,
		headerSlice:       make([]byte, MsgHeadSize),
		secondHeaderBytes: make([]byte, SecondHeaderSize),
		writeCh:           make(chan *NetMessage),
		userData:          nil,
		exitWrite:         true,  // 默认是没有开启的
		reConnect:         false, // 是否重连
		fromReconn:        false, // 是否来自重连
	}
}

func (s *NetSession) Close() {
	if s.status == SessionStatusActive {
		s.netConn.Close()
		close(s.writeCh)
		s.status = SessionStatusClose
	}
}

func (s *NetSession) IsClosed() bool {
	return s.status != SessionStatusActive
}

func (s *NetSession) run(hub *EventHub) {
	// 开启写协程
	s.startWriteWork(hub)
	// 开启读协程
	s.startReadWork(hub)
}

func (s *NetSession) startWriteWork(hub *EventHub) {
	hub.goWg.Add(1)
	s.exitWrite = false
	go func() {
		for msg := range s.writeCh {
			wData := msg.Encode()
			_, err := s.netConn.Write(wData)
			if err != nil {
				logrus.WithField("sessionID", s.sessionID).Errorf("send data to client failed, %s", err)
				// s.netConn.Close()
			}
		}
		logrus.Infof("writeloop exited sessionID=%d", s.sessionID)
		hub.PostCommand(HubCmdNetEvent, NetEventExitWrite, s)
		hub.goWg.Done()
	}()
}

func (s *NetSession) startReadWork(hub *EventHub) {
	hub.goWg.Add(1)
	go func() {
		for {
			s.netConn.SetReadDeadline(time.Now().Add(time.Second * 40))
			_, err := io.ReadFull(s.netConn, s.headerSlice)
			if err != nil {
				logrus.Errorf("sessionID:%d read failed, %s", s.sessionID, err)
				hub.PostCommand(HubCmdNetEvent, NetEventClose, s)
				break
			}

			hd := &NetMsgHead{}
			hd.Decode(s.headerSlice)
			var hd2 *NetMsgSecondHead = nil
			if hd.HasSecond > 0 {
				_, err := io.ReadFull(s.netConn, s.secondHeaderBytes)
				if err != nil {
					logrus.Errorf("sessionID:%d read failed2, %s,len:%d,msg(%d_%d)", s.sessionID, err.Error(), hd.Len, hd.MsgClass, hd.MsgType)
					hub.PostCommand(HubCmdNetEvent, NetEventClose, s)
					break
				}
				hd2 = &NetMsgSecondHead{}
				hd2.Decode(s.secondHeaderBytes)
			}

			messageBody := make([]byte, hd.Len-MsgHeadSize)
			_, err = io.ReadFull(s.netConn, messageBody)
			if err != nil {
				logrus.WithField("sessionID", s.sessionID).Errorf("read msgbody failed, %s,len:%d,msg(%d_%d)", err, hd.Len, hd.MsgClass, hd.MsgType)
				hub.PostCommand(HubCmdNetEvent, NetEventClose, s)
				break
			}
			netMsg := NewNetMessage(hd, messageBody)
			netMsg.SecondHead = hd2
			netMsg.RecvTime = time.Now().UnixNano() / 1000000
			sessionMsg := &SessionMsg{
				session: s,
				msg:     netMsg,
			}
			hub.PostCommand(HubCmdNetEvent, NetEventMessage, sessionMsg)
		}
		hub.goWg.Done()
	}()
}

func (s *NetSession) GetUserData() interface{} {
	return s.userData
}

func (s *NetSession) DetachResource() {
	s.netConn = nil
	s.writeCh = nil
	s.userData = nil
}

func (s *NetSession) Send(msg *NetMessage) {
	if !s.IsClosed() {
		s.writeCh <- msg
	} else {
		logrus.Errorf("session is closed,can not write msg,sessionID:%d", s.sessionID)
	}
}

func (s *NetSession) RemoteAddr() string {
	if s.IsClosed() {
		return "err addr, because session is closed"
	}
	if s.netConn != nil {
		return s.netConn.RemoteAddr().String()
	}
	return "unkown addr"
}

func (s *NetSession) GetSessionType() uint8 {
	return s.sessionType
}

func (s *NetSession) GetSessionID() int32 {
	return s.sessionID
}
func (s *NetSession) IsFullClosed() bool {
	return (s.status == SessionStatusClose) && s.exitWrite
}
