package tserver

import "roomcell/pkg/evhub"

type TServerSession struct {
	netSession *evhub.NetSession
	nodeType   int32
}

func NewServerSession(netSess *evhub.NetSession) *TServerSession {
	return &TServerSession{
		netSession: netSess,
	}
}

func (ss *TServerSession) GetSessionID() int32 {
	return ss.netSession.GetSessionID()
}
