package halldatamain

import "roomcell/pkg/tbobj"

type HallDataPlayer struct {
	RoleID         int64
	DataTbRoleBase *tbobj.TbRoleBase
	VisitTime      int64
}

func NewHalldataPlayer() *HallDataPlayer {
	return &HallDataPlayer{}
}
