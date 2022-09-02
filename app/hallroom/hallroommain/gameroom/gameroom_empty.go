package gameroom

import (
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"

	"google.golang.org/protobuf/proto"
)

type EmptyRoom struct {
	RoomID     int64
	RoomStatus int32
	PlayerList map[int64]iroom.IGamePlayer
	RoomType   int32
	RoomGlobal iroom.IRoomGlobal
	Creator    int64 // 房间的创建人
	MasterID   int64 // 房间当前的管理人
}

func NewEmptyRoom(roomID int64, hallroomGlobal iroom.IRoomGlobal) *EmptyRoom {
	return &EmptyRoom{
		RoomID:     roomID,
		RoomStatus: sconst.ERoomStatusIdle,
		PlayerList: make(map[int64]iroom.IGamePlayer),
		RoomGlobal: hallroomGlobal,
		RoomType:   sconst.EGameRoomTypeNone,
	}
}
func (room *EmptyRoom) GetRoomType() int32 {
	return room.RoomType
}
func (room *EmptyRoom) GetRoomID() int64 {
	return room.RoomID
}

func (room *EmptyRoom) Update(curTime int64) {

}

func (room *EmptyRoom) JoinPlayer(p iroom.IGamePlayer) {
	if _, ok := room.PlayerList[p.GetRoleID()]; ok {
		return
	}
	room.PlayerList[p.GetRoleID()] = p
	p.SetRoomID(room.RoomID)
	loghlp.Infof("player(%d) join room(%d),room player num:%d", p.GetRoleID(), room.RoomID, len(room.PlayerList))
	room.BroadCastPlayerEnter(p.GetRoleID())
}

func (room *EmptyRoom) LeavePlayer(roleID int64) {
	if _, ok := room.PlayerList[roleID]; ok {
		delete(room.PlayerList, roleID)
		room.BroadCastPlayerLeave(roleID)
	}
}

func (room *EmptyRoom) ToRoomDetail() *pbclient.RoomData {
	roomData := &pbclient.RoomData{
		RoomID:   room.RoomID,
		GameType: room.RoomType,
		MasterID: room.MasterID,
	}
	// 玩家列表
	for _, p := range room.PlayerList {
		// p := &pbclient.RoomPlayer{
		// 	RoleID:   p.GetRoleID(),
		// 	Nickname: p.GetName(),
		// }
		playerInfo := p.ToClientPlayerInfo()
		roomData.Players = append(roomData.Players, playerInfo)
	}
	return roomData
}

func (room *EmptyRoom) IsEmpty() bool {
	return len(room.PlayerList) == 0
}
func (room *EmptyRoom) OnPlayerOffline(p iroom.IGamePlayer) {
	room.BroadCastPlayerOffline(p.GetRoleID())
	room.LeavePlayer(p.GetRoleID())
	loghlp.Infof("player(%d) offline on room(%d),room player num:%d", p.GetRoleID(), room.RoomID, len(room.PlayerList))

}

func (room *EmptyRoom) BroadCastPlayerEnter(newRoleID int64) {
	pushMsg := &pbclient.ECMsgRoomPushPlayerEnterNotify{
		RoleID: newRoleID,
		// PlayerInfo: p.ToClientPlayerInfo(),
	}
	pushList := make([]int64, 0)
	for k, p := range room.PlayerList {
		if k == newRoleID {
			pushMsg.PlayerInfo = p.ToClientPlayerInfo()
			continue
		}
		pushList = append(pushList, p.GetRoleID())
		loghlp.Debugf("broad cast player(%d) enter to player:%d", newRoleID, p.GetRoleID())
	}
	if len(pushList) > 0 {
		room.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerEnter, pushMsg, pushList)
	}
}
func (room *EmptyRoom) BroadCastPlayerLeave(leaveRoleID int64) {
	pushMsg := &pbclient.ECMsgRoomPushPlayerLeaveNotify{
		RoleID: leaveRoleID,
	}
	pushList := make([]int64, 0)
	for k, p := range room.PlayerList {
		if k == leaveRoleID {
			continue
		}
		pushList = append(pushList, p.GetRoleID())
		//p.SendToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerLeave, pushMsg)
		loghlp.Debugf("broad cast player(%d) leave to player:%d", leaveRoleID, p.GetRoleID())
	}
	if len(pushList) > 0 {
		room.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerLeave, pushMsg, pushList)
	}
}
func (room *EmptyRoom) BroadCastPlayerOffline(newRoleID int64) {
	pushMsg := &pbclient.ECMsgRoomPushPlayerOfflineNotify{
		RoleID: newRoleID,
	}
	pushList := make([]int64, 0)
	for k, p := range room.PlayerList {
		if k == newRoleID {
			continue
		}
		pushList = append(pushList, p.GetRoleID())
		//p.SendToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerOffline, pushMsg)
		loghlp.Debugf("broad cast player(%d) offline to player:%d", newRoleID, p.GetRoleID())
	}
	if len(pushList) > 0 {
		room.RoomGlobal.BroadcastMsgToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerOffline, pushMsg, pushList)
	}
}
func (room *EmptyRoom) BroadCastRoomMsg(excludeRoleID int64, msgClass int32, msgType int32, pbMsg proto.Message) {
	pushList := make([]int64, 0)
	for k, p := range room.PlayerList {
		if k == excludeRoleID {
			continue
		}
		pushList = append(pushList, p.GetRoleID())
		//p.SendToClient(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerOffline, pushMsg)
		loghlp.Debugf("broad cast msg(%d_%d) to player:%d", msgClass, msgType, p.GetRoleID())
	}
	if len(pushList) > 0 {
		room.RoomGlobal.BroadcastMsgToClient(msgClass, msgType, pbMsg, pushList)
	}
}

func (room *EmptyRoom) IsCanJoin() bool {
	return true
}
func (room *EmptyRoom) IsPlayerFull() bool {
	return len(room.PlayerList) >= sconst.MaxRoomPlayerNum
}
func (room *EmptyRoom) SetMasterID(masterID int64) {
	room.MasterID = masterID
}

func (room *EmptyRoom) GetMasterID() int64 {
	return room.MasterID
}
