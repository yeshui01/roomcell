package iroom

import "google.golang.org/protobuf/proto"

type IRoomGlobal interface {
	BroadcastMsgToClient(msgClass int32, msgType int32, pbMsg proto.Message, roleList []int64)
}
