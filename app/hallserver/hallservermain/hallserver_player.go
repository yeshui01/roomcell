package hallservermain

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/ormdef"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/tbobj"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe/trnode"
	"time"
)

type HallPlayer struct {
	RoleID       int64
	LastSaveTime int64 // 最近保存的时间
	OfflineTime  int64 // 最近离线时间
	IsOnline     bool  // 是否在线
	// 玩家DB数据,必须私有化
	tbRoleBase *tbobj.TbRoleBase
	// 网络节点
	GateNode      *trnode.TRNodeInfo
	RoomNode      *trnode.TRNodeInfo
	lastHeartTime int64 // 最近心跳时间
}

func NewHallPlayer() *HallPlayer {
	p := &HallPlayer{
		tbRoleBase:    tbobj.NewTbRoleBase(),
		lastHeartTime: 0,
	}
	return p
}

// 加载数据
func (player *HallPlayer) LoadData(pbHallRoleData *pbserver.HallRoleData) {
	for _, tbData := range pbHallRoleData.RoleTables {
		switch tbData.TableID {
		case ormdef.ETableRoleBase:
			{
				player.tbRoleBase.FromBytes(tbData.Data)
				player.RoleID = player.tbRoleBase.GetRoleID()
				break
			}
		default:
			{
				loghlp.Errorf("unknown table id:%d", tbData.TableID)
			}
		}
	}
	player.LastSaveTime = time.Now().Unix()
}

// 获取角色返回给客户端的数据
func (player *HallPlayer) ToClientRoleInfo() *pbclient.RoleInfo {
	return &pbclient.RoleInfo{
		RoleID: player.tbRoleBase.GetRoleID(),
		Level:  player.tbRoleBase.GetLevel(),
		Name:   player.tbRoleBase.GetRoleName(),
	}
}

func (player *HallPlayer) GetBaseData() *tbobj.TbRoleBase {
	return player.tbRoleBase
}

func (player *HallPlayer) Online() {
	loghlp.Infof("player(%d) online", player.RoleID)
	player.GetBaseData().SetLoginTime(timeutil.NowTime())
	player.UpdateHeartTime(timeutil.NowTime())
	player.IsOnline = true
}

func (player *HallPlayer) Offline() {
	loghlp.Infof("player(%d) offline", player.RoleID)
	player.OfflineTime = timeutil.NowTime()
	player.GetBaseData().SetOfflineTime(player.OfflineTime)
	player.IsOnline = false
}

func (player *HallPlayer) SecUpdate(curTime int64) {

}

// 获取角色返回给客户端的数据
func (player *HallPlayer) ToRoomPlayerInfo() *pbserver.RoomPlayerData {
	return &pbserver.RoomPlayerData{
		RoleID:   player.tbRoleBase.GetRoleID(),
		Nickname: player.tbRoleBase.GetRoleName(),
	}
}

func (player *HallPlayer) UpdateHeartTime(heartTime int64) {
	player.lastHeartTime = heartTime
}

func (player *HallPlayer) GetHeartTime() int64 {
	return player.lastHeartTime
}
