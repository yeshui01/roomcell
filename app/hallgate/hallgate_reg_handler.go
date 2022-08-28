package hallgate

import (
	"roomcell/app/hallgate/hallgatehandler"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

func (hg *HallGate) RegisterMsgHandler() {
	// frame
	trframe.RegWorkMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgPushMsgToClient,
		hallgatehandler.HandleFramePushMsgToClient,
	)
	trframe.RegWorkMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgBroadcastMsgToClient,
		hallgatehandler.HandleFrameBroadcastMsgToClient,
	)
	// player
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerLoginHall,
		hallgatehandler.HandlePlayerLoginHall)

	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		hallgatehandler.HandlePlayerHeart,
	)
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerKickOut,
		hallgatehandler.HandlePlayerKickout,
	)

	// room
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassRoom,
		protocol.ECMsgRoomCreate,
		hallgatehandler.HandlePlayerCreateRoom,
	)
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassRoom,
		protocol.ECMsgRoomLeave,
		hallgatehandler.HandlePlayerLeaveRoom,
	)
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassRoom,
		protocol.ECMsgRoomEnter,
		hallgatehandler.HandlePlayerJoinRoom,
	)
}
