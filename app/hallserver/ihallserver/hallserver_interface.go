package ihallserver

import "roomcell/app/hallserver/hallservermain"

type IHallServer interface {
	GetHallGlobal() *hallservermain.HallServerGlobal
}
