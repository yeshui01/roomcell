package pbcmd

import (
	"roomcell/pkg/evhub"

	"github.com/gorilla/websocket"
)

type CmdTypeWebsocketMessageData struct {
	WsConn     *websocket.Conn
	WsMsgType  int32
	MsgData    []byte
	HubMsg     *evhub.NetMessage
	RecvTimeMs int64
}
