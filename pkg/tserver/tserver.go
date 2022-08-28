package tserver

import (
	"errors"

	"roomcell/pkg/appconfig"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
)

type tserver_loop_func_t func(servTimeMs int64)

type TServer struct {
	evHub          *evhub.EventHub
	runStep        tserver_step_t
	loopFuncList   []tserver_loop_func_t
	netSessionList map[int32]*TServerSession

	// 当前服务器节点
	nodeType    int32
	nodeIndex   int32
	zoneID      int32
	userStepRun []func(int64) bool
}

func NewTServer(serverID int32, nodeIdx int32, nType int32) *TServer {
	s := &TServer{
		evHub:       evhub.NewHub(),
		runStep:     EServerRunStepCheck,
		zoneID:      serverID,
		nodeType:    nType,
		nodeIndex:   nodeIdx,
		userStepRun: make([]func(int64) bool, 0),
	}
	s.Init()
	return s
}

func (serv *TServer) Init() bool {
	// 注册命令处理
	serv.evHub.OnUserHubCommand(func(hubCmd *evhub.HubCommand) {
		var servCmd TServCommand = TServCommand{
			hub:     serv.evHub,
			userCmd: hubCmd,
		}
		serv.OnUserCommand(&servCmd)
	})
	serv.evHub.AddFrameLoopFunc(func(curTimeMs int64) {
		serv.frameRun(curTimeMs)
	})
	// 注册回调函数
	serv.RegCallback()
	for i := 0; i <= int(EServerRunStepExit); i++ {
		serv.userStepRun = append(serv.userStepRun, nil)
	}
	return true
}

func (serv *TServer) Run() {
	serv.evHub.Run()
}

func (serv *TServer) Stop() {
	serv.evHub.Stop()
}

func (serv *TServer) changeNextStep() {
	serv.runStep++
	loghlp.Infof("change server_step(%d) before step:%d", serv.runStep, serv.runStep-1)
}

func (serv *TServer) AddLoopFunc(loopFunc tserver_loop_func_t) {
	serv.loopFuncList = append(serv.loopFuncList, loopFunc)
}

func (serv *TServer) Listen(listenMode int32, listenAddr string) error {
	return serv.evHub.Listen(listenMode, listenAddr)
}

func (serv *TServer) ListenFromConfig() error {
	var errListen error = nil
	var listenAddr string = ""
	switch serv.nodeType {
	case TServerNodeTypeGate:
		{
			if serv.nodeIndex >= 0 && int(serv.nodeIndex) < len(appconfig.Instance().ZoneNodeCfgs.GateCfgs) {
				listenAddr = appconfig.Instance().ZoneNodeCfgs.GateCfgs[serv.nodeIndex].ListenAddr
			} else {
				loghlp.Errorf("listenAddr node index not match config,nodeIndex:%d", serv.nodeIndex)
			}
			break
		}
	default:
		{
			errListen = errors.New("not find listen config")
			loghlp.Errorf("listen error:%s", errListen.Error())
			break
		}
	}
	if errListen != nil {
		loghlp.Errorf("listen error:%s", errListen.Error())
	} else if len(listenAddr) > 0 {
		loghlp.Infof("listen from config %s", listenAddr)
		errListen = serv.Listen(evhub.ListenModeTcp, listenAddr)
	} else {
		loghlp.Errorf("listen addr error!!!")
	}
	return errListen
}
func (serv *TServer) SetUserStepRun(step int32, f func(int64) bool) {
	serv.userStepRun[step] = f
}
