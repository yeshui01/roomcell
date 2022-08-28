package cellclient

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
)

func (c *CellClient) HandleCmd(doCmd *ClientCmd) {
	loghlp.Infof("HandleClientCmd:%d", doCmd.CmdType)
	switch doCmd.CmdType {
	case EClientCmdCreateRoom:
		{
			c.SendCreateRoom()
			break
		}
	case EClientCmdRecvMessage:
		{
			sMsg := doCmd.CmdData.(*evhub.NetMessage)
			c.HandleRecvMessage(sMsg)
			break
		}
	case EClientCmdLeaveRoom:
		{
			c.SendLeaveRoom()
			break
		}
	case EClientCmdEnterRoom:
		{
			roomID := doCmd.CmdData.(int64)
			c.SendEnterRoom(roomID)
			break
		}
	case EClientCmdKeepHeart:
		{
			c.SendKeepHeart()
			break
		}
	case EClientCmdDrawReady:
		{
			ready := doCmd.CmdData.(bool)
			if ready {
				c.SendDrawReady(1)
			} else {
				c.SendDrawReady(0)
			}
		}
	case EClientCmdSyncPainData:
		{
			painData := doCmd.CmdData.(*pbclient.DrawPainData)
			c.SendDrawPaint(painData)
		}
	case EClientCmdChatMessage:
		{
			talkContent := doCmd.CmdData.(string)
			c.SendChatMessage(talkContent)
			break
		}
	}

}
