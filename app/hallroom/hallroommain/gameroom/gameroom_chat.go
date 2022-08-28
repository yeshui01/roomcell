package gameroom

import "roomcell/pkg/sconst"

type GameRoomChat struct {
	RoomID int64
}

func (room *GameRoomChat) GetRoomType() int32 {
	return sconst.EGameRoomTypeChat
}

func (room *GameRoomChat) GetRoomID() int64 {
	return room.RoomID
}

func (room *GameRoomChat) Update(curTime int64) {

}

func NewChatRoom(roomID int64) *GameRoomChat {
	return &GameRoomChat{
		RoomID: roomID,
	}
}
