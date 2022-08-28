package hallroommgrmain

type RoomMgrGlobal struct {
	RoomInfoMgr       *RoomInfoManager
	lastSecUpdateTime int64
}

func NewRommMgr() *RoomMgrGlobal {
	g := &RoomMgrGlobal{
		RoomInfoMgr: newRoomInfoMgr(),
	}

	return g
}

func (g *RoomMgrGlobal) Update(curTimeMs int64) {
	curTime := curTimeMs / 1000

	if curTime >= g.lastSecUpdateTime {
		g.RoomInfoMgr.SecUpdate(curTime)
		g.lastSecUpdateTime = curTime
	}
}
