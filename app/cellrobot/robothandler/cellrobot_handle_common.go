package robothandler

import (
	"roomcell/app/cellrobot/robotcore"
	"roomcell/pkg/evhub"
)

func RobotHandleCommonNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	// notifyMsg := &pbclient.ECMsgGamePushPlayerReadyStatusNotify{}
	// robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
	robotObj := robotIns.(*robotcore.CellRobot)
	robotObj.Infof("robot(%s)(%d) recv notify msg(%d_%d)", robotObj.RobotName, robotObj.UserID, sMsg.Head.MsgClass, sMsg.Head.MsgType)
}
