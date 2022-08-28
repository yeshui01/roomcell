package hallgatehandler

import "roomcell/app/hallgate/ihallgate"

var (
	hallGateServe ihallgate.IHallGate
)

func InitHallGateServe(serv ihallgate.IHallGate) {
	hallGateServe = serv
}
