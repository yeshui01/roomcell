package cellclient

import (
	"fmt"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

func (c *CellClient) HandleGameMessage(sMsg *evhub.NetMessage) {
	switch sMsg.Head.MsgType {
	case protocol.ECMsgGameReadyOpt:
		{
			break
		}
	case protocol.ECMsgGameDrawPaint:
		{
			break
		}
	case protocol.ECMsgGamePushDrawPaint:
		{
			c.HandleDrawPainPush(sMsg)

			break
		}
	default:
		{
			c.refreshLogContent(fmt.Sprintf("unhandled game message(%d_%d)", sMsg.Head.MsgClass, sMsg.Head.MsgType))
			break
		}
	}
}
func (c *CellClient) HandleDrawPainPush(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgGamePushDrawPaintNotify{}
	if !trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomCreateRsp fail"))
		return
	}
	c.refreshLogContent(fmt.Sprintf("recv ECMsgGamePushDrawPaintNotify:%+v", pbRep))
	// 更新
	gameRoomObj := c.gameRoom.(*RoomDrawGuess)
	gameRoomObj.drawBoard.DrawByPaindata(pbRep.CurPain)
}
