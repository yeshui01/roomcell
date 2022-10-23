package hallroommain

import (
	"roomcell/app/hallroom/hallroommain/iroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
)

type RoomPlayer struct {
	RoleID    int64  // 角色id
	RoomID    int64  // 所在房间
	Nickname  string // 昵称
	Icon      int32  // 图像
	RoomPtr   iroom.IGameRoom
	HallNode  *trnode.TRNodeInfo // 所在大厅节点
	GateNode  *trnode.TRNodeInfo // 所在网关节点
	HeartTime int64
	Ready     int32
}

func NewRoomPlayer(roleID int64, nickName string) *RoomPlayer {
	p := &RoomPlayer{
		RoleID:    roleID,
		Nickname:  nickName,
		HeartTime: timeutil.NowTime(),
	}
	return p
}

func (p *RoomPlayer) GetRoleID() int64 {
	return p.RoleID
}
func (p *RoomPlayer) GetName() string {
	return p.Nickname
}

func (p *RoomPlayer) SetRoomID(roomID int64) {
	p.RoomID = roomID
	if roomID == 0 {
		p.RoomPtr = nil
	}
}
func (p *RoomPlayer) SetRoomPtr(roomObj iroom.IGameRoom) {
	p.RoomPtr = roomObj
}
func (p *RoomPlayer) GetRoomID() int64 {
	return p.RoomID
}
func (p *RoomPlayer) SendToClient(msgClass int32, msgType int32, pbMsg proto.Message) {
	if p.GateNode == nil {
		loghlp.Errorf("player(%d)GateNode is nil!!!")
		return
	}

	clientMsg := &pbframe.EFrameMsgPushMsgToClientReq{
		RoleID:   p.RoleID,
		MsgType:  msgType,
		MsgClass: msgClass,
	}
	msgData, _ := proto.Marshal(pbMsg)
	clientMsg.MsgData = msgData
	trframe.PushZoneMessage(protocol.EMsgClassFrame, protocol.EFrameMsgPushMsgToClient, clientMsg, p.GateNode.NodeType, p.GateNode.NodeIndex)
}

func (p *RoomPlayer) ToClientPlayerInfo() *pbclient.RoomPlayer {
	return &pbclient.RoomPlayer{
		RoleID:   p.RoleID,
		Nickname: p.Nickname,
		Icon:     p.Icon,
		Ready:    p.Ready,
	}
}

func (p *RoomPlayer) GetIcon() int32 {
	return p.Icon
}

func (p *RoomPlayer) GetReady() int32 {
	return p.Ready
}
func (p *RoomPlayer) SetReady(r int32) {
	p.Ready = r
}
