package cellclient

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/protocol"
)

func (c *CellClient) HandleRecvMessage(sMsg *evhub.NetMessage) {
	switch sMsg.Head.MsgClass {
	case protocol.ECMsgClassPlayer:
		{
			c.HandlePlayerMessage(sMsg)
			break
		}
	case protocol.ECMsgClassRoom:
		{
			c.HandleRoomMessage(sMsg)
			break
		}
	case protocol.ECMsgClassGame:
		{
			c.HandleGameMessage(sMsg)
			break
		}
	}
}
