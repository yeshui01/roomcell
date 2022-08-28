package hallroommain

import (
	"roomcell/app/hallroom/hallroommain/gameroom"
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

func (s *HallRoomGlobal) CreateRoom(roomID int64, roomType int32, roleID int64) iroom.IGameRoom {
	var newRoom iroom.IGameRoom = nil
	switch roomType {
	case sconst.EGameRoomTypeDrawGuess:
		{
			newRoom = gameroom.NewRoomDrawGuess(roomID, s)
			newRoom.SetMasterID(roleID)
			break
		}
	case sconst.EGameRoomTypeChat:
		{
			// newRoom = gameroom.NewChatRoom(roomID)
			// loghlp.Infof("create EGameRoomTypeChat success,roomid:%d", roomID)
			break
		}
	default:
		{
			loghlp.Errorf("unhandled room type:%d,create a empty room:%d, roleID:%d", roomType, roomID, roleID)
			newRoom = gameroom.NewEmptyRoom(roomID, s)
		}
	}
	if newRoom != nil {
		s.roomList[roomID] = newRoom
	}
	return newRoom
}
func (s *HallRoomGlobal) FindRoom(roomID int64) iroom.IGameRoom {
	if r, ok := s.roomList[roomID]; ok {
		return r
	}
	return nil
}

func (mgr *HallRoomGlobal) updateEmptyRoom(curTime int64) {
	for i, r := range mgr.roomList {
		if r.IsEmpty() {
			loghlp.Infof("remove empty room(%d)", r.GetRoomID())
			mgr.RemoveRoom(i)
		}
	}
}

func (mgr *HallRoomGlobal) RemoveRoom(roomID int64) {
	delete(mgr.roomList, roomID)

	req := &pbserver.ESMsgRoomAutoDeleteReq{
		RoomID: roomID,
	}
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("RemoveRoom(%d) cbsucc:%d", roomID, okCode)
	}
	env := trframe.MakeMsgEnv(0, nil)
	trframe.ForwardZoneMessage(
		protocol.ESMsgClassRoom,
		protocol.ESMsgRoomAutoDelete,
		req,
		trnode.ETRNodeTypeHallRoomMgr,
		0,
		cb,
		env,
	)
}
