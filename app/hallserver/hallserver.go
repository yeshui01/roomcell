package hallserver

import (
	"roomcell/app/hallserver/hallserverhandler"
	"roomcell/app/hallserver/hallservermain"
)

type HallServer struct {
	hallServerGlobal *hallservermain.HallServerGlobal
}

func NewHallServer() *HallServer {
	s := &HallServer{
		hallServerGlobal: hallservermain.NewHallServerGlobal(),
	}

	hallserverhandler.InitHallServerObj(s)
	s.RegisterMsgHandler()
	return s
}

func (serv *HallServer) FrameRun(curTimeMs int64) {
	serv.hallServerGlobal.Update(curTimeMs)
}

func (serv *HallServer) GetHallGlobal() *hallservermain.HallServerGlobal {
	return serv.hallServerGlobal
}
