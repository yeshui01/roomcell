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
	// 谁是卧底-发言
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameUndercoverTalk,
		hallroomhandler.HandlePlayerChatUndertalk,
	)
	// 数字炸弹-猜数字
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameNumberBombGuess,
		hallroomhandler.HandlePlayerGuessNumberTalk,
	)
	// 数字炸弹-设定轮数
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameNumberBombSetting,
		hallroomhandler.HandlePlayerECMsgGameNumberBombSetting,
	)
	// 数字炸弹-选择惩罚
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameNumberBombChoosePunishment,
		hallroomhandler.HandlePlayerECMsgGameNumberBombChoosePunishment,
	)

	// 拯救玩家-设定
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRescueSetting,
		hallroomhandler.HandlePlayerECMsgGameRescueSetting,
	)
	// 拯救玩家-收礼物
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRescueRecvGift,
		hallroomhandler.HandlePlayerECMsgGameRescueRecvGift,
	)
	// 拯救玩家-更改hp
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRescueChangeHp,
		hallroomhandler.HandlePlayerECMsgGameRescueChangeHp,
	)

	// 热血奔跑-发射炸弹
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRunningSendItem,
		hallroomhandler.HandlePlayerECMsgGameRunningSendItem,
	)
	// 热血奔跑-到达终点
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRunningReachEnd,
		hallroomhandler.HandlePlayerECMsgGameRunningReachEnd,
	)
	// 热血奔跑-游戏设定
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassGame,
		protocol.ECMsgGameRunningSetting,
		hallroomhandler.HandlePlayerECMsgGameRunningSetting,
	)

}
