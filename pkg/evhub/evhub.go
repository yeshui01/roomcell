/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-03-02 15:00
 * @LastEditTime: 2022-03-02 15:00
 */
package evhub

import (
	"container/list"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type CloseCallback func(netSession *NetSession)
type ClientConnectCallBack func(netSession *NetSession, userData interface{})
type ServerConnectCallBack func(netSession *NetSession)
type ErrorCallBack func(code int32, err error)
type NetMsgCallBack func(netSession *NetSession, netMsg *NetMessage)
type UserHubCommandHandle func(hubCmd *HubCommand)

const (
	ServerSessionWriteChanSize = 100
	HubCmdChanSize             = 1024
)

const (
	HubCodeConnectError    = 1000 // 连接错误
	HubCodeSessionAllocate = 1001 // session分配失败
)

const (
	ListenModeTcp  = 0 // tcp
	ListenModeUnix = 1 // unix
)

type atomicBool int32
type EvhubFrameLoopFunc func(curTimeMs int64)

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }
func (b *atomicBool) setTrue()    { atomic.StoreInt32((*int32)(b), 1) }
func (b *atomicBool) setFalse()   { atomic.StoreInt32((*int32)(b), 0) }

type EventHub struct {
	sessionMap           map[int32]*NetSession
	sessionPool          list.List
	onNetClose           CloseCallback // 服务器端的连接断开
	onClientClose        CloseCallback // 客户端连接断开
	onClientConnect      ClientConnectCallBack
	onSessionConnect     ServerConnectCallBack
	onError              ErrorCallBack
	onMessage            NetMsgCallBack
	onUserHubCommandCall UserHubCommandHandle
	// userCommandHandle
	// tcp监听
	tcpListener net.Listener

	// 命令队列
	cmdChain chan *HubCommand
	// 停止
	stopChain         chan bool
	goWg              sync.WaitGroup // 记录开启的线程,防止优雅退出时发生线程逃逸
	muSessionPool     sync.Mutex
	isStoping         atomicBool
	sessionIdx        int32
	stopingOpt        int32
	waitCloseSessions map[int32]*NetSession
	reconnectSession  map[int32]*NetSession
	reconnLastTime    int64
	frameFuncs        []EvhubFrameLoopFunc
	frameMs           int32
}

func (hub *EventHub) Listen(listenType int32, addr string) error {
	if ListenModeTcp == listenType {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		hub.tcpListener = ln
		hub.runTcpLiten()
	}

	return nil
}

func (hub *EventHub) runTcpLiten() {
	hub.goWg.Add(1)
	go func(h *EventHub) {
		var tempDelay time.Duration
		for {
			conn, err := h.tcpListener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					// log.Release("accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				break
			}
			hub.PostCommand(HubCmdNetEvent, NetEventConnected, conn)
		}
		hub.goWg.Done()
	}(hub)
}

func (hub *EventHub) postCommandInner(netCmd *HubCommand) {
	if netCmd != nil {
		hub.cmdChain <- netCmd
	}
}
func (hub *EventHub) PostCommand(cmdClass int32, cmdType int32, cmdData interface{}) {
	netCmd := NewHubCommand(cmdClass, cmdType, cmdData)
	hub.postCommandInner(netCmd)
}

// 主循环
func (hub *EventHub) Run() {
	sectick := time.NewTicker(time.Second)
	watchTick := time.NewTicker(time.Second * 60)
	frameTick := time.NewTicker(time.Millisecond * time.Duration(hub.frameMs))
	for {
		select {
		case netCmd := <-hub.cmdChain:
			hub.handleHubCmd(netCmd)
		case <-hub.stopChain:
			logrus.Info("stop hub")
			hub.shutDown()
		case <-watchTick.C:
			logrus.Info("hub watch tick")
			hub.watchReport()
		case <-frameTick.C:
			curTimeMs := time.Now().UnixNano() / 1000000
			for _, f := range hub.frameFuncs {
				f(curTimeMs)
			}
		case <-sectick.C:
			// 检测是都关闭完了
			if hub.stopingOpt > 0 {
				if hub.isStopFinish() {
					logrus.Info("hub stop! exit hub!")
					time.Sleep(time.Second * 3)
					hub.goWg.Wait() // 等待所有开启的线程都退出
					hub.watchReport()
					sectick.Stop()
					watchTick.Stop()
					return
				} else {
					logrus.Info("wait hub stoping")
				}
			} else {
				hub.runReconnect(time.Now().Unix())
			}
		}
	}
}

// 停止
func (hub *EventHub) Stop() {
	if hub.isStoping.isSet() {
		return
	}
	hub.stopChain <- true
	hub.isStoping.setTrue()
}
func (hub *EventHub) handleHubCmd(netCmd *HubCommand) {
	if netCmd.cmdClass > HubCmdUserBase {
		if hub.onUserHubCommandCall != nil {
			hub.onUserHubCommandCall(netCmd)
		} else {
			logrus.Errorf("not handled user hubcommand(%d,%d)", netCmd.cmdClass, netCmd.cmdType)
			return
		}
		return
	}
	// 内部hub command
	switch netCmd.cmdClass {
	case HubCmdNetEvent:
		hub.handleNetEvent(netCmd)
	}
}
func (hub *EventHub) RecycleSession(session *NetSession) {
	hub.muSessionPool.Lock()
	defer hub.muSessionPool.Unlock()
	logrus.Debugf("RecycleSession:%d", session.GetSessionID())
	session.status = SessionStatusIdle
	session.DetachResource()
	hub.sessionPool.PushBack(session)
}
func (hub *EventHub) AllocateSession(sessionType uint8) *NetSession {
	hub.muSessionPool.Lock()
	defer hub.muSessionPool.Unlock()
	if hub.sessionPool.Len() > 0 {
		s := hub.sessionPool.Front()
		hub.sessionPool.Remove(s)
		session := s.Value.(*NetSession)
		session.sessionType = sessionType
		if sessionType == NetSessionTypeServer {
			session.writeCh = make(chan *NetMessage, ServerSessionWriteChanSize)
		} else {
			session.writeCh = make(chan *NetMessage)
		}

		return session
	}
	hub.sessionIdx++
	if hub.sessionIdx >= 10000000 {
		logrus.Error("sessionIdx is full")
		return nil
	}
	s := NewNetSession(sessionType, hub.sessionIdx)
	if sessionType == NetSessionTypeServer {
		s.writeCh = make(chan *NetMessage, ServerSessionWriteChanSize)
	}
	return s
}

// 连接服务器
func (hub *EventHub) Connect(listenType int32, addr string, keepReConn bool, userData interface{}) bool {
	if listenType == ListenModeTcp {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			logrus.Error("connect server error:" + err.Error())
			if hub.onError != nil {
				hub.onError(HubCodeConnectError, err)
			}
			return false
		}
		session := hub.AllocateSession(NetSessionTypeClient)
		if session == nil {
			if hub.onError != nil {
				hub.onError(HubCodeSessionAllocate, errors.New("allocate session fail"))
			}
			return false
		}
		session.status = SessionStatusActive
		session.netConn = conn
		session.userData = userData
		session.connAddr = addr
		session.connMode = listenType
		session.reConnect = keepReConn
		session.run(hub)
		hub.onClientComeIn(session)
		if hub.onClientConnect != nil {
			hub.onClientConnect(session, userData)
		}
		logrus.Infof("connect server %s succ!!!", addr)
	} else {
		logrus.Errorf("param error, unhandled listen type:%d", listenType)
		return false
	}
	return true
}

// 连接服务器
func (hub *EventHub) Connect2(listenType int32, addr string, keepReConn bool, userData interface{}, connDo func()) bool {
	if listenType == ListenModeTcp {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			logrus.Error("connect server error:" + err.Error())
			if hub.onError != nil {
				hub.onError(HubCodeConnectError, err)
			}
			return false
		}
		session := hub.AllocateSession(NetSessionTypeClient)
		if session == nil {
			if hub.onError != nil {
				hub.onError(HubCodeSessionAllocate, errors.New("allocate session fail"))
			}
			return false
		}
		session.status = SessionStatusActive
		session.netConn = conn
		session.userData = userData
		session.connAddr = addr
		session.connMode = listenType
		session.reConnect = keepReConn
		session.connTODO = append(session.connTODO, connDo)
		session.run(hub)
		hub.onClientComeIn(session)
		if hub.onClientConnect != nil {
			hub.onClientConnect(session, userData)
		}
		logrus.Infof("connect server %s succ!!!", addr)
		connDo()
	} else {
		logrus.Errorf("param error, unhandled listen type:%d", listenType)
		return false
	}
	return true
}
func (hub *EventHub) GetSession(sessionID int32) *NetSession {
	if s, ok := hub.sessionMap[sessionID]; ok {
		return s
	}
	return nil
}

// 客户端连接
func (hub *EventHub) OnClientConnection(callback ClientConnectCallBack) {
	hub.onClientConnect = callback
}

// 服务器连接
func (hub *EventHub) OnSessionConnection(callback ServerConnectCallBack) {
	hub.onSessionConnect = callback
}

// 收到网络消息
func (hub *EventHub) OnMessage(callback NetMsgCallBack) {
	hub.onMessage = callback
}

// 错误
func (hub *EventHub) OnError(callback ErrorCallBack) {
	hub.onError = callback
}

// 关闭处理(服务端)
func (hub *EventHub) OnNetDisconnect(callback CloseCallback) {
	hub.onNetClose = callback
}
func (hub *EventHub) OnClientDisconnect(callback CloseCallback) {
	hub.onClientClose = callback
}

// 用户命令
func (hub *EventHub) OnUserHubCommand(callback UserHubCommandHandle) {
	hub.onUserHubCommandCall = callback
}

// 停服处理
func (hub *EventHub) shutDown() {
	// 管理比
	hub.stopingOpt = 1
}

// 是否停服完毕
func (hub *EventHub) isStopFinish() bool {
	if hub.tcpListener != nil {
		hub.tcpListener.Close()
		hub.tcpListener = nil
		return false
	}
	// 连接是否都关闭
	if len(hub.sessionMap) > 0 {
		for _, session := range hub.sessionMap {
			session.Close()
		}
		return false
	}
	return true
}

func NewHub() *EventHub {
	return &EventHub{
		sessionMap:        make(map[int32]*NetSession),
		stopChain:         make(chan bool),
		sessionIdx:        1000,
		cmdChain:          make(chan *HubCommand, HubCmdChanSize),
		stopingOpt:        0,
		waitCloseSessions: make(map[int32]*NetSession),
		reconnectSession:  make(map[int32]*NetSession),
		reconnLastTime:    0,
		frameFuncs:        make([]EvhubFrameLoopFunc, 0),
		frameMs:           100,
	}
}

func (hub *EventHub) runReconnect(curTime int64) {
	if curTime-hub.reconnLastTime < 5 {
		return
	}
	// 尝试重连
	for _, session := range hub.reconnectSession {
		logrus.Debugf("session_reconnect,sessionid:%d", session.GetSessionID())
		if session.connMode == ListenModeTcp {
			conn, err := net.Dial("tcp", session.connAddr)
			if err != nil {
				logrus.Error("reconnect server error:" + err.Error())
				if hub.onError != nil {
					hub.onError(HubCodeConnectError, err)
				}
				continue
			}
			session.status = SessionStatusActive
			session.netConn = conn
			session.writeCh = make(chan *NetMessage)
			session.fromReconn = true
			session.run(hub)
			hub.onClientComeIn(session)
			if hub.onClientConnect != nil {
				hub.onClientConnect(session, session.userData)
			}
			logrus.Infof("reconnect server %s succ!!!", session.connAddr)
			for _, do := range session.connTODO {
				do()
			}
		} else {
			logrus.Errorf("reconnect param error, unhandled listen type:%d", session.connMode)
		}
	}

	hub.reconnLastTime = curTime
}

func (hub *EventHub) watchReport() {
	logrus.Infof("hubWatchReport,sessionMap.size:%d,sessionPool.size:%d,waitCloseSessions.size:%d,reconnectSession.size:%d", len(hub.sessionMap), hub.sessionPool.Len(), len(hub.waitCloseSessions), len(hub.reconnectSession))
}

func (hub *EventHub) AddFrameLoopFunc(userFunc EvhubFrameLoopFunc) {
	hub.frameFuncs = append(hub.frameFuncs, userFunc)
}
