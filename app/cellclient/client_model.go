package cellclient

import "roomcell/app/account/accrouter"

type AccountLoginRsp struct {
	Code int32                      `json:"code"`
	Msg  string                     `json:"msg"`
	Data *accrouter.AccountLoginRsp `json:"data"`
}
