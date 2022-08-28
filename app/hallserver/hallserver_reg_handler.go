package hallserver

import (
	"roomcell/app/hallserver/hallserverhandler"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

// 注册消息处理
func (serv *HallServer) RegisterMsgHandler() {
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginHall,
		hallserverhandler.HandlePlayerLoginHall)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerDisconnect,
		hallserverhandler.HandlePlayerDisconnect)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		hallserverhandler.HandleRoomCreate)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomLeave,
		hallserverhandler.HandleHallRoomLeave)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassRoom,
		protocol.ESMsgRoomEnter,
		hallserverhandler.HandleRoomEnter)

	trframe.RegWorkMsgHandler(protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		hallserverhandler.HandlePlayerHeartTime)

}
