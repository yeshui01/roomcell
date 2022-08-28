package evhub

import "encoding/binary"

// 客户端消息结构
type ClientMessageDef struct {
	MsgClass uint16 // 消息大类别
	MsgType  uint16 // 消息大类别下的小类型
	SeqID    int32  // 客户端用请求序列id(用作客户端回调)
	Code     uint16 // 错误码
	MsgData  []byte // proto的序列化数据
}

// 客户端消息和服务器消息的转化
func EncodeServerMsgToClientMsg(sMsg *NetMessage) []byte {
	buf := make([]byte, 10)
	binary.LittleEndian.PutUint16(buf[0:2], sMsg.Head.MsgClass)
	binary.LittleEndian.PutUint16(buf[2:4], sMsg.Head.MsgType)
	if sMsg.SecondHead != nil {
		binary.LittleEndian.PutUint32(buf[4:8], uint32(sMsg.SecondHead.ReqID))
	} else {
		binary.LittleEndian.PutUint32(buf[4:8], uint32(0))
	}

	binary.LittleEndian.PutUint16(buf[8:10], sMsg.Head.Result)
	buf = append(buf, sMsg.Data...)
	return buf
}

// 解析客户端消息结构为本地服务器消息结构
func DecodeClientMsgToServerMsg(fullClientData []byte) *NetMessage {
	sMsg := MakeEmptyMessage()
	sMsg.Head.MsgClass = binary.LittleEndian.Uint16(fullClientData[0:2])
	sMsg.Head.MsgType = binary.LittleEndian.Uint16(fullClientData[2:4])
	SeqID := binary.LittleEndian.Uint32(fullClientData[4:8])
	sMsg.Head.Result = binary.LittleEndian.Uint16(fullClientData[8:10])
	sMsg.SecondHead = &NetMsgSecondHead{
		ReqID: uint64(SeqID),
	}
	sMsg.Data = append(sMsg.Data, fullClientData[10:]...)
	return sMsg
}
