package cellclient

import (
	"fmt"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

func (c *CellClient) HandlePlayerMessage(sMsg *evhub.NetMessage) {
	switch sMsg.Head.MsgType {
	case protocol.ECMsgPlayerLoginHall:
		{
			c.HandleLoginHall(sMsg)
			break
		}
	case protocol.ECMsgPlayerKeepHeart:
		{
			c.HandleKeepHeart(sMsg)
			break
		}
	default:
		{
			c.refreshLogContent(fmt.Sprintf("unhandled player message(%d_%d)", sMsg.Head.MsgClass, sMsg.Head.MsgType))
			break
		}
	}
}

func (c *CellClient) HandleLoginHall(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgPlayerLoginHallRsp{}
	//errParse := proto.Unmarshal(sMsg.Data, pbRep)
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgPlayerLoginHallRsp:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgPlayerLoginHallRsp fail"))
	}
	if sMsg.Head.Result == protocol.ECodeSuccess {
		c.roleData = pbRep.RoleData
	}
}

func (c *CellClient) HandleKeepHeart(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgPlayerKeepHeartRsp{}
	//errParse := proto.Unmarshal(sMsg.Data, pbRep)
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgPlayerKeepHeartRsp:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgPlayerKeepHeartRsp fail"))
	}
}
