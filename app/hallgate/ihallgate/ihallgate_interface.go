package ihallgate

import (
	"roomcell/app/hallgate/hallgatemain"
	"roomcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
)

type IHallGate interface {
	GetUserManager() *hallgatemain.HGateUserManager
	GetGateConnMgr() *hallgatemain.HGateClientManager
	SendWSClientReplyMessage(okCode int32, cltRep proto.Message, env *iframe.TRRemoteMsgEnv)
}
