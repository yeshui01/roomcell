package hallgatemain

import (
	"encoding/binary"
	"roomcell/pkg/evhub"
	"roomcell/pkg/funchelp"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	WsMsgEncodeHeadSize = 9
	WsMsgIdClassRatio   = 1000
	WsMsgVersion        = 1
)

// evhub的data为pb数据,这里实现pb和json消息互转
func ZipCompressHubMessage(msg *evhub.NetMessage) ([]byte, error) {
	jsonData := msg.Data // TODO:这里需要根据pb消息反序列化
	var lenJosnOrigin uint32 = uint32(len(jsonData))
	var msgTag uint8 = WsMsgVersion
	var clientSeq uint32 = 0
	var msgID uint32 = uint32(msg.Head.MsgClass)*WsMsgIdClassRatio + uint32(msg.Head.MsgType)
	dataPreCompress := make([]byte, len(jsonData)+WsMsgEncodeHeadSize)
	binary.LittleEndian.PutUint32(dataPreCompress[0:4], msgID)
	dataPreCompress[4] = msgTag
	binary.LittleEndian.PutUint32(dataPreCompress[5:9], clientSeq)
	copy(dataPreCompress[9:], jsonData)
	compressData := funchelp.ZlibCompress(dataPreCompress)
	finalData := make([]byte, len(compressData)+4)
	binary.LittleEndian.PutUint32(finalData[0:4], lenJosnOrigin)
	return finalData, nil
}

func UnzipToHubMessage(msgData []byte) (*evhub.NetMessage, error) {
	var originJsonLen uint32 = binary.LittleEndian.Uint32(msgData[0:4])
	unCompressData := funchelp.ZlibUnCompress(msgData[4:])
	var msgID uint32 = binary.LittleEndian.Uint32(unCompressData[0:4])
	var msgTag uint8 = unCompressData[4]
	if msgTag != WsMsgVersion {
		loghlp.Errorf("msg version tag not match(%d,%d)", msgTag, WsMsgVersion)
	}
	var clientSeq uint32 = binary.LittleEndian.Uint32(unCompressData[5:9])
	jsonData := unCompressData[9:]

	evMsg := evhub.MakeMessage(int32(msgID/WsMsgIdClassRatio), int32(msgID%WsMsgIdClassRatio), jsonData)
	evMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(clientSeq),
	}
	loghlp.Debugf("unzipMessage(%d_%d),originJsonLen:%d, nowLen:%d",
		evMsg.Head.MsgClass,
		evMsg.Head.MsgType,
		originJsonLen,
		len(unCompressData)-WsMsgEncodeHeadSize,
	)
	return evMsg, nil
}

// --------------------------------[proto msg reflect]--------------------------------------
func ConvertHubMsgPbDataToJsonData(hubPbData []byte) []byte {
	// TODO
	return nil
}

func ConvertJsonDataToHubMsgPbData() []byte {
	// TODO
	return nil
}

var (
	reqParseMap map[int32]func() proto.Message
	repParseMap map[int32]func() proto.Message
)

func getMsgKey(msgClass int32, msgType int32) int32 {
	return msgClass*1000 + msgType
}

// 请求
func RegisterProtoReqParseObj(msgClass int32, msgType int32, objFun func() proto.Message) {
	if reqParseMap == nil {
		reqParseMap = make(map[int32]func() protoreflect.ProtoMessage)
	}
	var msgKey int32 = getMsgKey(msgClass, msgType)
	reqParseMap[msgKey] = objFun
}

func GetReqPbObjectByMsgType(msgClass int32, msgType int32) proto.Message {
	var msgKey int32 = getMsgKey(msgClass, msgType)
	if fun, ok := reqParseMap[msgKey]; ok {
		return fun()
	}
	return nil
}

// 回复
func RegisterProtoRepParseObj(msgClass int32, msgType int32, objFun func() proto.Message) {
	if reqParseMap == nil {
		reqParseMap = make(map[int32]func() protoreflect.ProtoMessage)
	}
	var msgKey int32 = getMsgKey(msgClass, msgType)
	reqParseMap[msgKey] = objFun
}

func GetRepPbObjectByMsgType(msgClass int32, msgType int32) proto.Message {
	var msgKey int32 = getMsgKey(msgClass, msgType)
	if fun, ok := reqParseMap[msgKey]; ok {
		return fun()
	}
	return nil
}

// 注册所有解析对象分配器
func InitProtoParseObject() {
	//
	RegisterProtoReqParseObj(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall,
		func() proto.Message { return &pbframe.FrameMsgRegisterServerInfoReq{} })
	RegisterProtoRepParseObj(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall,
		func() proto.Message { return &pbframe.FrameMsgRegisterServerInfoRep{} })
	//
}
