package tserver

import "roomcell/pkg/loghlp"

func (serv *TServer) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("TServer::RunStepCheck")
	if serv.userStepRun[EServerRunStepCheck] != nil {
		return serv.userStepRun[EServerRunStepCheck](curTimeMs)
	}
	return true
}
func (serv *TServer) RunStepInit(curTimeMs int64) bool {
	if serv.userStepRun[EServerRunStepInit] != nil {
		return serv.userStepRun[EServerRunStepInit](curTimeMs)
	}
	return true
}
func (serv *TServer) RunStepPreRun(curTimeMs int64) bool {
	if serv.userStepRun[EServerRunStepPreRun] != nil {
		return serv.userStepRun[EServerRunStepPreRun](curTimeMs)
	}
	loghlp.Info("TServer::RunStepPreRun")
	return true
}
func (serv *TServer) RunStepNormalRun(curTimeMs int64) bool {
	if serv.userStepRun[EServerRunStepNormalRun] != nil {
		return serv.userStepRun[EServerRunStepNormalRun](curTimeMs)
	}
	return true
}
func (serv *TServer) RunStepStop(curTimeMs int64) bool {
	if serv.userStepRun[EServerRunStepStop] != nil {
		return serv.userStepRun[EServerRunStepStop](curTimeMs)
	}
	return true
}
func (serv *TServer) RunStepEnd(curTimeMs int64) bool {
	if serv.userStepRun[EServerRunStepEnd] != nil {
		return serv.userStepRun[EServerRunStepEnd](curTimeMs)
	}
	return true
}

// frame run
func (serv *TServer) frameRun(curTimeMs int64) {
	switch serv.runStep {
	case EServerRunStepCheck:
		{
			if serv.RunStepCheck(curTimeMs) {
				serv.changeNextStep()
			}
			break
		}
	case EServerRunStepInit:
		{
			if serv.RunStepInit(curTimeMs) {
				if serv.ListenFromConfig() == nil {
					serv.changeNextStep()
				} else {
					loghlp.Errorf("listen fail")
					panic("listen from config error")
				}
			}
			break
		}
	case EServerRunStepPreRun:
		{
			if serv.RunStepPreRun(curTimeMs) {
				serv.changeNextStep()
			}
			break
		}
	case EServerRunStepNormalRun:
		{
			serv.RunStepNormalRun(curTimeMs)
			break
		}
	case EServerRunStepStop:
		{
			if serv.RunStepStop(curTimeMs) {
				serv.changeNextStep()
			}
			break
		}
	case EServerRunStepEnd:
		{
			if serv.RunStepEnd(curTimeMs) {
				serv.changeNextStep()
			}
			break
		}
	}
}
