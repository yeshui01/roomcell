package hallserverhandler

import "roomcell/app/hallserver/ihallserver"

var (
	hallServe ihallserver.IHallServer
)

func InitHallServerObj(iserv ihallserver.IHallServer) {
	hallServe = iserv
}
