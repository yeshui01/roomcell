package robothandler

import (
	"roomcell/app/cellrobot/robotcore"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
)

// 玩家进入房间
func RobotHandleECMsgRoomPushPlayerEnterNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	notifyMsg := &pbclient.ECMsgRoomPushPlayerEnterNotify{}
	robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}

// 玩家离开房间
func RobotHandleECMsgRoomPushPlayerLeaveNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	notifyMsg := &pbclient.ECMsgRoomPushPlayerLeaveNotify{}
	robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}

// 玩家房间聊天
func RobotHandleECMsgRoomPushChatNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	notifyMsg := &pbclient.ECMsgRoomPushChatNotify{}
	robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}
