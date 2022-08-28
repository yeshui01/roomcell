package gate

import (
	"roomcell/pkg/appconfig"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/tserver"
)

type GateServer struct {
	tserver.TServer
}

func NewGateServer(servIndex int32) *GateServer {
	s := &GateServer{}
	ts := tserver.NewTServer(appconfig.Instance().ServerID, servIndex, tserver.TServerNodeTypeGate)
	s.TServer = *ts
	// 设置运行函数
	s.SetUserStepRun(int32(tserver.EServerRunStepPreRun), s.RunStepPreRun)
	s.SetUserStepRun(int32(tserver.EServerRunStepNormalRun), s.RunStepNormalRun)
	return s
}

// 运行前一刻
func (gate *GateServer) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("GateServer::RunStepPreRun")
	return true
}
func (gate *GateServer) RunStepNormalRun(curTimeMs int64) bool {
	// loghlp.Info("GateServer::RunStepNormalRun")
	return true
}
