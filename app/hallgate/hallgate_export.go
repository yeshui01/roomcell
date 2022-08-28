package hallgate

import "roomcell/app/hallgate/hallgatemain"

func (hg *HallGate) GetUserManager() *hallgatemain.HGateUserManager {
	return hg.UserMgr
}
