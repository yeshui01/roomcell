package robothandler

import (
	"roomcell/app/cellrobot/robotcore"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
)

// 玩家进入房间
func RobotHandleECMsgGamePushPlayerReadyStatusNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	notifyMsg := &pbclient.ECMsgGamePushPlayerReadyStatusNotify{}
	robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}
