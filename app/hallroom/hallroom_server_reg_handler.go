package hallroom

import (
	"roomcell/app/hallroom/hallroomhandler"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
)

func (s *HallRoom) RegisterMsgHandler() {
	// 创建房间
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomCreate,
		hallroomhandler.HandleRoomCreate,
	)
	// 离开房间
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomLeave,
		hallroomhandler.HandlePlayerRoomLeave,
	)
	// 进入房间
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomEnter,
		hallroomhandler.HandlePlayerRoomEnter,
	)
	// 玩家下线
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerDisconnect,
		hallroomhandler.HandleRoomPlayerOffline,
	)
	// 玩家心跳
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		hallroomhandler.HandleRoomPlayerKeepHeart,
	)
	// 房间聊天
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassRoom,
		protocol.ECMsgRoomChat,
		hallroomhandler.HandleRoomPlayerChat,
	)
	// 你画我猜-画画
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameDrawPaint,
		hallroomhandler.HandlePlayerDrawPaint,
	)
	// 你画我猜-猜词
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameDrawGuessWords,
		hallroomhandler.HandlePlayerGuessWords,
	)
	// 你画我猜-房间设定
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameGrawSetting,
		hallroomhandler.HandleDrawGuessSetting,
	)
	// 你画我猜-选词
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameDrawSelectWords,
		hallroomhandler.HandleDrawSelectWords,
	)
	// game-准备
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameReadyOpt,
		hallroomhandler.HandleDrawReadyOpt,
	)

	// 谁是卧底-投票
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameUndercoverVote,
		hallroomhandler.HandlePlayerUndercoverVote,
	)
	// // 谁是卧底-发言
	// trframe.RegWorkMsgHandler(
	// 	protocol.ECMsgClassGame,
	// 	protocol.ECMsgGameUndercoverVote,
	// 	hallroomhandler.HandlePlayerChatUndertalk,
	// )
}
