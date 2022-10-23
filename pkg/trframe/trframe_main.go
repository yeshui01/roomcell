/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-10-09 11:32:38
 * @Brief:
 */
package trframe

import (
	"errors"
	"fmt"

	"roomcell/pkg/configdata"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/protocol"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/tframeconfig"
	"roomcell/pkg/trframe/tframedispatcher"
	"roomcell/pkg/trframe/trnode"
	"roomcell/pkg/trframe/trnode/tgate"
	"roomcell/pkg/trframe/trnode/thalldata"
	"roomcell/pkg/trframe/trnode/thallgate"
	"roomcell/pkg/trframe/trnode/thallroom"
	"roomcell/pkg/trframe/trnode/thallroommgr"
	"roomcell/pkg/trframe/trnode/thallserver"
	"roomcell/pkg/trframe/trnode/troot"

	"github.com/sirupsen/logrus"
)

type EFrameStep int32

const (
	ETRFrameStepCheck  EFrameStep = 0 // 运行检测
	ETRFrameStepInit   EFrameStep = 1 // 初始化
	ETRFrameStepPreRun EFrameStep = 2 // 准备运行
	ETRFrameStepRun    EFrameStep = 3 // 正常运行
	ETRFrameStepStop   EFrameStep = 4 // 停止
	ETRFrameStepEnd    EFrameStep = 5 // 结束
	ETRFrameStepExit   EFrameStep = 6 // 退出
	ETRFrameStepFinal  EFrameStep = 7 // 边界值
)

type FrameRunFunc func(curTimeMs int64)

// 包装一下hubcommand
type TRFrameCommand struct {
	UserCmd *evhub.HubCommand
	Hub     *evhub.EventHub
}
type TRFrame struct {
	frameNodeMgr   *FrameNodeMgr
	evHub          *evhub.EventHub
	runStep        EFrameStep
	loopFuncList   []FrameRunFunc
	curWorkNode    ITRFrameWorkNode // 当前工作节点
	frameConfig    *tframeconfig.FrameConfig
	nodeType       int32
	nodeIndex      int32
	netSessionList map[int32]*FrameSession
	userStepRun    []func(curTimeMs int64) bool
	userCmdHandle  func(userCmd *TRFrameCommand)
	remoteMsgMgr   *RemoteMsgCallMgr
	msgDispatcher  *tframedispatcher.FrameMsgDispatcher
	keepNodeTime   int64
	nowFrameTimeMs int64
	msgDoneList    []func() // 消息处理后的执行列表
}

func newTRFrame(configPath string, nType int32, nIndex int32) *TRFrame {
	tf := &TRFrame{
		frameNodeMgr:   NewFrameNodeMgr(),
		evHub:          evhub.NewHub(),
		runStep:        ETRFrameStepCheck,
		frameConfig:    tframeconfig.NewFrameConfig(),
		nodeType:       nType,
		nodeIndex:      nIndex,
		netSessionList: make(map[int32]*FrameSession),
		userStepRun:    make([]func(curTimeMs int64) bool, int(ETRFrameStepFinal)),
		nowFrameTimeMs: timeutil.NowTimeMs(),
	}
	// 加载配置
	err := tframeconfig.ReadFrameConfigFromFile(configPath, tf.frameConfig)
	if err != nil {
		panic(fmt.Sprintf("load frame config error:%s", err.Error()))
	}

	// 分发器
	tf.msgDispatcher = tframedispatcher.NewFrameMsgDispatcher(tf)
	// 初始化当前节点
	tf.curWorkNode = tf.makeInitWorkNode(nType, nIndex)
	if tf.curWorkNode == nil {
		panic(fmt.Sprintf("worknode(%d,%d) is nil", nType, nIndex))
	}
	tf.initNodeSetting()
	tf.remoteMsgMgr = newRemoteMsgMgr(tf)
	tf.regCallback()
	tf.evHub.AddFrameLoopFunc(func(curTimeMs int64) {
		tf.frameRun(curTimeMs)
	})
	return tf
}

func (tf *TRFrame) frameRun(curTimeMs int64) {
	tf.nowFrameTimeMs = curTimeMs
	tf.remoteMsgMgr.update(curTimeMs)
	switch tf.runStep {
	case ETRFrameStepCheck:
		{
			if tf.RunStepCheck(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepInit:
		{
			if tf.RunStepInit(curTimeMs) {
				tf.listen()
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepPreRun:
		{
			if tf.RunStepPreRun(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepRun:
		{
			tf.RunStepRun(curTimeMs)
			break
		}
	case ETRFrameStepStop:
		{
			if tf.RunStepStop(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepEnd:
		{
			if tf.RunStepEnd(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	}
}

func (tf *TRFrame) Run() {
	tf.evHub.Run()
}

func (tf *TRFrame) Stop() {
	tf.evHub.Stop()
}

func (tf *TRFrame) changeNextStep() {
	tf.runStep++
	loghlp.Infof("change tframe_step(%d) before step:%d", tf.runStep, tf.runStep-1)
}

func (tf *TRFrame) listen() error {
	var listenMode int32
	var listenAddr string

	switch tf.nodeType {
	case trnode.ETRNodeTypeRoot:
		{
			rootCfg := tf.frameConfig.RootCfgs[tf.nodeIndex]
			if rootCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = rootCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeGate:
		{
			listenCfg := tf.frameConfig.GateCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeHallGate:
		{
			listenCfg := tf.frameConfig.HallGateCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeHallServer:
		{
			listenCfg := tf.frameConfig.HallServerCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeHallData:
		{
			listenCfg := tf.frameConfig.HallDataCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeHallRoom:
		{
			listenCfg := tf.frameConfig.HallRoomCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeHallRoomMgr:
		{
			listenCfg := tf.frameConfig.HallRoomMgrCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	default:
		{
			loghlp.Warnf("not find listen config,nodeType:%d", tf.nodeType)
			break
		}
	}
	if len(listenAddr) == 0 {
		return errors.New("not find listenaddr")
	}
	loghlp.Infof("listen addr:%s", listenAddr)
	return tf.evHub.Listen(listenMode, listenAddr)
}

func (tf *TRFrame) makeInitWorkNode(nType int32, index int32) ITRFrameWorkNode {
	switch nType {
	case trnode.ETRNodeTypeRoot:
		{
			return troot.New(tf, index)
		}
	case trnode.ETRNodeTypeGate:
		{
			return tgate.New(tf, index)
		}
	case trnode.ETRNodeTypeHallGate:
		{
			return thallgate.New(tf, index)
		}
	case trnode.ETRNodeTypeHallServer:
		{
			//return thallserver.New(tf, index)
			return thallserver.New(tf, index)
		}
	case trnode.ETRNodeTypeHallData:
		{
			return thalldata.New(tf, index)
		}
	case trnode.ETRNodeTypeHallRoom:
		{
			return thallroom.New(tf, index)
		}
	case trnode.ETRNodeTypeHallRoomMgr:
		{
			return thallroommgr.New(tf, index)
		}
	default:
		{
			loghlp.Info("unknonw init worknode type:%d", nType)
			break
		}
	}
	return nil
}
func (tf *TRFrame) initNodeSetting() {
	switch tf.nodeType {
	case trnode.ETRNodeTypeRoot:
		{

		}
	case trnode.ETRNodeTypeGate:
		{

		}
	case trnode.ETRNodeTypeHallGate:
		{
			nodeCfg := tf.frameConfig.HallGateCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("hallgate%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
			break
		}
	case trnode.ETRNodeTypeHallServer:
		{
			nodeCfg := tf.frameConfig.HallServerCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("hallserver%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
			break
		}
	case trnode.ETRNodeTypeHallData:
		{
			nodeCfg := tf.frameConfig.HallDataCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("halldata%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
			break
		}
	case trnode.ETRNodeTypeHallRoom:
		{
			nodeCfg := tf.frameConfig.HallRoomCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("hallroom%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
			if len(nodeCfg.CsvPath) > 0 {
				configdata.InitConfigData(nodeCfg.CsvPath) // 加载csv配置
			}
			break
		}
	case trnode.ETRNodeTypeHallRoomMgr:
		{
			nodeCfg := tf.frameConfig.HallRoomMgrCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("room_mgr%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
			break
		}
	default:
		{
			loghlp.Info("unknonw init worknode type:%d", tf.nodeType)
			break
		}
	}
}
func (tf *TRFrame) stopFrame() {
	if tf.runStep < ETRFrameStepStop {
		tf.runStep = ETRFrameStepStop
	}
}
func (tf *TRFrame) GetEvHub() *evhub.EventHub {
	return tf.evHub
}
func (tf *TRFrame) GetFrameConfig() *tframeconfig.FrameConfig {
	return tf.frameConfig
}
func (tf *TRFrame) AfterMsgJob(doJob func()) {
	tf.msgDoneList = append(tf.msgDoneList, doJob)
}
func (tf *TRFrame) updateKeepNodeAlive(curTimeMs int64) {
	if curTimeMs-tf.keepNodeTime < 3000 {
		return
	}
	tf.keepNodeTime = curTimeMs
	// 心跳
	for _, tn := range tf.frameNodeMgr.nodeList {
		if !tn.IsServerClient() {
			continue
		}
		if curTimeMs/1000-tn.LastHeartTime() >= 15 {
			sendMsg := MakeInnerMsg(protocol.EMsgClassFrame, protocol.EFrameMsgKeepNodeHeart, make([]byte, 0))
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("keep node heart callback suss:%d", okCode)
			}
			callInfo := tf.remoteMsgMgr.makeCallInfo(protocol.EMsgClassFrame,
				protocol.EFrameMsgKeepNodeHeart,
				cb,
				nil)
			setSecondHead(sendMsg, 0, callInfo.reqID, 0)
			tn.SendMsg(sendMsg)
			tn.SetHeartTime(curTimeMs / 1000)
		}
	}
}

func (tf *TRFrame) registerFrameHandler() {
	tf.msgDispatcher.RegisterMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgRegisterServerInfo,
		handleRegisterNodeInfo)
}
