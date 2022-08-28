package evhub

import (
	"encoding/binary"
	"errors"
)

const (
	MsgHeadSize      = 12
	SecondHeaderSize = 24
)

type INetMessage interface {
	MainType()
	SubType()
	Encode() []byte
	Decode(d []byte)
}

// 附加头信息
type NetMsgSecondHead struct {
	ID    int64
	ReqID uint64
	RepID uint64
}

type NetMsgHead struct {
	Len        uint32
	MsgClass   uint16
	MsgType    uint16
	MsgVersion uint8
	HasSecond  uint8
	Result     uint16
	// // 附加头部信息(服务器内部用)
	// SecondHead NetMsgSecondHead
}

func (head *NetMsgHead) Encode() []byte {
	buf := make([]byte, MsgHeadSize)
	binary.LittleEndian.PutUint32(buf, head.Len)
	binary.LittleEndian.PutUint16(buf[4:6], head.MsgClass)
	binary.LittleEndian.PutUint16(buf[6:8], head.MsgType)
	buf[8] = head.MsgVersion
	buf[9] = head.HasSecond
	binary.LittleEndian.PutUint16(buf[10:], head.Result)
	return buf
}
func (head *NetMsgHead) Decode(data []byte) {
	head.Len = binary.LittleEndian.Uint32(data[0:4])
	head.MsgClass = binary.LittleEndian.Uint16(data[4:6])
	head.MsgType = binary.LittleEndian.Uint16(data[6:8])
	head.MsgVersion = uint8(data[8])
	head.HasSecond = uint8(data[9])
	head.Result = binary.LittleEndian.Uint16(data[10:])
}

func (head2 *NetMsgSecondHead) Encode() []byte {
	buf := make([]byte, SecondHeaderSize)
	binary.LittleEndian.PutUint64(buf[0:8], uint64(head2.ID))
	binary.LittleEndian.PutUint64(buf[8:16], head2.ReqID)
	binary.LittleEndian.PutUint64(buf[16:24], head2.RepID)
	return buf
}
func (head2 *NetMsgSecondHead) Decode(data []byte) {
	head2.ID = int64(binary.LittleEndian.Uint64(data[0:8]))
	head2.ReqID = binary.LittleEndian.Uint64(data[8:16])
	head2.RepID = binary.LittleEndian.Uint64(data[16:24])
}

type NetMessage struct {
	Head *NetMsgHead
	// 附加头部信息(服务器内部用)
	SecondHead *NetMsgSecondHead
	Data       []byte
	RecvTime   int64
}

func (msg *NetMessage) Encode() []byte {
	msg.Head.Len = uint32(MsgHeadSize + len(msg.Data))
	hd := msg.Head.Encode()
	if msg.Head.HasSecond > 0 {
		hd2 := msg.SecondHead.Encode()
		hd = append(hd, hd2...)
	}
	return append(hd, msg.Data...)
}

func (msg *NetMessage) Decode(data []byte) error {
	if len(data) < MsgHeadSize {
		return errors.New("msg decode size error")
	}
	msg.Head.Decode(data[:MsgHeadSize])
	if msg.Head.HasSecond > 0 {
		msg.SecondHead = &NetMsgSecondHead{}
		msg.SecondHead.Decode(data[MsgHeadSize : MsgHeadSize+SecondHeaderSize])
	}
	msg.Data = data[MsgHeadSize:]
	return nil
}

func NewNetMessage(hd *NetMsgHead, data []byte) *NetMessage {
	return &NetMessage{
		Head: hd,
		Data: data,
	}
}

func MakeMessage(msgClass int32, msgType int32, data []byte) *NetMessage {
	hd := &NetMsgHead{
		MsgClass:  uint16(msgClass),
		MsgType:   uint16(msgType),
		HasSecond: 0,
	}
	msg := NewNetMessage(hd, data)
	return msg
}
func MakeEmptyMessage() *NetMessage {
	hd := &NetMsgHead{
		MsgClass:  0,
		MsgType:   0,
		HasSecond: 0,
	}
	msg := NewNetMessage(hd, nil)
	return msg
}
