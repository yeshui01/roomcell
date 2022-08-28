package robotcore

import (
	"fmt"
	"log"
	"net/url"
	"roomcell/app/account/accrouter"
	"roomcell/pkg/crossdef"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbtools"
	"roomcell/pkg/protocol"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe"
	"roomcell/pkg/webreq"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type ICellRobot interface {
	GetRobotName() string
	LogRecvMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message)
}
type RobotMsgHandler func(robotIns ICellRobot, sMsg *evhub.NetMessage)

type RobotCallEnv struct {
	CallbackFunc func(sMsg *evhub.NetMessage)
	BeginTime    int64
	MsgClass     int32
	MsgType      int32
}
type IRobotAI interface {
	Update(curTime int64)
}
type CellRobot struct {
	WsConn    *websocket.Conn
	recvMsgCh chan *evhub.NetMessage // 收到的消息队列
	//sendMsgCh   chan *evhub.NetMessage // 发送消息的队列
	RobotName    string
	UserID       int64
	Icon         int32
	SeqID        int32
	ClosedWrite  bool
	HeartTime    int64
	StopRun      bool
	MsgHandlers  map[int32]RobotMsgHandler
	AsyncCall    map[int32]*RobotCallEnv
	RoomID       int64 // 当前房间id
	TargetRoomID int64 // 指定进入的房间id
	UpdateFunc   func(curTime int64)
	// status
	RobotStatus  int32
	StatusTime   int64
	StatusStep   int32 // 状态阶段	0-进入 1-状态完成
	Token        string
	RealRobotIns ICellRobot
	//
	AiInstance IRobotAI
	// data
	HallAddr   string
	RoomDetail *pbclient.RoomData
	IsReady    bool // 是否准备游戏了
	// robot action
	LastActionTime int64  // 最近执行的行为时间
	ActionName     string // 行为名字标示
	ActionTime     int64  // 行为开始时间 0-没有执行任何行为
	ActionStep     int32  // 行为阶段 0--进行中 1-完成结束了
	ActionSeq      int32
}

func NewCellRobot(robotName string) *CellRobot {
	r := &CellRobot{
		recvMsgCh: make(chan *evhub.NetMessage),
		//sendMsgCh:   make(chan *evhub.NetMessage),
		RobotName:    robotName,
		UserID:       0,
		SeqID:        0,
		ClosedWrite:  false,
		HeartTime:    0,
		StopRun:      false,
		MsgHandlers:  make(map[int32]RobotMsgHandler),
		AsyncCall:    make(map[int32]*RobotCallEnv),
		RoomID:       0,
		RealRobotIns: nil,
		ActionTime:   0,
		ActionStep:   0,
		ActionSeq:    0,
		IsReady:      false,
	}
	//r.InitRegisterCommnHandle()
	return r
}

func (robotObj *CellRobot) WSConnect(servAddr string) error {
	u := url.URL{Scheme: "ws", Host: servAddr, Path: "/ws"}
	loghlp.Infof("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	robotObj.WsConn = c
	go func() {
		robotObj.startRead()
	}()

	return nil
}

// 日志操作
func (robotObj *CellRobot) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
func (robotObj *CellRobot) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}
func (robotObj *CellRobot) Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}
func (robotObj *CellRobot) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (robotObj *CellRobot) startRead() {
	for {
		msgType, message, err := robotObj.WsConn.ReadMessage()
		if err != nil {
			robotObj.Debugf("read error:%s", err.Error())
			break
		}
		switch msgType {
		case websocket.TextMessage:
			{
				robotObj.Debugf("recv txtmessage:%s", string(message))
				break
			}
		case websocket.BinaryMessage:
			{
				sMsg := evhub.DecodeClientMsgToServerMsg(message)
				//robotObj.Debugf("recv BinaryMessage(%d_%d)SeqID:%d", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.SecondHead.ReqID)
				if sMsg.Head.MsgClass == protocol.ECMsgClassPlayer && sMsg.Head.MsgType == protocol.ECMsgPlayerLoginHall {
					pbRep := &pbclient.ECMsgPlayerLoginHallRsp{}
					errParse := proto.Unmarshal(sMsg.Data, pbRep)
					if trframe.DecodePBMessage(sMsg, pbRep) {
						robotObj.Debugf("recv ECMsgPlayerLoginHallRsp:%+v", pbRep)
					} else {
						robotObj.Errorf("decode ECMsgPlayerLoginHallRsp error:%s", errParse.Error())
					}
				}
				robotObj.recvMsgCh <- sMsg
				break
			}
		}
	}
	robotObj.ClosedWrite = true
	robotObj.Infof("robot(%s) exit read", robotObj.RobotName)
}

// func (robotObj *CellRobot) startWrite() {
// 	for {
// 		select {
// 		case sMsg := <-robotObj.sendMsgCh:
// 			{
// 				robotObj.Debugf("will send msg:%+v", sMsg)
// 			}
// 		}
// 	}
// }

func (robotObj *CellRobot) GenActionSeq() int32 {
	robotObj.ActionSeq++
	if robotObj.ActionSeq >= 99999999 {
		robotObj.ActionSeq = 1
	}
	return robotObj.ActionSeq
}

// 执行一个行为(一连串的消息),设定,一个机器人，同时只能做一件事情
func (robotObj *CellRobot) DoAction(actionName string, doAnything func()) bool {
	if robotObj.ActionTime != 0 {
		robotObj.Errorf("robot(%s) do action fail, action(%s) is doing", robotObj.RobotName, robotObj.ActionName)
		return false
	}
	// 可以执行
	robotObj.ActionName = fmt.Sprintf("%s_%d", actionName, robotObj.GenActionSeq())
	robotObj.ActionTime = timeutil.NowTime()
	robotObj.ActionStep = 0
	robotObj.Infof("robot(%s) start action(%s)", robotObj.RobotName, robotObj.ActionName)
	doAnything()
	return true
}

func (robotObj *CellRobot) EndAction() {
	robotObj.Infof("robot(%s) end action(%s)", robotObj.RobotName, robotObj.ActionName)
	robotObj.LastActionTime = robotObj.ActionTime
	robotObj.ActionTime = 0
	robotObj.ActionStep = 0
	robotObj.ActionName = ""
}

func (robotObj *CellRobot) GenSeqID() int32 {
	robotObj.SeqID++
	if robotObj.SeqID >= 999999999 {
		robotObj.SeqID = 1
	}
	return robotObj.SeqID
}
func (robotObj *CellRobot) SendMsgToServer(msgClass int32, msgType int32, pbMsg proto.Message) {
	if robotObj.ClosedWrite {
		robotObj.Warnf("writeclosed, robot(%s) ingnore send msg(%d_%d)", robotObj.RobotName, msgClass, msgType)
		return
	}
	sMsg := trframe.MakePBMessage(msgClass, msgType, pbMsg, protocol.ECodeSuccess)
	sMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(robotObj.GenSeqID()),
	}
	sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
	err := robotObj.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
	robotObj.Debugf("robot(%s) send msg(%s_%s)[%s]:%+v", robotObj.RobotName, msgClass, msgType, pbtools.GetFullNameByMessage(pbMsg), pbMsg)
	if err != nil {
		loghlp.Errorf("write binary error:%s", err.Error())
	}
}

func (robotObj *CellRobot) LogCbMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	robotObj.Debugf("log robot(%s) cbmsg(%d_%d)[%d][%s]isok[%d]",
		robotObj.RobotName,
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		netMsg.SecondHead.ReqID,
		pbtools.GetFullNameByMessage(pbMsg),
		netMsg.Head.Result,
	)
}
func (robotObj *CellRobot) LogRecvMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	robotObj.Debugf("log robot(%s) recv msg(%d_%d)[%d][%s]isok[%d]",
		robotObj.RobotName,
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		netMsg.SecondHead.ReqID,
		pbtools.GetFullNameByMessage(pbMsg),
		netMsg.Head.Result,
	)
}
func (robotObj *CellRobot) SendKeepHeart() {
	reqMsg := &pbclient.ECMsgPlayerKeepHeartReq{}
	//robotObj.SendMsgToServer(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerKeepHeart, reqMsg)
	robotObj.RemoteCall(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerKeepHeart, reqMsg, func(sMsg *evhub.NetMessage) {
		rsp := &pbclient.ECMsgPlayerKeepHeartRsp{}
		if !trframe.DecodePBMessage(sMsg, rsp) {
			return
		}
		robotObj.LogCbMsgInfo(sMsg, rsp)
	})
}
func (robotObj *CellRobot) GetRobotName() string {
	return robotObj.RobotName
}
func (robotObj *CellRobot) Update(curTime int64) {
	if robotObj.AiInstance != nil {
		robotObj.AiInstance.Update(curTime)
	}
}
func (robotObj *CellRobot) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for !robotObj.StopRun {
		select {
		case <-ticker.C:
			{
				curTime := time.Now().Unix()
				//robotObj.Debugf("robot(%s) sectick?:%s", robotObj.RobotName, t.String())

				if curTime-robotObj.HeartTime >= 15 {
					// 发送心跳
					robotObj.SendKeepHeart()
					robotObj.HeartTime = curTime
				}
				// 异步调用超时处理
				for k, v := range robotObj.AsyncCall {
					if curTime-v.BeginTime >= 30 {
						robotObj.Errorf("asynccall(%d_%d),seqid(%d) timeout", v.MsgClass, v.MsgType, k)
						delete(robotObj.AsyncCall, k)
					}
				}
				// 机器人行为时间检测
				if robotObj.ActionTime > 0 {
					if robotObj.ActionStep == 0 {
						if curTime-robotObj.ActionTime > 30 {
							// 欣慰超时了,异常,解除锁定
							robotObj.LastActionTime = robotObj.ActionTime
							robotObj.ActionTime = 0
							robotObj.ActionStep = 0
						}
					}
				}

				// if robotObj.UpdateFunc != nil {
				// 	robotObj.UpdateFunc(curTime)
				// }
				robotObj.Update(curTime)
				break
			}
		case sMsg, ok := <-robotObj.recvMsgCh:
			{
				if ok {
					robotObj.Debugf("recv BinaryMessage(%d_%d)SeqID:%d", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.SecondHead.ReqID)
					robotObj.DispatchRobotMsg(sMsg)
				} else {
					robotObj.Errorf("maybe recvMsgCh closed !!!")
				}
				break
			}
		}
	}

	robotObj.Infof("robot(%s) stop", robotObj.RobotName)

	close(robotObj.recvMsgCh)
	time.Sleep(time.Second * 2)
	// reqMsg := &pbclient.ECMsgPlayerKeepHeartReq{}
	// robotObj.SendMsgToServer(0, 0, reqMsg)

	robotObj.Infof("robot(%s) exit run", robotObj.RobotName)
}

func (robotObj *CellRobot) ChangeRobotStatus(robotStatus int32, step int32) {
	robotObj.RobotStatus = robotStatus
	robotObj.StatusTime = timeutil.NowTime()
	robotObj.StatusStep = step
}

func (robotObj *CellRobot) RemoteCall(msgClass int32, msgType int32, pbMsg proto.Message, cbFunc func(sMsg *evhub.NetMessage)) {
	if robotObj.ClosedWrite {
		robotObj.Warnf("writeclosed, robot(%s) remotecall ingnore send msg(%d_%d)", robotObj.RobotName, msgClass, msgType)
		return
	}
	sMsg := trframe.MakePBMessage(msgClass, msgType, pbMsg, protocol.ECodeSuccess)
	sMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(robotObj.GenSeqID()),
	}
	robotObj.AsyncCall[int32(sMsg.SecondHead.ReqID)] = &RobotCallEnv{
		BeginTime:    timeutil.NowTime(),
		CallbackFunc: cbFunc,
		MsgClass:     msgClass,
		MsgType:      msgType,
	}
	sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
	err := robotObj.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
	robotObj.Debugf("robot(%s) send call msg(%d_%d)seqid(%s),userid(%d):%+v",
		robotObj.RobotName,
		msgClass,
		msgType,
		pbtools.GetFullNameByMessage(pbMsg),
		robotObj.UserID,
		pbMsg)
	if err != nil {
		loghlp.Errorf("write binary error:%s", err.Error())
	}
}

// 注册账号
func (robotObj *CellRobot) RegisterAccount(hostAddr string) bool {
	regReq := &accrouter.AccountRegisterReq{
		UserName: robotObj.RobotName,
		Pswd:     "123456",
	}
	err, repData := webreq.PostJson(fmt.Sprintf("http://%s/account/register", hostAddr), regReq)
	if err != nil {
		robotObj.Errorf("register account(%s) error:%s", robotObj.RobotName, err.Error())
		return false
	}
	type RegisterMsgRsp struct {
		Data *accrouter.AccountRegisterRsp
		Code int32
		Msg  string
	}
	regRep := &RegisterMsgRsp{}
	errJv := json.Unmarshal(repData, regRep)
	if errJv != nil {
		robotObj.Errorf("unmarshal RegisterMsgRsp error:%s", errJv.Error())
		return false
	}
	robotObj.UserID = regRep.Data.UserID
	robotObj.Infof("robot(%s) register success,user_id(%d)", robotObj.RobotName, robotObj.UserID)
	return true
}

// 登录账号
func (robotObj *CellRobot) LoginAccount(hostAddr string) bool {
	regReq := &accrouter.AccountLoginReq{
		UserName: robotObj.RobotName,
		Pswd:     "123456",
	}
	err, repData := webreq.PostJson(fmt.Sprintf("http://%s/account/login", hostAddr), regReq)
	if err != nil {
		robotObj.Errorf("login account(%s) error:%s", robotObj.RobotName, err.Error())
		return false
	}
	type LoginMsgRsp struct {
		Data *accrouter.AccountLoginRsp
		Code int32
		Msg  string
	}
	loginRep := &LoginMsgRsp{}
	errJv := json.Unmarshal(repData, loginRep)
	if errJv != nil {
		robotObj.Errorf("unmarshal RegisterMsgRsp error:%s", errJv.Error())
		return false
	}
	if loginRep.Code != protocol.ECodeSuccess {
		robotObj.Errorf("login account error(%d)", loginRep.Code)
		return false
	}
	// 验证touken
	ok, tokenRes := crossdef.TokenAuthClaims(loginRep.Data.Token, crossdef.SignKey)
	if !ok {
		robotObj.Errorf("token parse fail")
		return false
	} else {
		loghlp.Infof("parse player token success:%+v", *tokenRes)
	}
	robotObj.HallAddr = loginRep.Data.HallAddr
	robotObj.UserID = tokenRes.UserID
	robotObj.Token = loginRep.Data.Token
	robotObj.Infof("robot(%s) login account success,user_id(%d):%+v", robotObj.RobotName, robotObj.UserID, loginRep.Data)
	return true
}

// 登录大厅
func (robotObj *CellRobot) LoginHall() {
	robotObj.ChangeRobotStatus(ERobotStatusLoginHall, ERobotStatusStepIng)
	reqMsg := &pbclient.ECMsgPlayerLoginHallReq{
		Token: robotObj.Token,
	}
	robotObj.RemoteCall(protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerLoginHall,
		reqMsg,
		func(sMsg *evhub.NetMessage) {
			rsp := &pbclient.ECMsgPlayerLoginHallRsp{}
			if !trframe.DecodePBMessage(sMsg, rsp) {
				return
			}
			robotObj.LogCbMsgInfo(sMsg, rsp)
			if sMsg.Head.Result == protocol.ECodeSuccess {
				robotObj.Debugf("robot(%s) login hall succ", robotObj.RobotName)
				robotObj.Icon = rsp.RoleData.Icon
				robotObj.ChangeRobotStatus(ERobotStatusLoginHall, ERobotStatusStepFinish)
			}
		},
	)
}

// 进入房间
func (robotObj *CellRobot) SendEnterRoom(roomID int64) {
	reqMsg := &pbclient.ECMsgRoomEnterReq{
		RoomID: roomID,
	}
	robotObj.ChangeRobotStatus(ERobotStatusEnterRoom, ERobotStatusStepIng)
	robotObj.RemoteCall(protocol.ECMsgClassRoom,
		protocol.ECMsgRoomEnter,
		reqMsg,
		func(sMsg *evhub.NetMessage) {
			rsp := &pbclient.ECMsgRoomEnterRsp{}
			if !trframe.DecodePBMessage(sMsg, rsp) {
				return
			}
			robotObj.LogCbMsgInfo(sMsg, rsp)
			if sMsg.Head.Result == protocol.ECodeSuccess {
				robotObj.Debugf("robot(%s) enter room(%d) succ:%+v", robotObj.RobotName, roomID, rsp.RoomDetail)
				robotObj.RoomDetail = rsp.RoomDetail
				robotObj.RoomID = rsp.RoomDetail.RoomID
				robotObj.ChangeRobotStatus(ERobotStatusEnterRoom, ERobotStatusStepFinish)
			}
		},
	)
}

// 游戏准备 readyOpt: 0-取消 1-准备
func (robotObj *CellRobot) DoActGameReady(readyOpt int32) {
	robotObj.DoAction(ActGameReady, func() {
		reqMsg := &pbclient.ECMsgGameReadyOptReq{
			Ready: readyOpt,
		}
		robotObj.RemoteCall(protocol.ECMsgClassGame,
			protocol.ECMsgGameReadyOpt,
			reqMsg,
			func(sMsg *evhub.NetMessage) {
				rsp := &pbclient.ECMsgGameReadyOptRsp{}
				if !trframe.DecodePBMessage(sMsg, rsp) {
					return
				}
				robotObj.LogCbMsgInfo(sMsg, rsp)
				if sMsg.Head.Result == protocol.ECodeSuccess {
					robotObj.IsReady = true
					robotObj.EndAction()
				}
			},
		)
	})
}
