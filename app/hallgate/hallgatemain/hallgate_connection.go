package hallgatemain

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbcmd"
	"roomcell/pkg/protocol"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type HGateConnction struct {
	UserID        int64
	UserInfo      *HGateUser
	WSConn        *websocket.Conn
	LastHeartTime int64 // 最近心跳时间
	sendMsgChan   chan *evhub.NetMessage
	Closed        bool
}

func (hgc *HGateConnction) SendMsg(msg *evhub.NetMessage) bool {
	if hgc.Closed {
		return false
	}
	hgc.sendMsgChan <- msg
	return true
}
func (hgc *HGateConnction) runRead() {
	for {
		hgc.WSConn.SetReadDeadline(time.Now().Add(time.Second * 300))
		msgType, msgData, err := hgc.WSConn.ReadMessage()
		if err != nil {
			var errInfo string = err.Error()
			if strings.Index(errInfo, "closed by the remote host") != -1 {
				loghlp.Warnf("remote client(%s)[%d] closed!!!", hgc.WSConn.RemoteAddr(), hgc.UserID)
			} else {
				loghlp.Errorf("wsRecv error:%s", err.Error())
			}
			break
		}

		msgCmdData := &pbcmd.CmdTypeWebsocketMessageData{
			WsConn:     hgc.WSConn,
			WsMsgType:  int32(msgType),
			MsgData:    msgData,
			RecvTimeMs: timeutil.NowTimeMs(),
		}
		if msgType == websocket.BinaryMessage {
			// hubMsg, _ := UnzipToHubMessage(msgData)
			// hubMsg := evhub.MakeMessage(0, 0, nil)
			// hubMsg.Decode(msgData)
			hubMsg := evhub.DecodeClientMsgToServerMsg(msgData)
			msgCmdData.HubMsg = hubMsg
		}

		trframe.PostUserCommand(protocol.CellCmdClassWebsocket,
			protocol.CmdTypeWebsocketMessage,
			msgCmdData)
	}
	loghlp.Debugf("HGateConnection exitRead")
	trframe.PostUserCommand(protocol.CellCmdClassWebsocket, protocol.CmdTypeWebsocketClosed, hgc.WSConn)
}
func (hgc *HGateConnction) runWrite() {
	for sMsg := range hgc.sendMsgChan {
		if sMsg == nil {
			break
		}
		// 特殊处理测试消息
		if sMsg.Head.MsgClass == 0 && sMsg.Head.MsgType == 1 {
			// 回复文本消息
			err := hgc.WSConn.WriteMessage(websocket.TextMessage, sMsg.Data)
			if err != nil {
				loghlp.Errorf("hgateconnection send error:%s", err.Error())
			}
		} else if sMsg.Head.MsgClass > 0 && sMsg.Head.MsgType > 0 {
			// byData := sMsg.Encode()
			byData := evhub.EncodeServerMsgToClientMsg(sMsg)
			err := hgc.WSConn.WriteMessage(websocket.BinaryMessage, byData)
			if err != nil {
				loghlp.Errorf("hgateconnection send error2:%s", err.Error())
			}
			// byData, _ := ZipCompressHubMessage(sMsg)
			// err := hgc.WSConn.WriteMessage(websocket.BinaryMessage, byData)
			// if err != nil {
			// 	loghlp.Errorf("hgateconnection send error2:%s", err.Error())
			// }
		} else if sMsg.Head != nil && sMsg.SecondHead != nil && sMsg.Head.MsgClass == 0 && sMsg.Head.MsgType == 0 && sMsg.Head.HasSecond == 1 {
			// 关闭
			loghlp.Warnf("server send close message to client, role(%d)", sMsg.SecondHead.ID)
			hgc.WSConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		}
	}
	hgc.WSConn.Close()
	loghlp.Debugf("HGateConnection exitWrite")
}
func (hgc *HGateConnction) Start() {
	go hgc.runRead()
	go hgc.runWrite()
}
func (hgc *HGateConnction) Stop() {
	hgc.Closed = true
	close(hgc.sendMsgChan)
}

type HGateClientManager struct {
	connMap map[*websocket.Conn]*HGateConnction
}

func NewHGateClientManager() *HGateClientManager {
	return &HGateClientManager{
		connMap: make(map[*websocket.Conn]*HGateConnction),
	}
}
func (mgr *HGateClientManager) AddConnection(wsConn *websocket.Conn) {
	hgc := &HGateConnction{
		WSConn:      wsConn,
		UserID:      0,
		sendMsgChan: make(chan *evhub.NetMessage, 1),
		Closed:      false,
	}
	mgr.connMap[wsConn] = hgc
	hgc.Start()
}
func (mgr *HGateClientManager) RemoveConnection(wsConn *websocket.Conn) {
	if _, ok := mgr.connMap[wsConn]; ok {
		delete(mgr.connMap, wsConn)
	}
}

func (mgr *HGateClientManager) GetConnection(wsConn *websocket.Conn) *HGateConnction {
	if hgc, ok := mgr.connMap[wsConn]; ok {
		return hgc
	}
	return nil
}
