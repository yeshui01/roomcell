package iframe

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/tframeconfig"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	ESessionDataType        = 0
	ESessionDataTypeNetInfo = 1
)

// 消息上下文
type TRRemoteMsgEnv struct {
	SrcSessionID int32
	SrcMessage   *evhub.NetMessage
	UserData     interface{}
}

type MsgCallbackFunc func(okCode int32, msgData []byte, env *TRRemoteMsgEnv)

type ITRFrame interface {
	Run()
	GetEvHub() *evhub.EventHub
	GetFrameConfig() *tframeconfig.FrameConfig
	ForwardMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32, cb MsgCallbackFunc, env *TRRemoteMsgEnv) bool
}

// session数据
type SessionUserData struct {
	DataType       int32
	NodeType       int32
	NodeIndex      int32
	DesInfo        string
	IsServerClient bool
}
type ISession interface {
	SendMsg(msg *evhub.NetMessage) bool
}

// 消息分发器
type IHandleResultType int32

const (
	EHandleContent IHandleResultType = 0 // 返回消息
	EHandlePending IHandleResultType = 1 // 稍后处理
	EHandleNone    IHandleResultType = 2 // 不需要处理
)

// 消息上下文
type TMsgContext struct {
	FrameInstance ITRFrame
	Session       ISession
	NetMessage    *evhub.NetMessage
	CustomData    interface{}
}
type FrameMsgHandler func(tmsgCtx *TMsgContext) (isok int32, retData interface{}, rt IHandleResultType)
type IFrameMsgDispathcer interface {
	Dispatch(session ISession, msg *evhub.NetMessage, customData interface{})
	RegisterMsgHandler(msgClass int32, msgType int32, msgHander FrameMsgHandler)
}

func ToPbData(pbMsg protoreflect.ProtoMessage) ([]byte, error) {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("ToPbData fail, proto marsh error:%s", err.Error())
		return nil, err
	}
	return data, err
}
