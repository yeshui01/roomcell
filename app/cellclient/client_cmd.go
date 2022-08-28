package cellclient

import "roomcell/pkg/evhub"

const (
	EClientCmdCreateRoom   = 1
	EClientCmdEnterRoom    = 2
	EClientCmdRecvMessage  = 3
	EClientCmdLeaveRoom    = 4
	EClientCmdKeepHeart    = 5
	EClientCmdDrawReady    = 6
	EClientCmdSyncPainData = 7
	EClientCmdChatMessage  = 8
)

type ClientCmd struct {
	CmdType int32
	CmdData interface{}
}

// 数据
type ClientCmdCreateRoomData struct {
}

type ClientCmdRecvMessageData struct {
	recvMessage *evhub.NetMessage
}
