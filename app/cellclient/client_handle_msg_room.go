package cellclient

import (
	"fmt"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
)

func (c *CellClient) HandleRoomMessage(sMsg *evhub.NetMessage) {
	switch sMsg.Head.MsgType {
	case protocol.ECMsgRoomCreate:
		{
			c.HandleRoomCreate(sMsg)
			break
		}
	case protocol.ECMsgRoomLeave:
		{
			c.HandleRoomLeave(sMsg)
			break
		}
	case protocol.ECMsgRoomEnter:
		{
			c.HandleRoomEnter(sMsg)
			break
		}
	case protocol.ECMsgRoomPushPlayerEnter:
		{
			c.HandleRoomEnterNotify(sMsg)
			break
		}
	case protocol.ECMsgRoomPushPlayerLeave:
		{
			c.HandleRoomLeaveNotify(sMsg)
			break
		}
	case protocol.ECMsgRoomPushPlayerOffline:
		{
			c.HandleRoomPlayerOfflineNotify(sMsg)
			break
		}
	default:
		{
			c.refreshLogContent(fmt.Sprintf("unhandled player message(%d_%d)", sMsg.Head.MsgClass, sMsg.Head.MsgType))
			break
		}
	}
}
func (c *CellClient) HandleRoomCreate(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomCreateRsp{}
	if !trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomCreateRsp fail"))
		return
	}
	c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomCreateRsp:%+v", pbRep))
	if sMsg.Head.Result == protocol.ECodeSuccess {
		if pbRep.RoomDetail.GameType == sconst.EGameRoomTypeDrawGuess {
			drawRoom := NewRoomDrawGuess(c)
			// drawRoom.drawBoard.playerList = pbRep.GetRoomDetail().Players
			c.gameRoom = drawRoom
			for _, p := range pbRep.RoomDetail.Players {
				drawRoom.AddRoomPlayer(p)
			}
		} else if pbRep.RoomDetail.GameType == sconst.EGameRoomTypeUndercover {
			drawRoom := NewGameRoomBase(c, pbRep.RoomDetail.RoomID, "谁是卧底")
			c.gameRoom = drawRoom
			for _, p := range pbRep.RoomDetail.Players {
				drawRoom.AddRoomPlayer(p)
			}
		} else {
			c.refreshLogContent(fmt.Sprintf("unhandled room type:%d", pbRep.RoomDetail.GameType))
		}
	}
}
func (c *CellClient) HandleRoomLeave(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomLeaveRsp{}
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomLeaveRsp:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomLeaveRsp fail"))
	}
}
func (c *CellClient) HandleRoomEnter(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomEnterRsp{}
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomEnterRsp:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomEnterRsp fail"))
	}
	if sMsg.Head.Result == protocol.ECodeSuccess {
		if pbRep.RoomDetail.GameType == sconst.EGameRoomTypeDrawGuess {
			drawRoom := NewRoomDrawGuess(c)
			c.gameRoom = drawRoom
			c.gameRoom = drawRoom
			for _, p := range pbRep.RoomDetail.Players {
				drawRoom.AddRoomPlayer(p)
			}
		} else if pbRep.RoomDetail.GameType == sconst.EGameRoomTypeUndercover {
			drawRoom := NewGameRoomBase(c, pbRep.RoomDetail.RoomID, "谁是卧底")
			c.gameRoom = drawRoom
			for _, p := range pbRep.RoomDetail.Players {
				drawRoom.AddRoomPlayer(p)
			}
		} else {
			c.refreshLogContent(fmt.Sprintf("unhandled room type:%d", pbRep.RoomDetail.GameType))
		}
	}
}
func (c *CellClient) HandleRoomEnterNotify(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomPushPlayerEnterNotify{}
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.gameRoom.AddRoomPlayer(pbRep.PlayerInfo)
		c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomPushPlayerEnterNotify:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomPushPlayerEnterNotify fail"))
	}
}

func (c *CellClient) HandleRoomLeaveNotify(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomPushPlayerLeaveNotify{}
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomPushPlayerLeaveNotify:%+v", pbRep))
		c.gameRoom.DelRoomPlayer(pbRep.RoleID)
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomPushPlayerLeaveNotify fail"))
	}
}

func (c *CellClient) HandleRoomPlayerOfflineNotify(sMsg *evhub.NetMessage) {
	pbRep := &pbclient.ECMsgRoomPushPlayerOfflineNotify{}
	if trframe.DecodePBMessage(sMsg, pbRep) {
		c.refreshLogContent(fmt.Sprintf("recv ECMsgRoomPushPlayerOfflineNotify:%+v", pbRep))
	} else {
		c.refreshLogContent(fmt.Sprintf("decode ECMsgRoomPushPlayerOfflineNotify fail"))
	}
}
