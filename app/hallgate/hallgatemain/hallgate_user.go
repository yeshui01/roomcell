package hallgatemain

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/trnode"
)

type HGateUser struct {
	hallNode    *trnode.TRNodeInfo // 大厅节点
	roomNode    *trnode.TRNodeInfo // 房间节点
	userID      int64              // 账号id
	gateConnect *HGateConnction    // 对应的连接
}

func NewGateUser(id int64) *HGateUser {
	return &HGateUser{
		userID: id,
	}
}

func (u *HGateUser) SetHallNode(hallNode *trnode.TRNodeInfo) {
	u.hallNode = hallNode
}

func (u *HGateUser) SetRoomNode(hallNode *trnode.TRNodeInfo) {
	u.roomNode = hallNode
}
func (u *HGateUser) SetGateConnect(conn *HGateConnction) {
	u.gateConnect = conn
}
func (u *HGateUser) GetGateConnect() *HGateConnction {
	return u.gateConnect
}
func (u *HGateUser) GetHallNode() *trnode.TRNodeInfo {
	return u.hallNode
}
func (u *HGateUser) GetRoomNode() *trnode.TRNodeInfo {
	return u.roomNode
}
func (u *HGateUser) SendMessageToSelf(netMessage *evhub.NetMessage) {
	if u.gateConnect != nil {
		u.gateConnect.SendMsg(netMessage)
	}
}

type HGateUserManager struct {
	userList map[int64]*HGateUser
	// gateServ *HallGate
}

func NewHGateUserManager() *HGateUserManager {
	return &HGateUserManager{
		userList: make(map[int64]*HGateUser),
		// gateServ: hgServ,
	}
}

func (mgr *HGateUserManager) AddGateUser(userID int64, gateUser *HGateUser) {
	mgr.userList[userID] = gateUser
}

func (mgr *HGateUserManager) DelGateUser(userID int64) {
	delete(mgr.userList, userID)
	loghlp.Warnf("DelGateUser:%d", userID)
}
func (mgr *HGateUserManager) GetGateUser(userID int64) *HGateUser {
	if p, ok := mgr.userList[userID]; ok {
		return p
	}
	return nil
}
